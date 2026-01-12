package main

import (
	"errors"
	"testing"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

func TestDisplayStrings(t *testing.T) {
	binStrings := binaryDisplayStrings(64)
	if binStrings[4] != "16: 0 0 0 0  0 0 0 0  0 1 0 0  0 0 0 0" {
		t.Fail()
	}
	uintString := uintDisplayString(-64)
	if uintString != "Unsigned Decimal: 18446744073709551552" {
		t.Fail()
	}
	if ASCII(int64('a')) != "ASCII: a" {
		t.Fail()
	}
	strs := displayString(64, errors.New("random error"))
	errCheck := tcolor.Red.Foreground() + "Last input was invalid" + tcolor.Reset
	if strs[0] != errCheck {
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

func TestExec(t *testing.T) {
	testCases := []struct {
		expression string
		expected   int64
		shouldFail bool
	}{
		{"1 + 3 + 2", 6, false},
		{"1 * (3 + 2)", 5, false},
		{"(2 * (3 + 2) - 1)+ 1 / 1", 10, false},
		{"(2 * (3 + 2) - 1)+ 1 / 1*3", 30, false},
		{"1 << 5", 32, false},
		{"0+-1", 0, true},
		{"1&2", 0, false},
		{"1|2", 3, false},
		{"1^2", 3, false},
		{"1^3", 2, false},
		{"~1", -2, false},
		{"2>>1", 1, false},
	}
	for _, tc := range testCases {
		s := newState()
		err := s.Exec(tc.expression)
		if tc.shouldFail {
			if err == nil {
				t.Errorf("Expected failure for expression: %s", tc.expression)
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error for expression %s: %v", tc.expression, err)
		} else if s.ans != tc.expected {
			t.Errorf("For expression %s, expected %d but got %d", tc.expression, tc.expected, s.ans)
		}
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
			got := c.handleInput()
			if tt.want != got {
				t.Errorf("handleInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigHandleMouseInput(t *testing.T) {
	ap := ansipixels.NewAnsiPixels(30)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		ap   *ansipixels.AnsiPixels
		data []byte
		want bool
	}{
		{"test mouse click", ap, []byte{0x1b, 0x5b, 0x4d, 0x20, 0x21, 0x41}, true},
		{"test bit flip with click", ap, []byte{0x1b, 0x5b, 0x4d, 0x20, 0x21, 0x41}, true},
		{"test other mouse", ap, []byte{0x1b, 0x5b, 0x4d, 0x20, 0x21, 0x42}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configure(tt.ap)
			c.AP.Data = tt.data
			got := c.handleInput()
			if tt.want != got {
				t.Errorf("handleInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssign(t *testing.T) {
	s := newState()
	s.Exec("x=5")
	if s.variables["x"] != 5 {
		t.Fail()
	}
}

func TestDrawHistory(t *testing.T) {
	s := newState()
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
