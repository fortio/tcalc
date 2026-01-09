package main

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

type config struct {
	AP        *ansipixels.AnsiPixels
	state     *state
	input     string
	index     int
	bitset    int
	history   []historyRecord
	curRecord int
}

type historyRecord struct {
	evaluated  string
	finalValue int64
}

var validClickXs = []int{
	5, 7, 9, 11, 14, 16, 18, 20, 23, 25, 27, 29, 32, 34, 36, 38,
}

func main() {
	ap := ansipixels.NewAnsiPixels(30)
	c := config{ap, newState(), "", 0, -1, []historyRecord{{"0", 0}}, -1}
	err := c.AP.Open()
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
		if len(ap.Data) > 2 {
			switch string(ap.Data) {
			case "\x1b[H": // home
				c.index = 0
			case "\x1b[F": // home
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
						c.curRecord -= 1

					}

					c.input = c.history[c.curRecord].evaluated

					if c.curRecord > 0 {
						c.input = strings.Replace(c.input, "_ans_", strconv.Itoa(int(c.history[c.curRecord-1].finalValue)), 1)
					}
					c.index = len(c.input)
				}
			case "\x1b[B": // down
				if len(c.history) > 1 {
					c.curRecord = (c.curRecord + 1) % len(c.history)
					c.input = c.history[c.curRecord].evaluated
					if c.curRecord > 0 {
						c.input = strings.Replace(c.history[c.curRecord].evaluated, "_ans_", strconv.Itoa(int(c.history[c.curRecord-1].finalValue)), 1)
					}
					c.index = len(c.input)
				}
			case "\x1b[3~":
				before, after := c.input[:c.index], c.input[min(len(c.input), c.index+1):]
				c.input = before + after
			}
		}
		if len(ap.Data) == 1 {
			switch ap.Data[0] {
			case '\x03':
				return false
			case '\x7f':
				before, after := c.input[:max(0, c.index-1)], c.input[c.index:]
				c.input = before + after
				c.index = max(c.index-1, 0)
			case '\r', '\n':
				if (len(c.input) > 2 && slices.Contains(length2operators, doubleRuneOperator(c.input[:2]))) ||
					(len(c.input) > 0 && slices.Contains(length1operatorsInfix, operator(c.input[0]))) {
					c.input = "_ans_" + c.input
				}
				newRecord := historyRecord{
					evaluated: strings.Replace(c.input, "_ans_-", "-", 1),
				}
				ans := c.state.ans
				newRecord.evaluated = strings.ReplaceAll(newRecord.evaluated, "_ans_", strconv.Itoa(int(ans)))
				err = c.state.Exec(c.input)
				if err != nil {
					c.input = ""
					c.index = 0
					c.state.ans = ans
				}
				newRecord.finalValue = c.state.ans
				if newRecord.evaluated == "" {
					newRecord.evaluated = strconv.Itoa(int(newRecord.finalValue))
				}
				c.history = append(c.history, newRecord)
				c.input, c.index = "", 0
			default:
				c.curRecord = -1
				before, after := c.input[:c.index], c.input[c.index:]
				c.input = before + string(ap.Data) + after
				c.index = c.index + 1
			}
		}
		// fmt.Println(string(ap.Data))
		c.AP.ClearScreen()

		c.AP.WriteAtStr(0, ap.H-11, displayString(c.state.ans))
		c.AP.WriteAtStr(0, c.AP.H, "⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯")
		c.AP.WriteAtStr(0, c.AP.H-2, c.input)
		if c.AP.W > 76 {
			for i := range ap.H {
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
					c.AP.WriteAtStr(ap.W-len(runes), ap.H-((len(c.history)-i)*2)+1, tcolor.Green.Foreground()+string(runes)+tcolor.Reset)
				}

				c.AP.WriteAtStr(ap.W-len(line), ap.H-((len(c.history)-i)*2), line)
				c.AP.WriteAtStr(ap.W-len(runes), ap.H-((len(c.history)-i)*2)-1, tcolor.Green.Foreground()+string(runes)+tcolor.Reset)
			}
		}
		c.AP.MoveCursor(c.index, c.AP.H-2)
		if c.AP.LeftClick() && c.AP.MouseRelease() {
			x, y := c.AP.Mx, c.AP.My
			if slices.Contains(validClickXs, x) && y <= c.AP.H-2 && y >= c.AP.H-6 {
				bit := c.determineBitFromXY(x, c.AP.H-2-y)
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
	}
	c.bitset = bit
	return bit
}
