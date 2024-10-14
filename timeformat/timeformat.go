package timeformat

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
	"time"

	timefmt "github.com/itchyny/timefmt-go"
)

// ParseFunc will be a function that attempts to parse a string into a
// time.Time. A ParseFunc may try many different parse functions, or just one.
type ParseFunc func(string) (time.Time, error)

// FmtFunc is the opposite of ParseFunc; it transforms a time.Time into a
// string. FmtFuncs have no restrictions on behavior (e.g. they're allowed to
// return different results given the same input, if they desire).
type FmtFunc func(time.Time) (string, error)

// NewFuncs tries its best to give you funcs that parse your input and format your output the way
// you want based on the command-line flag inputs.
func NewFuncs(strptimefmt, gotimefmt string) (ParseFunc, FmtFunc, error) {
	if strptimefmt == "" && gotimefmt == "" {
		return parseUnixMillis, fmtUnixMillis, nil
	}

	if gotimefmt != "" {
		return makeParseGotime(gotimefmt), makeFmtGotime(gotimefmt), nil
	}
	if strptimefmt != "" {
		return makeParseStrptime(strptimefmt), makeFmtStrptime(strptimefmt), nil
	}
	return nil, nil, fmt.Errorf("logically you shouldn't be able to get this error; congratulations!")
}

func parseUnixMillis(s string) (time.Time, error) {
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(ts).UTC(), nil
}
func fmtUnixMillis(t time.Time) (string, error) {
	return fmt.Sprintf("%d", t.UnixMilli()), nil
}

func makeParseStrptime(sfmt string) ParseFunc {
	return func(s string) (time.Time, error) {
		return timefmt.Parse(s, sfmt)
	}
}

func makeFmtStrptime(sfmt string) FmtFunc {
	return func(t time.Time) (string, error) {
		return timefmt.Format(t, sfmt), nil
	}
}

func makeParseGotime(sfmt string) ParseFunc {
	return func(s string) (time.Time, error) {
		return time.Parse(sfmt, s)
	}
}

func makeFmtGotime(sfmt string) FmtFunc {
	return func(t time.Time) (string, error) {
		return t.Format(sfmt), nil
	}
}

// GuessStrptimeFormat checks if your timestamp can be parsed by any commonly
// used strptime formats. If the input line DOES match a commonly known strptime
// format, then the name and strptime-compatible formatting string (e.g.
// '%Y-%m-%d') is returned. Otherwise, an error is returned.
func GuessStrptimeFormat(line string) (string, string, error) {
	type tsfmt struct {
		Name string
		Fmt  string
	}
	potenials := []tsfmt{
		{"RFC3339", "%Y-%m-%dT%H:%M:%S.%f%z"},
		{"RFC3339", "%Y-%m-%dT%H:%M:%S%z"},
		{"YYYY-mm-dd HH:MM:SS.ms", "%Y-%m-%d %H:%M:%S.%f"},
		{"YYYY-mm-dd HH:MM:SS.ms TZ", "%Y-%m-%d %H:%M:%S.%f %z"},
		{"YYYY-mm-dd HH:MM:SS", "%Y-%m-%d %H:%M:%S"},
		{"YYYY-mm-dd HH:MM:SS TZ", "%Y-%m-%d %H:%M:%S %z"},
		{"YYYY-mm-dd", "%Y-%m-%d"},
	}
	for _, candidate := range potenials {
		_, err := timefmt.Parse(line, candidate.Fmt)
		if err == nil {
			return candidate.Name, candidate.Fmt, nil
		}
	}
	return "", "", fmt.Errorf("could not parse %q with any known strptime formats", line)
}

func GuessGoTimeFormat(line string) (string, string, error) {
	type tsfmt struct {
		Name string
		Fmt  string
	}
	potenials := []tsfmt{
		// Note that we do not need a timestamp specifier for RFC3339 with
		// micro/mili/nano seconds, just one with whole seconds. This is
		// because of a documented quirk of how the official Go time library
		// parses formats; to quote the documentation:
		//
		// > When parsing (only), the input may contain a fractional second field
		// > immediately after the seconds field, even if the layout does not
		// > signify its presence. In that case either a comma or a decimal point
		// > followed by a maximal series of digits is parsed as a fractional
		// > second.
		//
		// So basically, even if you don't include fractional seconds in the
		// timestamp specifier, Go will still parse the fractional second.
		{"RFC3339", "2006-01-02T15:04:05.999Z07:00"},
		{"RFC1123", time.RFC1123},
		{"RFC1123Z", time.RFC1123Z},
		{"RFC822", time.RFC822},
		{"RFC822Z", time.RFC822Z},
		{"YYYY-mm-dd HH:MM:SS.ms", "2006-01-02 15:04:05.999"},
		{"YYYY-mm-dd HH:MM:SS.ms TZ", "2006-01-02 15:04:05.999 -07:00"},
		{"YYYY-mm-dd", "2006-01-02"},
	}
	for _, candidate := range potenials {
		_, err := time.Parse(candidate.Fmt, line)
		if err == nil {
			return candidate.Name, candidate.Fmt, nil
		}
	}
	return "", "", fmt.Errorf("could not parse %q with any known gotime formats", line)
}

func GuessTimestampFormat(line string) string {
	hinttempl := `
HINT: It looks like the timestamp '{{.UnknownLine}}' is in a format named
'{{.FormatName}}' with a format specification of '{{.FormatContents}}'. To parse
all the incoming timestamps as format '{{.FormatName}}', provide the following option:

	{{.FlagName}} '{{.FormatContents}}'

`

	var tmpl = template.Must(template.New("FormatHintTemplate").Parse(hinttempl))
	type guessfmt struct {
		FlagName string
		FmtFunc  func(string) (string, string, error)
	}
	guessers := []guessfmt{
		{"--gotime-fmt", GuessGoTimeFormat},
		{"--strptime-fmt", GuessStrptimeFormat},
	}

	for _, g := range guessers {
		fmtname, format, err := g.FmtFunc(line)
		if err == nil {
			buf := &bytes.Buffer{}
			tmpl.Execute(buf, map[string]string{
				"UnknownLine":    line,
				"FormatName":     fmtname,
				"FormatContents": format,
				"FlagName":       g.FlagName,
			})
			return buf.String()
		}
	}
	return "HINT: Use the '--strptime-format' flag to indicate the format of the incoming timestamps\n\n"
}
