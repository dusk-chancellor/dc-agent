package utils

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrMismatchParenthesis = errors.New("mismatched parenthesis")
	ErrParseOperand = errors.New("failed to parse operand")
)

type Utils interface {
	ToPostfix(infix string) (postfix string, err error)
	Evaluate(queue string) (result float64, err error)
}

type utils struct {
	l *zap.Logger
}

func New(l *zap.Logger) Utils {
	return &utils{l: l}
}
