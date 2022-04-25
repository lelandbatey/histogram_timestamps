<p align="center">
  <a href="https://github.com/lelandbatey/histogram_timestamps">
    <img src="https://user-images.githubusercontent.com/1964720/165003930-89c6f4ef-b481-4c38-9036-b898c5093ffe.png" width="400px"/>
  </a>
</p>


# Histogram Timestamps

`histogram_timestamps` takes a series of timestamps (no need to be sorted) in any format and shows you an interactive histogram of those timestamps.

## Use Case

You ever find yourself working with a dataset and you want to graph a
particular aspect of that data over time, but due to lots of data and/or weird
formats, there's no good tool to graph and visualize that data? This is the
problem `histogram_timestamps` is meant to solve! If you can get the timestamps
of anything, you can pipe them into `histogram_timestamps` to view trends in
that data over time. Here are example questions answerable with
`histogram_timestamps`:

- *"Hmmmm, we had a lot of broken data created last night; how can I see when it started vs when it grew out of control? There's no metric for this broken data because it's caused by a bug, but I can find all the broken rows in the DB with a query. If only I could graph that somehow..."*
	- Well with `histogram_timestamps`, you can! Run a `select created_time from table ...` command and pipe those timestamps straight into `histogram_timestamps` in order to instantly generate an interactive graph of all that data through time.

# Installation

On any OS, you can install `histogram_timestamps` by downloading the latest binary for your platform from [the latest release page](https://github.com/lelandbatey/histogram_timestamps/releases/tag/latest), extracting that binary, then placing that binary into your `$PATH` or equivalent.

## Linux (Bash)

```
# Installs in /usr/local/bin
curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/lelandbatey/histogram_timestamps/master/install.sh | sudo bash
```

## From source
You'll need Go 1.18+ and NPM installed to build from source, though the
generated binary is totally stand-alone and can be transfered to other systems
which do not have Go or NPM installed.
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
