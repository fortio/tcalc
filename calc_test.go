package main

import (
	"fmt"
	"testing"
)

func BenchmarkExec(b *testing.B) {
	s := state{
		variables: make(map[string]int64),
	}
	b.Run("operators", func(b *testing.B) {
		err := s.Exec("1 + 3 + 2")
		if err != nil {
			b.Fail()
		}
		if s.ans != 6 {
			b.Fail()
		}
		err = s.Exec("1 * (3 + 2)")
		fmt.Println("execed")
		if err != nil || s.ans != 5 {
			b.Fail()
		}
		fmt.Println(s.ans)
		tokens := s.Tokenize("(2 * (3 + 2) - 1)+ 1 / 1")
		fmt.Println(tokens)
		err = s.Exec("(2 * (3 + 2) - 1)+ 1 / 1")
		fmt.Println(s.ans)
		if err != nil || s.ans != 10 {
			b.Fail()
		}
		err = s.Exec("1 << 5")

		fmt.Println(s.ans)
		if err != nil || s.ans != 1<<5 {
			b.Fail()
		}
		tokens = s.Tokenize("2**5")
		fmt.Println(tokens)
		tokens = s.Tokenize("+-2")
		fmt.Println(tokens)
	})
	b.Run("assignment", func(b *testing.B) {
		err := s.Exec("x=1")
		if err != nil || s.variables["x"] != 1 {
			b.Fail()
		}
	})
}

func TestExec(t *testing.T) {
	s := state{
		variables: make(map[string]int64),
	}
	t.Run("operators", func(b *testing.T) {
		err := s.Exec("1 + 3 + 2")
		if err != nil {
			b.Fail()
		}
		if s.ans != 6 {
			b.Fail()
		}
		err = s.Exec("1 * (3 + 2)")
		fmt.Println("execed")
		if err != nil || s.ans != 5 {
			b.Fail()
		}
		fmt.Println(s.ans)
		tokens := s.Tokenize("(2 * (3 + 2) - 1)+ 1 / 1")
		fmt.Println(tokens)
		err = s.Exec("(2 * (3 + 2) - 1)+ 1 / 1")
		fmt.Println(s.ans)
		if err != nil || s.ans != 10 {
			b.Fail()
		}
		err = s.Exec("1 << 5")

		fmt.Println(s.ans)
		if err != nil || s.ans != 1<<5 {
			b.Fail()
		}
		tokens = s.Tokenize("2**5")
		fmt.Println(tokens)
		err = s.Exec("0+-1")
		if err == nil {
			t.Fail()
		}
		fmt.Println(tokens)
	})
	t.Run("assignment", func(b *testing.T) {
		err := s.Exec("x=1")
		if err != nil || s.variables["x"] != 1 {
			b.Fail()
		}
	})

	t.Run("bitwise operators", func(t *testing.T) {
		err := s.Exec("1&2")
		if err != nil || s.ans != 0 {
			t.Fail()
		}
		err = s.Exec("1|2")
		if err != nil || s.ans != 3 {
			t.Fail()
		}
		err = s.Exec("1^2")
		if err != nil || s.ans != 3 {
			t.Fail()
		}
		err = s.Exec("1^3")
		if err != nil || s.ans != 2 {
			t.Fail()
		}
		err = s.Exec("~1")
		if err != nil || s.ans != -2 {
			t.Fail()
		}
		err = s.Exec("2>>1")
		if err != nil || s.ans != 1 {
			t.Fail()
		}
	})

	t.Run("display strings", func(t *testing.T) {
		if decimalDisplayString(64) != "Decimal: 64\n" {
			t.Fail()
		}
		if hexDisplayString(64) != "Hex: 40\n" {
			t.Fail()
		}
	})
}
