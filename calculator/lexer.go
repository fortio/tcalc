package calculator

import (
	"errors"
	"slices"
	"strings"
)

func (s *State) Tokenize(input string) ([]string, error) {
	if strings.Count(input, "=") > 1 {
		return nil, errors.New("invalid double assignment")
	}
	tokens := make([]string, 0, len(input))
	cur := ""
	for _, char := range input {
		numTokens := len(tokens)
		if numTokens > 0 && tokens[numTokens-1] == "*" && char == '*' {
			tokens[numTokens-1] = "**"
			continue
		}
		if char == '(' ||
			char == ')' ||
			slices.Contains(Length1operatorsInfix, Operator(char)) ||
			slices.Contains(Length1operatorsPrefix, Operator(char)) {
			if len(cur) > 0 {
				tokens = append(tokens, cur)
				cur = ""
			}
			tokens = append(tokens, string(char))
			continue
		}
		switch char {
		case ' ', '\r', '\n':
			if len(cur) > 0 {
				tokens = append(tokens, cur)
				cur = ""
			}
		case '>', '<', '*':
			if cur == string(char) {
				tokens = append(tokens, cur+string(char))
				cur = ""
				continue
			}
			if len(cur) > 0 {
				tokens = append(tokens, cur)
			}
			cur = string(char)
		default:
			cur += string(char)
			if slices.Contains(Length2operators, DoubleRuneOperator(cur)) {
				tokens = append(tokens, cur)
				cur = ""
			}
		}
	}
	tokens = tokens[:len(tokens):len(tokens)]
	if cur != "" {
		tokens = append(tokens, cur)
	}
	return tokens, nil
}
