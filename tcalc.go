package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

type config struct {
	AP           *ansipixels.AnsiPixels
	state        *state
	input        string
	index        int
	bitset       int
	history      []historyRecord
	curRecord    int
	clicked      bool
	clickedValue int64
}

type historyRecord struct {
	evaluated  string
	finalValue int64
}

var validClickXs = []int{
	5, 7, 9, 11, 14, 16, 18, 20, 23, 25, 27, 29, 32, 34, 36, 38,
}

var instructions = []string{
	"Type expressions to evaluate.",
	"SUM +   SUB -   MUL *   DIV /",
	"MOD %   AND &   OR |   XOR ^",
	"POW **  LSHIFT <<   RSHIFT >>",
	"NOT ~   ASSIGN =",
	"Click on individual bits to flip them.",
	"up and down arrows to navigate history.",
	"Press ctrl+c to quit.",
}

func configure(ap *ansipixels.AnsiPixels) config {
	return config{ap, newState(), "", 0, -1, []historyRecord{{"0", 0}}, -1, false, 0}
}

func main() {
	ap := ansipixels.NewAnsiPixels(30)
	c := configure(ap)
	err := c.AP.Open()
	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
	slog.SetDefault(log)
	if err != nil {
		slog.Error("couldn't open terminal", "error", err)
		return
	}
	defer func() {
		c.AP.ShowCursor()
		c.AP.MouseClickOff()
		c.AP.Restore()
		c.AP.ClearScreen()
		fmt.Println(c.state.ans)
	}()

	c.AP.MouseClickOn()
	c.AP.ClearScreen()
	fmt.Println(c)
	if c.AP.W < 38 || c.AP.H < 11 {
		slog.Error("terminal window not large enough")
		return
	}

	err = c.AP.FPSTicks(func() bool {
		c.AP.MoveCursor(c.index+1, c.AP.H-2)
		diff := len(c.history) - (c.AP.H / 2) + 1
		if diff > 0 {
			c.history = c.history[diff:]
		}
		if !c.handleInput() {
			return false
		}
		c.AP.ClearScreen()
		if c.AP.H > 17 {
			for i, str := range instructions {
				c.AP.WriteAtStr(0, i, str)
			}
		}
		strings := displayString(c.state.ans, c.state.err)
		y := ap.H - 13
		for i, str := range strings {
			c.AP.WriteAtStr(0, y+i, str)
		}
		c.AP.WriteAtStr(0, c.AP.H, "⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯")
		c.AP.WriteAtStr(0, c.AP.H-2, c.input)
		c.DrawHistory()
		c.AP.MoveCursor(c.index, c.AP.H-2)
		if c.AP.LeftClick() && c.AP.MouseRelease() {
			x, y := c.AP.Mx, c.AP.My
			if slices.Contains(validClickXs, x) && y < c.AP.H-2 && y >= c.AP.H-6 {
				bit := c.determineBitFromXY(x, c.AP.H-2-y)
				c.clicked = true
				c.state.ans = (c.state.ans) ^ (1 << bit)
			}
		}
		return true
	})
	if err != nil {
		slog.Error("error running fpsticks", "error", err)
	}
}

func (c *config) determineBitFromXY(x, y int) int {
	index := slices.Index(validClickXs, x)
	bit := 0
	if index != -1 {
		bit += (15 - index)
		bit += (16 * (y - 1))
		c.bitset = bit
		return bit
	}
	return -1
}

func (c *config) handleInput() bool {
	switch len(c.AP.Data) {
	case 0:
		return true
	case 1:
		switch c.AP.Data[0] {
		case '\x03':
			return false
		case '\x7f':
			before, after := c.input[:max(0, c.index-1)], c.input[c.index:]
			c.input = before + after
			c.index = max(c.index-1, 0)
		case '\r', '\n':
			c.handleEnter()
		default:
			c.curRecord = -1
			before, after := c.input[:c.index], c.input[c.index:]
			c.input = before + string(c.AP.Data) + after
			c.index++
		}
	default:
		switch string(c.AP.Data) {
		case "\x1b[H": // home
			c.index = 0
		case "\x1b[F": // end
			c.index = len(c.input)
		case "\x1b[C": // right
			c.index = min(c.index+1, len(c.input))
		case "\x1b[D": // left
			c.index = max(c.index-1, 0)
		case "\x1b[A": // up
			if len(c.history) > 1 {
				switch c.curRecord {
				case -1:
					c.curRecord += len(c.history)
				case 0:
					c.curRecord += len(c.history) - 1
				default:
					c.curRecord--
				}
				c.input = c.history[c.curRecord].evaluated
				c.index = len(c.input)
			}
		case "\x1b[B": // down
			if len(c.history) > 1 {
				c.curRecord = (c.curRecord + 1) % len(c.history)
				c.input = c.history[c.curRecord].evaluated
				if c.curRecord > 0 {
					c.input = strings.Replace(c.history[c.curRecord].evaluated, "_ans_",
						strconv.Itoa(int(c.history[c.curRecord-1].finalValue)), 1)
				}
				c.index = len(c.input)
			}
		case "\x1b[3~":
			before, after := c.input[:c.index], c.input[min(len(c.input), c.index+1):]
			c.input = before + after
		}
	}
	return true
}

func (c *config) handleEnter() {
	defer func() { c.clicked = false }()
	if c.input == "" {
		if c.clicked {
			c.input = "(" + strconv.Itoa(int(c.state.ans)) + ")"
		} else {
			c.input = c.history[len(c.history)-1].evaluated
		}
	}
	ansValue := "_ans_"
	if c.clicked {
		ansValue = strconv.Itoa(int(c.state.ans))
	}
	if (len(c.input) >= 2 && slices.Contains(length2operators, doubleRuneOperator(c.input[:2]))) ||
		(len(c.input) > 0 && slices.Contains(length1operatorsInfix, operator(c.input[0]))) {
		c.input = ansValue + c.input
	}
	newRecord := historyRecord{
		evaluated: c.input,
	}
	if len(c.history) > 1 {
		ans := c.history[len(c.history)-2].finalValue
		stringToReplace := strconv.Itoa(int(ans))
		if stringToReplace[0] == '-' {
			stringToReplace = "(" + stringToReplace + ")"
		}
		c.history[len(c.history)-1].evaluated = strings.ReplaceAll(c.history[len(c.history)-1].evaluated, "_ans_", stringToReplace)
	}
	err := c.state.Exec(c.input)
	if err != nil {
		c.input = ""
		c.index = 0
		c.state.ans = c.history[len(c.history)-1].finalValue
		return
	}
	newRecord.finalValue = c.state.ans
	if newRecord.evaluated == "" {
		newRecord.evaluated = strconv.Itoa(int(newRecord.finalValue))
	}
	c.history = append(c.history, newRecord)
	c.input, c.index = "", 0
}

func (c *config) DrawHistory() {
	if c.AP.W > 76 {
		c.AP.WriteAtStr(c.AP.W-27, c.AP.H, "⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯")
		for i := range c.AP.H {
			c.AP.WriteAtStr(c.AP.W/2, i, "⏐")
		}
		for i, record := range c.history {
			line := record.evaluated + ": " + strconv.Itoa(int(record.finalValue))
			runes := make([]rune, len(line), c.AP.W/2-1)
			for i := range line {
				runes[i] = '⎯'
			}
			if c.curRecord == i {
				for j := len(line); j < c.AP.W/2-1; j++ {
					runes = append(runes, '⎯')
				}
				c.AP.WriteAtStr(c.AP.W-len(runes), c.AP.H-((len(c.history)-i)*2)+1, tcolor.Green.Foreground()+string(runes))
			}
			if c.curRecord != i-1 {
				c.AP.WriteAtStr(c.AP.W-len(runes), c.AP.H-((len(c.history)-i)*2)-1, string(runes)+tcolor.Reset)
			}
			c.AP.WriteAtStr(c.AP.W-len(line), c.AP.H-((len(c.history)-i)*2), tcolor.Reset+line)
		}
	}
}
