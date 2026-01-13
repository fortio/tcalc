package calculator

type CalcNode struct {
	left       *CalcNode
	right      *CalcNode
	value      *string
	assignment *assignment
}

type State struct {
	Variables map[string]int64
	Ans       int64
	Err       error
}

func NewState() *State {
	return &State{
		Variables: make(map[string]int64),
	}
}

//go:generate stringer -type=Operator
type (
	Operator           rune
	DoubleRuneOperator string
)

const (
	SUM  Operator = '+'
	SUB  Operator = '-'
	PROD Operator = '*'
	DIV  Operator = '/'
	XOR  Operator = '^'
	OR   Operator = '|'
	NOT  Operator = '~'
	AND  Operator = '&'
	MOD  Operator = '%'

	ASSIGN Operator = '='

	LPAREN Operator = '('
	RPAREN Operator = ')'
	// two rune operators

	LEFTSHIFT  DoubleRuneOperator = "<<"
	RIGHTSHIFT DoubleRuneOperator = ">>"
	EXP        DoubleRuneOperator = "**"
)

var Length1operatorsInfix = []Operator{
	SUM, SUB, PROD, DIV, XOR, AND, MOD, ASSIGN, OR,
}

var Length1operatorsPrefix = []Operator{
	NOT,
}

var Length2operators = []DoubleRuneOperator{
	LEFTSHIFT, RIGHTSHIFT,
}

func (s *State) Exec(input string) error {
	tokens, err := s.Tokenize(input)
	if err != nil {
		s.Err = err
		return err
	}
	node, err := s.Parse(tokens)
	if err != nil {
		s.Err = err
		return err
	}
	value, err := s.Eval(node)
	s.Err = err
	if err != nil {
		return err
	}

	s.Ans = value
	return nil
}

type assignment struct {
	name  string
	right CalcNode
}
