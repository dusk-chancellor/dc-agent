package utils

import (
	"testing"

	"github.com/dusk-chancellor/dc-agent/zaplog"
)

func TestEvaluate(t *testing.T) {
	cases := []struct{
		name 	   string
		expression string
		wantResult float64
		err 	   error
	} {
		{
			name: "successful simple",
			expression: "2 2 2 * + ",
			wantResult: float64(6),
			err: nil,
		},
		{
			name: "succesful tough",
			expression: "18 12 / 16 * 11 14 + + 25 -",
			wantResult: float64(24),
			err: nil,
		},
	}

	log := zaplog.New()
	u := New(log)

	for _, tt := range cases {
		res, err := u.Evaluate(tt.expression)
		if err != tt.err {
			t.Errorf("expected: %v, got: %v", tt.err, err)
		} else if res != tt.wantResult {
			t.Errorf("expected: %f, got: %f", tt.wantResult, res)
		}

		t.Logf("res: %f, err: %v", res, err)
	}
}
