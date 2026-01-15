package calculator

import (
	"testing"
)

func TestExec(t *testing.T) {
	testCases := []struct {
		expression string
		expected   int64
		shouldFail bool
	}{
		{"1 + 3 + 2", 6, false},
		{"1 * (3 + 2)*(5+6)", 55, false},
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
		s := NewState()
		err := s.Exec(tc.expression)
		if tc.shouldFail {
			if err == nil {
				t.Errorf("Expected failure for expression: %s", tc.expression)
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error for expression %s: %v", tc.expression, err)
		} else if s.Ans != tc.expected {
			t.Errorf("For expression %s, expected %d but got %d", tc.expression, tc.expected, s.Ans)
		}
	}
}
