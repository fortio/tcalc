package main

import (
	"errors"
	"fmt"
	"testing"

	"fortio.org/tcalc/calculator"
	"fortio.org/terminal/ansipixels"
)

func TestDisplayStrings(t *testing.T) {
	binStrings := binaryDisplayStrings(64, 0)
	row4, _ := ansipixels.AnsiClean([]byte(binStrings[4]))
	if string(row4) != "16: 0 0 0 0  0 0 0 0  0 1 0 0  0 0 0 0" {
		t.Fail()
	}
	uintString := uintDisplayString(-64)
	if uintString != "Unsigned : 18446744073709551552" {
		t.Fail()
	}
	if unicodeDisplayString(int64('a')) != "Unicode  : a" {
		t.Fail()
	}
	strs := displayString(64, 0, errors.New("random error"))
	fmt.Println(strs)
	errCheck := "Last input was invalid"
	checkString, _ := ansipixels.AnsiClean([]byte(strs[0]))
	if string(checkString) != errCheck {
		t.Fail()
	}
	octal := octalDisplayString(-64)
	if octal != "Octal    : 0o1777777777777777777700" {
		fmt.Println(octal)
		t.Fail()
	}
}

func TestBitPosition(t *testing.T) {
	c := configure(ansipixels.NewAnsiPixels(30))
	index := c.determineBitFromXY(14, 5)
	if index != 75 {
		t.Fail()
	}
	index = c.determineBitFromXY(15, 5)
	if index != -1 {
		t.Fail()
	}
}

func TestConfigHandleInput(t *testing.T) {
	ap := ansipixels.NewAnsiPixels(30)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		ap   *ansipixels.AnsiPixels
		data []byte
		want bool
	}{
		{"test up", ap, []byte{0x1b, 0x5b, 0x41}, true},
		{"test down", ap, []byte{0x1b, 0x5b, 0x42}, true},
		{"test right", ap, []byte{0x1b, 0x5b, 0x43}, true},
		{"test left", ap, []byte{0x1b, 0x5b, 0x44}, true},
		{"test enter", ap, []byte("\n"), true},
		{"test home", ap, []byte{0x1b, 0x5b, 0x48}, true},
		{"test end", ap, []byte{0x1b, 0x5b, 0x46}, true},
		{"test backspace", ap, []byte{0x7f}, true},
		{"test other", ap, []byte{'a'}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configure(tt.ap)
			c.AP.Data = tt.data
			got, quit := c.handleInput()
			if tt.want != got || quit {
				t.Errorf("handleInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigHandleMouseInput(t *testing.T) {
	ap := ansipixels.NewAnsiPixels(30)
	ap.Mouse = true

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		ap       *ansipixels.AnsiPixels
		data     []byte
		want     bool
		mbuttons int
	}{
		{"test scroll up ", ap, []byte{}, false, ansipixels.MouseWheelUp},
		{"test scroll down ", ap, []byte{}, false, ansipixels.MouseWheelDown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configure(tt.ap)
			c.AP.Data = tt.data
			c.AP.Mbuttons = tt.mbuttons
			got, quit := c.handleInput()
			if tt.want != got || quit {
				t.Errorf("handleInput() = %v, want %v", got, tt.want)
			}
			c.handleMouse()
		})
	}
	t.Run("test clicks", func(_ *testing.T) {
		c := configure(ap)
		for y := 88; y < 94; y++ {
			c.AP.Mouse = true
			c.AP.Mrelease = true
			c.AP.Mbuttons = 0
			c.AP.H = 100
			c.AP.W = 100
			c.AP.My = y
			c.AP.Mx = 30
			c.strings = displayString(5, 20, nil)
			c.history = append(c.history, []historyRecord{{"", 5}, {"daf", 0}}...)
			c.handleMouse()
		}
		c.AP.Mouse = true
		c.AP.Mrelease = true
		c.AP.Mbuttons = 0
		c.AP.H = 100
		c.AP.W = 100
		c.AP.My = 88
		c.AP.Mx = 86
		c.strings = displayString(5, 20, nil)
		c.history = append(c.history, []historyRecord{{"", 5}, {"daf", 0}}...)
		c.handleMouse()
		c.AP.Mbuttons = ansipixels.MouseRight
		c.handleMouse()
	})
}

func TestAssign(t *testing.T) {
	s := calculator.NewState()
	err := s.Exec("x=5")
	if err != nil || s.Variables["x"] != 5 {
		t.Fail()
	}
}

func TestDrawHistory(t *testing.T) {
	s := calculator.NewState()
	err1 := s.Exec("1+1")
	err2 := s.Exec("2+2")
	err3 := s.Exec("3+3")
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fail()
	}
	c := configure(ansipixels.NewAnsiPixels(30))
	c.AP.H = 100
	c.AP.W = 100
	c.state = s
	c.history = append(
		c.history,
		historyRecord{
			evaluated:  "1+1",
			finalValue: 2,
		},
		historyRecord{
			evaluated:  "2+2",
			finalValue: 4,
		},
		historyRecord{
			evaluated:  "3+3",
			finalValue: 6,
		})
	c.curRecord = 2
	c.DrawHistory()
}
