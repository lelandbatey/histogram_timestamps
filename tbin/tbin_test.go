package tbin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSpec(t *testing.T) {
	type tcase struct {
		Spec    string
		ExpMult int64
		ExpDelt int64
		ExpErr  error
	}

	for idx, test := range []tcase{
		{
			Spec:    "30D",
			ExpMult: 30,
			ExpDelt: TD_1_day,
			ExpErr:  nil,
		},
		{
			Spec:    "30m",
			ExpMult: 30,
			ExpDelt: TD_1_min,
			ExpErr:  nil,
		},
		{
			Spec:    "5Y",
			ExpMult: 5,
			ExpDelt: TD_1_day * 365,
			ExpErr:  nil,
		},
		{
			Spec:    "1W",
			ExpMult: 1,
			ExpDelt: TD_1_week,
			ExpErr:  nil,
		},
		{
			Spec:    "m",
			ExpMult: 1,
			ExpDelt: TD_1_min,
			ExpErr:  nil,
		},
		{
			Spec:    "1h",
			ExpMult: 1,
			ExpDelt: TD_1_hr,
			ExpErr:  nil,
		},
		{
			Spec:    "1",
			ExpMult: 0,
			ExpDelt: 0,
			ExpErr:  fmt.Errorf("no timedelta configured for abbreviation of \"\""),
		},
	} {
		mult, delt, err := ParseSpec(test.Spec)
		require.Equal(t, test.ExpErr, err, "for test #%d", idx)
		require.Equal(t, test.ExpMult, mult, "for test #%d", idx)
		require.Equal(t, test.ExpDelt, delt, "for test #%d", idx)
	}
}
