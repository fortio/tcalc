package calculator

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"unicode/utf8"
)

func (s *State) Eval(curNode CalcNode) (int64, error) { //nolint:funlen,gocyclo // evaluation will be hairy
	if curNode.assignment != nil {
		num, err := s.Eval(curNode.assignment.right)
		if err != nil {
			return -1, err
		}
		s.Variables[curNode.assignment.name] = num
		return num, nil
	}
	if curNode.value == nil {
		return -1, errors.New("bad value")
	}
	if (*curNode.value)[0] == '\'' && (*curNode.value)[len(*curNode.value)-1] == '\'' {
		r, _ := utf8.DecodeRuneInString((*curNode.value)[1 : len(*curNode.value)-1])
		return int64(r), nil
	}
	if *curNode.value == "-" && (curNode.left == nil || curNode.left.value == nil) {
		num, err := s.Eval(*curNode.right)
		if err != nil {
			return -1, err
		}
		return -1 * num, nil
	}
	if slices.Contains(Length1operatorsInfix, Operator((*curNode.value)[0])) {
		l, err := s.Eval(*curNode.left)
		if err != nil {
			return 0, err
		}
		if curNode.right == nil {
			return 0, errors.New("invalid operator")
		}

		r, err := s.Eval(*curNode.right)
		if err != nil {
			return 0, err
		}
		switch *curNode.value {
		case "+":
			return l + r, nil
		case "-":
			return l - r, nil
		case "*":
			return l * r, nil
		case "/":
			return l / r, nil
		case "&":
			return l & r, nil
		case "^":
			return l ^ r, nil
		case "|":
			return l | r, nil
		case "%":
			return l % r, nil
		case "]":
			return int64(pext(uint64(l), uint64(r))), nil //nolint:gosec // we just want to display the unsigned representation of our number
		case "[":
			return int64(pdep(uint64(l), uint64(r))), nil //nolint:gosec // we just want to display the unsigned representation of our number
		default:
			return -1, errors.New("invalid operator")
		}
	}
	if slices.Contains(Length1operatorsPrefix, Operator((*curNode.value)[0])) {
		num, err := s.Eval(*curNode.right)
		if err != nil {
			return -1, err
		}
		switch *curNode.value {
		case "~":
			return ^num, nil
		default:
			return -1, errors.New("bad prefix operator")
		}
	}
	if slices.Contains(Length2operators, DoubleRuneOperator(*curNode.value)) {
		l, err := s.Eval(*curNode.left)
		if err != nil {
			return 0, err
		}
		r, err := s.Eval(*curNode.right)
		if err != nil {
			return 0, err
		}
		switch *curNode.value {
		case "<<":
			return l << r, nil
		case ">>":
			return l >> r, nil
		case "**":
			f := math.Pow(float64(l), float64(r))
			return int64(f), nil
		default:
			return -1, errors.New("bad double rune operator")
		}
	}
	num, err := strconv.ParseUint(*curNode.value, 0, 64)
	if err != nil {
		if *curNode.value == "ans" {
			return s.Ans, nil
		}
		return s.Variables[*curNode.value], nil
	}
	return int64(num), nil //nolint:gosec // must cast so that we don't discard sign bit
}

func pdep(num uint64, mask uint64) uint64 {
	var result uint64
	numIndex := 0
	for i := range 64 {
		if (mask & (1 << i)) != 0 {
			if (num & (1 << numIndex)) != 0 {
				result |= (1 << i)
			}
			numIndex++
		}
	}
	return result
}

func pext(num uint64, mask uint64) uint64 {
	var result uint64
	numIndex := 0
	for i := range 64 {
		if (mask & (1 << i)) != 0 {
			if (num & (1 << i)) != 0 {
				result |= (1 << numIndex)
			}
			numIndex++
		}
	}
	return result
}
