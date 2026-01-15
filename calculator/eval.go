package calculator

import (
	"errors"
	"math"
	"slices"
	"strconv"
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
	if (*curNode.value)[0] == '\'' && (*curNode.value)[2] == '\'' {
		return int64((*curNode.value)[1]), nil
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
	num, err := strconv.ParseInt(*curNode.value, 0, 64)
	if err != nil {
		if *curNode.value == "_ans_" {
			return s.Ans, nil
		}
		return s.Variables[*curNode.value], nil
	}
	return num, nil
}
