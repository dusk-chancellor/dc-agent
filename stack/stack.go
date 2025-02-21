package stack

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
)

// Stack data structure implementation

var (
    ErrStackEmpty = errors.New("stack is empty")
)

type Stack struct {
    l *list.List
}

func NewStack() *Stack {
    return &Stack{l: list.New()}
}

func (s *Stack) Push(v any) {
    if v == nil {
        return
    }

    s.l.PushBack(v)
}

func (s *Stack) Pop() (any, error) {
    if s.IsEmpty() {
        return nil, ErrStackEmpty
    }

    tail := s.l.Back()
    val := tail.Value
    s.l.Remove(tail)

    return val, nil
}

func (s *Stack) Peek() (any, error) {
    if s.IsEmpty() {
        return nil, ErrStackEmpty
    }

    tail := s.l.Back()
    val := tail.Value
    return val, nil
}


func (s *Stack) IsEmpty() bool {
    if s.l.Len() == 0 {
        return true
    }

    return false
}

func (s *Stack) String() (res string) {
    if s.IsEmpty() {
        return ""
    }

    var sb strings.Builder

    for e := s.l.Front(); e != nil; e = e.Next() {
        if sb.Len() > 0 {
            sb.WriteString(" ")
        }

        v := fmt.Sprintf("%v", e.Value)

        sb.WriteString(v)
    }

    return sb.String()
}
