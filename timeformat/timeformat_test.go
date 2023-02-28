package timeformat_test

import (
	"fmt"
	"testing"

	"github.com/lelandbatey/histogram_timestamps/timeformat"

	"github.com/stretchr/testify/require"
)

func TestGuessGoTimeFormat(t *testing.T) {
	type tcase struct {
		TS           string
		ExpectedFmt  string
		ExpectedName string
	}

	for idx, tc := range []tcase{
		// Due to some oddness in the Golang time.Parse(), fractional seconds
		// are always parsed whether specified or not.
		{"2023-02-27T15:38:17.847773933-08:00", "2006-01-02T15:04:05.999Z07:00", "RFC3339"},
		{"2023-02-27T15:38:17-08:00", "2006-01-02T15:04:05.999Z07:00", "RFC3339"},
		{"Tue, 28 Feb 2023 12:11:26 PST", "Mon, 02 Jan 2006 15:04:05 MST", "RFC1123"},
		{"Tue, 28 Feb 2023 12:11:26 -0800", "Mon, 02 Jan 2006 15:04:05 -0700", "RFC1123Z"},
		{"28 Feb 23 12:16 PST", "02 Jan 06 15:04 MST", "RFC822"},
		{"28 Feb 23 12:16 -0800", "02 Jan 06 15:04 -0700", "RFC822Z"},
		{"2023-02-28 12:16:42.002", "2006-01-02 15:04:05.999", "YYYY-mm-dd HH:MM:SS.ms"},
		{"2023-02-28 12:16:42.002 -07:00", "2006-01-02 15:04:05.999 -07:00", "YYYY-mm-dd HH:MM:SS.ms TZ"},
		{"2023-02-28", "2006-01-02", "YYYY-mm-dd"},
	} {
		t.Run(fmt.Sprintf("GuessGoTimeFormat case #%d", idx), func(t *testing.T) {
			name, format, err := timeformat.GuessGoTimeFormat(tc.TS)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedFmt, format)
			require.Equal(t, tc.ExpectedName, name)
		})
	}
}

func TestGuessStrptimeFormat(t *testing.T) {
	type tcase struct {
		TS           string
		ExpectedFmt  string
		ExpectedName string
	}

	for idx, tc := range []tcase{
		{"2023-02-28T12:24:13.926780-0800", "%Y-%m-%dT%H:%M:%S.%f%z", "RFC3339"},
		{"2023-02-28T12:24:13-0800", "%Y-%m-%dT%H:%M:%S%z", "RFC3339"},
		{"2023-02-28 12:24:13.926841", "%Y-%m-%d %H:%M:%S.%f", "YYYY-mm-dd HH:MM:SS.ms"},
		{"2023-02-28 12:24:13.926845 -0800", "%Y-%m-%d %H:%M:%S.%f %z", "YYYY-mm-dd HH:MM:SS.ms TZ"},
		{"2023-02-28 12:24:13", "%Y-%m-%d %H:%M:%S", "YYYY-mm-dd HH:MM:SS"},
		{"2023-02-28 12:24:13 -0800", "%Y-%m-%d %H:%M:%S %z", "YYYY-mm-dd HH:MM:SS TZ"},
		{"2023-02-28", "%Y-%m-%d", "YYYY-mm-dd"},
	} {
		t.Run(fmt.Sprintf("GuessGoTimeFormat case #%d", idx), func(t *testing.T) {
			name, format, err := timeformat.GuessStrptimeFormat(tc.TS)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedFmt, format)
			require.Equal(t, tc.ExpectedName, name)
		})
	}
}
