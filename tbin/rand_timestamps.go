package tbin

import (
	"fmt"
	"math/rand"
)

/*
# Start time:
#     1572347470840
#     Tuesday, October 29, 2019 11:11:10.840 AM (UTC)
def gen_rand_dates(count=100, seed=1, start_ts=1572347470840, stop_ts=None, dist=None):
    if dist is None:
        dist = 'rand'
    if stop_ts is None:
        stop_ts = start_ts + (12 * _1_hr)
    rnd = random.Random(seed)
    range_ = stop_ts - start_ts
    getrand = lambda : int(rnd.random() * range_) + start_ts
    if dist == 'rand':
        getrand = lambda : rnd.randint(start_ts, stop_ts)
    elif dist == 'normal':
        getrand = lambda : min(max(int(rnd.normalvariate(0.5, 0.15) * range_), 0), range_) + start_ts

    for _ in range(count):
        yield getrand()
*/

func clamp(v, lo, hi int64) int64 {
	rv := v
	if v > hi {
		rv = hi
	}
	if v < lo {
		rv = lo
	}
	return rv
}

func GenRandomTimestamps(count int, seed int64, start_ts int64, stop_ts int64, dist string) ([]int64, error) {
	rnd := rand.New(rand.NewSource(seed))
	range_ := stop_ts - start_ts
	var getrand func() int64
	if dist == "random" {
		getrand = func() int64 { return clamp(int64(rnd.Float64()*float64(range_))+start_ts, start_ts, stop_ts) }
	} else if dist == "normal" {
		getrand = func() int64 {
			return clamp(int64(rnd.NormFloat64()*float64(range_/8))+start_ts+(range_/2), start_ts, stop_ts)
		}
	} else {
		return nil, fmt.Errorf("distribution %q specified in dist is not valid", dist)
	}

	tss := []int64{}
	for i := 0; i < count; i++ {
		tss = append(tss, getrand())
	}
	return tss, nil
}

func SimpleRandomTimestamps(count int, hours_duration int) ([]int64, error) {
	// Start time:
	// 1572347470840
	// Tuesday, October 29, 2019 11:11:10.840 AM (UTC)
	var start_ts int64 = 1572347470840
	stop_ts := start_ts + (int64(hours_duration) * TD_1_hr)
	return GenRandomTimestamps(count, 1, start_ts, stop_ts, "normal")
}
