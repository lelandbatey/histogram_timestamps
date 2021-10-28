# Histogram Timestamps

You ever find yourself working with a dataset and you want to graph a
particular aspect of that data over time, but because you're dealing with a lot
of data and in a weird format, there's no good tool to graph and visualize that
data? This is the problem `histogram_timestamps` is meant to solve! If you can
get the timestamps of anything, you can pipe them into `histogram_timestamps`
to view trends in that data over time. Here are example questions answerable
with `histogram_timestamps`:

- *"Hmmmm, we had a lot of broken data created last night; how can I see when it started vs when it grew out of control? There's no metric for this broken data because it's caused by a bug, but I can find all the broken rows in the DB with a query. If only I could graph that somehow..."*
	- Well with `histogram_timestamps`, you can! Run a `select created_time from table ...` command and pipe those timestamps straight into `histogram_timestamps` in order to instantly generate an interactive graph of all that data through time.

# Installation

```
git clone https://github.com/lelandbatey/histogram_timestamps.git
cd histogram_timestamps
make install # Does require you have NPM installed to build the JS portion, but once pre-built the binary is totaly self-contained
```

# Usage

As an example of usage, you can have `histogram_timestamps` generate some fake data which you feed back into `histogram_timestamps`. Examples:

```
./histogram_timestamps --generate-fake-data | ./histogram_timestamps
./histogram_timestamps --generate-fake-data | ./histogram_timestamps --unit minute
./histogram_timestamps --generate-fake-data | ./histogram_timestamps --unit hour
```
