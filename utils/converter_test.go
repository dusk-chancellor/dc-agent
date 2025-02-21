package utils

import (
	"testing"

	"github.com/dusk-chancellor/dc-agent/zaplog"
)

func TestToPostfix(t *testing.T) {
	cases := []struct{
		name string
		in	 string
		err	 error
	}{
		{
			name: "successful normal",
			in:   "2+2*2",
			err:  nil,
		},
		{
			name: "successful parenthesis",
			in:	  "(2+2)*2",
			err:  nil,
		},
		{
			name: "successful tough",
			in:   "18/12*16+(11+14)-25",
			err:  nil,
		},
		{
			name: "mismatch parenthesis",
			in:	  "(2+2*2(",
			err:  ErrMismatchParenthesis,
		},
	}

	log := zaplog.New()
	u := New(log)

	for _, tt := range cases {
		q, err := u.ToPostfix(tt.in)
		if err != tt.err {
			t.Errorf("expected: %v, got: %v", tt.err, err)
		}

		t.Logf("queue: %v, error: %v", q, err)
	}
}
