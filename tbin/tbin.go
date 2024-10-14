package tbin

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

const TD_1_ms int64 = 1
const TD_1_sec int64 = 1000 * TD_1_ms
const TD_1_min int64 = 60 * TD_1_sec
const TD_1_hr int64 = 60 * TD_1_min
const TD_1_day int64 = 24 * TD_1_hr
const TD_1_week int64 = 7 * TD_1_day

var TIMEDELTA_ABBREVS map[string]string = map[string]string{
	"Y":            "Y", // year
	"y":            "Y",
	"W":            "W", // week
	"w":            "W",
	"D":            "D", // day
	"d":            "D",
	"days":         "D",
	"day":          "D",
	"hours":        "h",
	"hour":         "h",
	"hr":           "h",
	"h":            "h",
	"m":            "m",
	"minute":       "m",
	"min":          "m",
	"minutes":      "m",
	"t":            "m",
	"s":            "s",
	"seconds":      "s",
	"sec":          "s",
	"second":       "s",
	"ms":           "ms",
	"milliseconds": "ms",
	"millisecond":  "ms",
	"milli":        "ms",
	"millis":       "ms",
	"l":            "ms",
}

var ABBREV_TO_DELT map[string]int64 = map[string]int64{
	"Y":  365 * TD_1_day,
	"W":  TD_1_week,
	"D":  TD_1_day,
	"h":  TD_1_hr,
	"m":  TD_1_min,
	"s":  TD_1_sec,
	"ms": TD_1_ms,
}

var ABBREV_LARGE_TO_SMALL []string = []string{"Y", "W", "D", "h", "m", "s", "ms"}

// Maps the time abbreviations originally taken from Pandas onto the time
// abbreviations needed by ChartJS:
// https://www.chartjs.org/docs/3.0.2/axes/cartesian/time.html#time-units
var ABBREV_TO_CHARTJS_UNIT map[string]string = map[string]string{
	"Y":  "year",
	"W":  "week",
	"D":  "day",
	"h":  "hour",
	"m":  "minute",
	"s":  "second",
	"ms": "millisecond",
}

// BinTimestamp takes a timestamp in epoch_ms format and returns that same
// timestamp floor-ed down to the nearest 'frequency' you provided, effectively
// giving you the "bin" where this timestamp belongs in a histogram with bins
// of size 'frequency'. If 'frequency' does not stand for a known bin-size,
// then an error is returned.
func BinTimestamp(ts int64, spec string) (int64, error) {
	mult, delt, err := ParseSpec(spec)
	if err != nil {
		return 0, err
	}
	d := delt * mult
	return (ts / d) * d, nil
}

func BinTimestamps(tss []int64, spec string) (map[int64]int64, error) {
	hist := map[int64]int64{}
	for _, ts := range tss {
		bin, err := BinTimestamp(ts, spec)
		if err != nil {
			return nil, err
		}
		if _, ok := hist[bin]; !ok {
			hist[bin] = 0
		}
		hist[bin] = hist[bin] + 1
	}
	sort.SliceStable(tss, func(i, j int) bool { return tss[i] < tss[j] })
	mult, basedelt, err := ParseSpec(spec)
	if err != nil {
		return nil, err
	}
	delt := mult * basedelt
	minbin, _ := BinTimestamp(tss[0], spec)
	maxbin, _ := BinTimestamp(tss[len(tss)-1], spec)
	cur := minbin
	for cur < maxbin {
		cur += delt
		cb, _ := BinTimestamp(cur, spec)
		if _, ok := hist[cb]; !ok {
			hist[cb] = 0
		}
	}
	return hist, nil
}

type ChartJSDatapoint struct {
	X interface{} `json:"x"`
	Y interface{} `json:"y"`
}
type ChartJSCtx struct {
	Unit string             `json:"unit"`
	Data []ChartJSDatapoint `json:"data"`
}

func FormatBinDataForChartJS(bins map[int64]int64) (ChartJSCtx, error) {
	ctx := ChartJSCtx{}
	keys := []int64{}
	for k := range bins {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })
	_, jsunit := EstimateBinSize(keys)
	for _, k := range keys {
		v := bins[k]
		ctx.Data = append(ctx.Data, ChartJSDatapoint{X: k, Y: v})
	}
	ctx.Unit = jsunit
	return ctx, nil
}

// EstimateBinSize returns two abbreviations for duration. The first is an
// abbreviation appropriate to get a timedelta duration from ABBREV_TO_DELT,
// while the second is a ChartJS compatible abbreviation.
// Note that this isn't a very good way to estimate the bin size, as it doesn't
// account for the density of data over the minimum time span. For example, if
// the data spans a period of 2 days, this'll bin that data into just two bins,
// one for each day, because the 48 hour duration of the data divides nicely
// into day long segments.
func EstimateBinSize(tss []int64) (string, string) {
	sort.SliceStable(tss, func(i, j int) bool { return tss[i] < tss[j] })
	var dur int64 = tss[len(tss)-1] - tss[0]
	unit := "ms"
	for _, abrv := range ABBREV_LARGE_TO_SMALL {
		var delt int64 = ABBREV_TO_DELT[abrv]
		// The smallest unit which still divides the duration of the bins into
		// a whole integer
		if (dur / delt) < 1 {
			continue
		}
		unit = abrv
		break
	}
	jsunit := ABBREV_TO_CHARTJS_UNIT[unit]
	return unit, jsunit
}

func ParseSpec(unit string) (mult int64, delt int64, err error) {
	rs := []rune(unit)
	var numbers []rune
	var letters []rune
	for _, r := range rs {
		if unicode.IsNumber(r) {
			numbers = append(numbers, r)
		} else {
			letters = append(letters, r)
		}
	}
	if len(numbers) == 0 {
		mult = 1
	} else {
		mult, err = strconv.ParseInt(string(numbers), 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	abbrev, ok := TIMEDELTA_ABBREVS[string(letters)]
	if !ok {
		return 0, 0, fmt.Errorf("no timedelta configured for abbreviation of %q", string(letters))
	}
	delt, ok = ABBREV_TO_DELT[abbrev]
	if !ok {
		return 0, 0, fmt.Errorf("no timedelta configured for frequency of %q leading to abbrev %q", string(letters), abbrev)
	}
	return mult, delt, nil
}
