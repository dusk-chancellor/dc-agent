package utils

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/dusk-chancellor/dc-agent/stack"
	"go.uber.org/zap"
)

// evaluates postfix (RPN) expression
func (u *utils) Evaluate(q string) (float64, error) {
    // create new stack of operands
    o := stack.NewStack()
    // split each element
    chars := strings.Split(q, " ")
    u.l.Debug("", zap.Any("chars", chars))

    for _, ch := range chars {
        if ch == "" {
            break
        // if number
        } else if !isOperator(ch) &&
        !isFunction(ch) {
            operand, err := strconv.ParseFloat(ch, 64)
            if err != nil {
                u.l.Error("failed to parse operand", zap.Error(err), zap.String("operand", ch))
                return 0, ErrParseOperand
            }

            u.l.Debug("operand pushed to stack", zap.Float64("op", operand))
            o.Push(operand)
        // if operator or function
        } else {
            // top element as 2nd operand
            op2, err := o.Peek()
            if err != nil {
                u.l.Error("failed to peek",  zap.Any("stack", o), zap.Error(err))
                return 0, err
            }

            u.l.Debug("operand №2 poped from stack", zap.Any("op", op2))
            o2 := op2.(float64)
            o.Pop() // pop it up
            
            // next top element as 1st operand
            op1, err := o.Peek()
            if err != nil {
                u.l.Error("failed to peek",  zap.Any("stack", o), zap.Error(err))
                return 0, err
            }

            u.l.Debug("operand №1 poped from stack", zap.Any("op", op1))
            o1 := op1.(float64)
            o.Pop() // pop it up

            c, err := calculate(ch, o1, o2)
            if err != nil {
                u.l.Error("failed to calculate",  zap.String("operator", ch), zap.Error(err))
                return 0, err
            }

            u.l.Debug("result pushed to stack", zap.Float64("c", c))
            o.Push(c)
        }
    }

    last, err := o.Peek()
    if err != nil {
        u.l.Error("failed to peek",  zap.Any("stack", o), zap.Error(err))
        return 0, err
    }

    res := last.(float64)

    return res, nil
}

// calculates operation between two operands where:
// o  - operator;
// o1 - operand №1;
// o2 - operand №2
func calculate(o string, o1, o2 float64) (float64, error) {
    switch o {
    case "+":
        return o1 + o2, nil
    case "-":
        return o1 - o2, nil
    case "*":
        return o1 * o2, nil
    case "/", ":":
        return o1 / o2, nil
    case "^":
        return math.Pow(o1, o2), nil
    default:
        return 0, errors.ErrUnsupported
    }
}
