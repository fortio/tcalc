package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"fortio.org/tcalc/calculator"
	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

type config struct {
	AP           *ansipixels.AnsiPixels
	state        *calculator.State
	input        string
	index        int
	bitset       int
	history      []historyRecord
	curRecord    int
	clicked      bool
	clickedValue int64
	strings      []string
	notification string
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
	return config{ap, calculator.NewState(), "", 0, -1, []historyRecord{{"0", 0}}, -1, false, 0, nil, ""}
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

func (c *config) handleMouse() {
	c.AP.MoveCursor(c.index+1, c.AP.H-2)
	switch {
	case c.AP.MouseWheelUp():
		c.handleUp()
	case c.AP.MouseWheelDown():
		c.handleDown()
	case c.AP.LeftClick() && c.AP.MouseRelease():
		x, y := c.AP.Mx, c.AP.My
		if slices.Contains(validClickXs, x) && y < c.AP.H-2 && y >= c.AP.H-6 {
			bit := c.determineBitFromXY(x, c.AP.H-2-y)
			c.clicked = true
			c.state.Ans = (c.state.Ans) ^ (1 << bit)
			return
		}
		if c.AP.W <= 76 {
			return
		}
		if x <= c.AP.W/2 {
			switch y {
			case c.AP.H - 11:
				c.AP.CopyToClipboard(c.strings[1][len(unicodeString):])
				c.notification = GREEN + "Unicode value copied to clipboard" + tcolor.Reset
			case c.AP.H - 10:
				c.AP.CopyToClipboard(c.strings[2][len(decimalString):])
				c.notification = GREEN + "Decimal value copied to clipboard" + tcolor.Reset
			case c.AP.H - 9:
				c.AP.CopyToClipboard(c.strings[3][len("Unsigned "+decimalString):])
				c.notification = GREEN + "Unsigned decimal value copied to clipboard" + tcolor.Reset
			case c.AP.H - 8:
				c.AP.CopyToClipboard(c.strings[4][len(hexString):])
				c.notification = GREEN + "Hexadecimal value copied to clipboard" + tcolor.Reset
			case c.AP.H - 7:
				if c.clicked {
					c.AP.CopyToClipboard(fmt.Sprintf("%b", c.clickedValue))
				} else {
					c.AP.CopyToClipboard(fmt.Sprintf("%b", c.state.Ans))
				}
				c.notification = GREEN + "Binary value copied to clipboard" + tcolor.Reset
			}
			return
		}
		// we know x > midline
		index := c.recordFromYValue(y)
		if index != -1 {
			c.curRecord = index
			c.input = c.history[c.curRecord].evaluated
			c.index = len(c.input)
		}
	case c.AP.W > 76 && c.AP.RightClick() && c.AP.MouseRelease() && c.AP.Mx > c.AP.W/2:
		index := c.recordFromYValue(c.AP.My)
		if index != -1 {
			c.AP.CopyToClipboard(strconv.Itoa(int(c.history[index].finalValue)))
			c.notification = GREEN + "History copied to clipboard" + tcolor.Reset
		}
	}
}

func (c *config) recordFromYValue(y int) int {
	index := (c.AP.H - y) / 2
	if index < len(c.history) {
		return len(c.history) - 1 - index
	}
	return -1
}

func (c *config) handleInput() bool {
	c.handleMouse()
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
			c.notification = ""
			c.handleEnter()
		default:
			c.notification = ""
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
			c.handleUp()
		case "\x1b[B": // down
			c.handleDown()
		case "\x1b[3~":
			before, after := c.input[:c.index], c.input[min(len(c.input), c.index+1):]
			c.input = before + after
		default:
			before, after := c.input[:c.index], c.input[c.index:]
			c.input = before + string(c.AP.Data) + after
			c.index += len(string(c.AP.Data))
		}
	}
	return true
}

func (c *config) handleDown() {
	if len(c.history) > 1 {
		c.curRecord = (c.curRecord + 1) % len(c.history)
		c.input = c.history[c.curRecord].evaluated
		if c.curRecord == 1 {
			c.input = strings.Replace(c.history[c.curRecord].evaluated, "_ans_",
				strconv.Itoa(int(c.history[c.curRecord-1].finalValue)), 1)
		}
		c.index = len(c.input)
	}
}

func (c *config) handleUp() {
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
}

func (c *config) handleEnter() {
	defer func() { c.clicked = false }()
	if c.input == "" {
		if c.clicked {
			c.input = "(" + strconv.Itoa(int(c.state.Ans)) + ")"
		} else {
			c.input = c.history[len(c.history)-1].evaluated
		}
	}
	trimmed := strings.Trim(c.input, " ")
	lengthTrimmed := len(trimmed)
	if lengthTrimmed >= 2 && (trimmed[lengthTrimmed-2:] == "<<" || trimmed[lengthTrimmed-2:] == ">>") {
		c.input += "1"
	}
	ansValue := "_ans_"
	if c.clicked {
		ansValue = strconv.Itoa(int(c.state.Ans))
	}
	if (len(c.input) >= 2 && slices.Contains(calculator.Length2operators, calculator.DoubleRuneOperator(c.input[:2]))) ||
		(len(c.input) > 0 && slices.Contains(calculator.Length1operatorsInfix, calculator.Operator(c.input[0]))) {
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
	c.curRecord = -1
	err := c.state.Exec(c.input)
	if err != nil {
		c.input = ""
		c.index = 0
		c.state.Ans = c.history[len(c.history)-1].finalValue
		return
	}
	newRecord.finalValue = c.state.Ans
	if newRecord.evaluated == "" {
		newRecord.evaluated = strconv.Itoa(int(newRecord.finalValue))
	}
	c.history = append(c.history, newRecord)
	c.input, c.index = "", 0
}

func (c *config) DrawHistory() {
	if c.AP.W > 76 {
		for i := range 27 {
			c.AP.WriteAtStr(c.AP.W-i, c.AP.H, ansipixels.Horizontal)
			c.AP.WriteAtStr(c.AP.W-i, c.AP.H-((len(c.history))*2)-1, ansipixels.Horizontal)
		}
		for i := range c.AP.H {
			c.AP.WriteAtStr(c.AP.W/2, i, ansipixels.Vertical)
		}
		for i, record := range c.history {
			line := record.evaluated + ": " + strconv.Itoa(int(record.finalValue))
			lengthToUse := len(line)
			line = strings.ReplaceAll(line, "_ans_", italicPrefix+GREEN+"_ans_"+tcolor.Reset)
			runes := make([]string, lengthToUse, c.AP.W)
			for i := range lengthToUse {
				runes[i] = ansipixels.Horizontal
			}
			if c.curRecord == i {
				for j := lengthToUse; j < c.AP.W/2-1; j++ {
					runes = append(runes, ansipixels.Horizontal)
				}
				c.AP.WriteAtStr(c.AP.W-len(runes), c.AP.H-((len(c.history)-i)*2)+1, GREEN+strings.Join(runes, ""))
			}
			if c.curRecord != i-1 {
				c.AP.WriteAtStr(c.AP.W-len(runes), c.AP.H-((len(c.history)-i)*2)-1, strings.Join(runes, "")+tcolor.Reset)
			}
			c.AP.WriteAtStr(c.AP.W-lengthToUse, c.AP.H-((len(c.history)-i)*2), tcolor.Reset+line)
		}
	}
}
