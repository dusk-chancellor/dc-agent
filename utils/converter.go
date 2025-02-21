package utils

import (
	"errors"

	"github.com/dusk-chancellor/dc-agent/stack"
	"go.uber.org/zap"
)

// converts basic expressions (infix notation)
// to reverse polish notation (postfix notation)
// using Dijkstra's Shunting-yard algorithm:
// https://en.wikipedia.org/wiki/Shunting_yard_algorithm#The_algorithm_in_detail
func (u *utils) ToPostfix(infix string) (string, error) {
	// create new stack of operators and queue for output
	o, q := stack.NewStack(), stack.NewStack()
	// for collecting digits
	var digits string

	// iterate over infix string
	for _, r := range infix {
		// stringify
		s := string(r)

		switch {
		case isOperand(s): // if 0-9
			u.l.Debug("operand digit added to digits", zap.String("s", s))
			digits += s // collect digit

		case isFunction(s):
			u.l.Debug("digits pushed to queue", zap.String("digits", digits))
			q.Push(digits)
			digits = ""

			u.l.Debug("function pushed to operators stack", zap.String("s", s))
			o.Push(s)

		case isOperator(s):
			u.l.Debug("digits pushed to queue", zap.String("digits", digits))
			q.Push(digits)
			digits = ""

			// retrieving top element
			var top string

			peek, err := o.Peek()
			if err != nil {
				u.l.Debug("no peek", zap.Any("peek", peek))
				top = "0"
			} else {
				top = peek.(string)
			}
			// start popping operators
			for !o.IsEmpty() && // stack is not empty
			!isLeftParenthesis(top) && // top elemtent is not left parenthesis
			(precedence(top) > precedence(s)) || // precedence of top element higher than current element
			(precedence(top) == precedence(s) && isLeftAssociative(s)) { // precedence is equal AND left associative
				pop, err := o.Pop()
				if err != nil {
					u.l.Debug("no pop", zap.Any("pop", pop))
					break
				}

				v := pop.(string)

				u.l.Debug("operator pushed to queue", zap.String("v", v))
				q.Push(v)
			}

			u.l.Debug("operator pushed to operators stack", zap.String("s", s))
			o.Push(s)

		case isLeftParenthesis(s):
			u.l.Debug("digits pushed to queue", zap.String("digits", digits))
			q.Push(digits)
			digits = ""

			u.l.Debug("left parenthesis pushed to operators stack", zap.String("s", s))
			o.Push(s)

		case isRigthParenthesis(s):
			u.l.Debug("digits pushed to queue", zap.String("digits", digits))
			q.Push(digits)
			digits = ""

			for {
				// retrieving top element
				var top string
			
				peek, err := o.Peek()
				if err != nil {
					u.l.Debug("no peek", zap.Any("peek", peek))
					break
				}
			
				top = peek.(string)

				if isLeftParenthesis(top) {
					u.l.Debug("left parenthesis poped from operators stack", zap.String("top", top))
					o.Pop()
					break
				}
				
				pop, _ := o.Pop()

				v := pop.(string)

				u.l.Debug("element pushed to queue", zap.String("v", v))
				q.Push(v)
			}

		default:
			u.l.Debug("digits pushed to queue", zap.String("digits", digits))
			q.Push(digits)
			digits = ""

			u.l.Debug("handling unsupported operation", zap.String("s", s))
			return "", errors.ErrUnsupported
		}
	}
	// push last number
	q.Push(digits)

	// pop up the rest of operators
	for !o.IsEmpty() {
		pop, _ := o.Pop()

		e := pop.(string)

		if isLeftParenthesis(e) || isRigthParenthesis(e) {
			u.l.Debug("handling mismatch parenthesis error")
			return "", ErrMismatchParenthesis
		}

		u.l.Debug("element pushed to queue", zap.String("e", e))
		q.Push(e)
	}

	return q.String(), nil
}

// defines operation precedence
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/", ":":
		return 2
	case "^": // left associative
		return 3
	default:
		return 0
	}
}

// check if operation is left associative
func isLeftAssociative(op string) bool {
	if op == "^" {
		return false
	}

	return true
}

// check if element is a number
func isOperand(s string) bool {
	if (s >= "0" && s <= "9") || s == "." || s == "," {
		return true
	}

	return false
}

// check if element is a function
func isFunction(s string) bool {
	// for now only `exp` is supported
	if s == "^" {
		return true
	}

	return false
}

// check if element is an operator
func isOperator(s string) bool {
	switch s {
	case "+", "-", "*", "/", ":":
		return true
	default:
		return false
	}
}

// check if element is left parenthesis
func isLeftParenthesis(s string) bool {
	switch s {
	case "(", "{", "[":
		return true
	default:
		return false
	}
}

// check if element is right parenthesis
func isRigthParenthesis(s string) bool {
	switch s {
	case ")", "}", "]":
		return true
	default:
		return false
	}
}
