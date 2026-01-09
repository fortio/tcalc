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

func TestDraw(t *testing.T) {
	fmt.Println(decimalDisplayString(63))
	fmt.Println(hexDisplayString(63))
	fmt.Println(binaryDisplayString(63))
}
