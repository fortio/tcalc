package calculator

import (
	"errors"
	"slices"
)

func (s *State) Parse(tokens []string) (CalcNode, error) {
	node := s.parse(tokens, 0, nil)

	if node == nil {
		return CalcNode{}, errors.New("nil error")
	}
	return *node, nil
}

func (s *State) parse(tokens []string, index int, cur *CalcNode) *CalcNode {
	if index >= len(tokens) || len(tokens) == 0 {
		return cur
	}
	token := tokens[index]
	newNode := CalcNode{value: &token}
	if slices.Contains(Length1operatorsInfix, Operator(token[0])) && token[0] != '=' {
		newNode.left = cur
		return s.parse(tokens, index+1, &newNode)
	}
	switch token {
	case "=":
		if index == 0 || index == len(tokens)-1 {
			return nil
		}
		name := *cur.value
		newNode = *s.parse(tokens[index+1:], 0, nil)
		return &CalcNode{assignment: &assignment{name: name, right: newNode}}

	case "<<", ">>", "**":
		newNode.left = cur
		return s.parse(tokens, index+1, &newNode)
	case "(":
		rParenIndex := slices.Index(tokens[index:], ")")
		if rParenIndex == -1 {
			return nil
		}
		inner := innerParentheses(tokens[index+1:])
		node := s.parse(inner, 0, nil)
		if cur != nil {
			cur.right = node
			return s.parse(tokens[index+1+len(inner):], 0, cur)
		}
		return s.parse(tokens[index+1+len(inner):], 0, node)
	case ")":
		return s.parse(tokens, index+1, cur)
	default:
		if cur != nil {
			cur.right = &newNode
			return s.parse(tokens, index+1, cur)
		}
		return s.parse(tokens, index+1, &newNode)
	}
}

func innerParentheses(tokens []string) []string {
	score := 0
	for i, token := range tokens {
		switch token {
		case ")":
			if score == 0 {
				return tokens[:i]
			}
			score--
		case "(":
			score++
		}
	}
	return nil
}
