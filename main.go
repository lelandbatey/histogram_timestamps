package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	timefmt "github.com/itchyny/timefmt-go"
	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/pflag"

	"github.com/lelandbatey/histogram_timestamps/tbin"
	"github.com/lelandbatey/histogram_timestamps/timeformat"
)

var (
	outputpath   = pflag.StringP("output-path", "o", "./", "Path to the directory to write out the HTML file visualizing the timeseries data")
	title        = pflag.StringP("title", "t", "Timeseries data", "Title of the generated HTML page")
	generateData = pflag.BoolP("generate-fake-data", "g", false, "If provided, all the program will do is generate a bunch of fake timestamps and print them on stdout. Useful as a way to feed known input to another histogram_timestamps")
	unit         = pflag.StringP("unit", "u", "auto", "The duration of each 'bin' to group timestamps into: https://pandas.pydata.org/pandas-docs/stable/user_guide/timeseries.html#offset-aliases")
	strptimefmt  = pflag.StringP("strptime-fmt", "f", "", "A strptime-compatible date format specifier. Use if your data isn't formatted as integer milliseconds since epoch.")
	gotimefmt    = pflag.StringP("gotime-fmt", "", "", "A go time compatible date format specifier. Use if your data isn't formatted as integer milliseconds since epoch.")
	helpFlag     = pflag.BoolP("help", "h", false, "Print usage and exit")
)

func smax(v interface{}, l int) string {
	s := fmt.Sprintf("%v", v)
	if len(s) <= l {
		return s
	}
	return s[:l]
}

func PrintUsage() {
	// This copy-paste of the code from pflag.Usage() is done so we can
	// wrap our usage messages automatically.
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	usages := pflag.CommandLine.FlagUsagesWrapped(90)
	fmt.Fprintf(os.Stderr, usages)
	fmt.Fprintf(os.Stderr, `
Examples:

	# Generate a graph with some fake data
	$ %s --generate-fake-data | %s

	# Graph the data in 1-minute wide bins
	$ %s --generate-fake-data | %s --unit 1minute

	# Parse timestamps in a custom format
	$ cat /tmp/file_with_timestamps | %s --strptime-fmt "%%Y-%%m-%%dT%%H:%%M:%%S.%%f"

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	pflag.Parse()
	if *helpFlag {
		PrintUsage()
		os.Exit(0)
	}
	if *generateData {
		tss, err := tbin.SimpleRandomTimestamps(10000, 36)
		if err != nil {
			fmt.Printf("cannot generate random timestamps: %q", err.Error())
			os.Exit(2)
		}
		for _, ts := range tss {
			if *strptimefmt != "" {
				fmt.Printf("%s\n", timefmt.Format(time.UnixMilli(ts).UTC(), *strptimefmt))
			} else {
				fmt.Printf("%d\n", ts)
			}
		}
		os.Exit(0)
	}
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Printf("You must pipe the timestamps into this program\non stdin; since stdin is a terminal, exiting.\n")
		fmt.Printf(`
    HINT: to see an example interactive graph, run the following command

        %s --generate-fake-data | %s

`, os.Args[0], os.Args[0])
		PrintUsage()
		os.Exit(1)
	}

	parsefunc, _, err := timeformat.NewFuncs(*strptimefmt, *gotimefmt)
	if err != nil {
		fmt.Printf("cannot figure out how to parse the ")
	}

	// 1. read stdin for lines of text
	// 2. Attempt to parse lines of text into dates then into epoch_ms formats
	// 3. Bin each timestamp by the interval
	// 4. Assemble the JSON data to be graphed
	// 5. Render the HTML/JS+JSON data to a tmp file
	// 6. Serve the tmp file from a port
	// 7. Launch a web-browser to view the localhost port

	tss, err := read_lines_to_integers(os.Stdin, parsefunc)
	if err != nil {
		fmt.Printf("cannot divide timestamps into bins: %q", err.Error())
		os.Exit(2)
	}

	*unit = strings.ToLower(*unit)
	if *unit == "auto" {
		*unit, _ = tbin.EstimateBinSize(tss)
	}

	bins, err := tbin.BinTimestamps(tss, *unit)
	if err != nil {
		fmt.Printf("cannot divide timestamps into bins: %q", err.Error())
		os.Exit(2)
	}

	ctx, err := tbin.FormatBinDataForChartJS(bins)
	if err != nil {
		fmt.Printf("cannot convert binned timestamp data into ChartJS data: %q", err.Error())
		os.Exit(2)
	}
	ctxjson, err := json.MarshalIndent(ctx, "", "    ")
	if err != nil {
		fmt.Printf("cannot marshal ChartJS data into JSON format: %q", err.Error())
		os.Exit(2)
	}

	// Asterisk tell CreateTemp where to put a random filename component, which
	// we want to avoid collisions.
	tmpfn := fmt.Sprintf("%d_*_histogram_timestamps.html", time.Now().Unix())
	tmpdir := *outputpath
	f, err := os.CreateTemp(tmpdir, tmpfn)
	if err != nil {
		fmt.Printf("cannot open temporary file for recording HTML: %q", err.Error())
		os.Exit(2)
	}
	defer f.Close()

	jslib := BundleJS
	html_tmplfile := IndexHTML

	html_tmplfile = strings.ReplaceAll(html_tmplfile, "REPLACE_ME_WITH_JS_CONTEXT", string(ctxjson))
	html_tmplfile = strings.ReplaceAll(html_tmplfile, "REPLACE_ME_WITH_BUNDLEJS", jslib)
	html_tmplfile = strings.ReplaceAll(html_tmplfile, "TITLE_HERE", *title)

	fmt.Fprintf(f, "%s\n", html_tmplfile)
	fmt.Printf("Wrote new HTML view file to file %q at path %q\n", f.Name(), tmpdir)
	f.Close()

	mux := http.NewServeMux()
	{
		absp, err := filepath.Abs(tmpdir)
		if err != nil {
			fmt.Printf("cannot determine abs path to temporary directory: %q", err.Error())
			os.Exit(2)
		}
		fullFP := filepath.Join(absp, f.Name())
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			http.ServeFile(w, req, fullFP)
		})
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	localURL := fmt.Sprintf("http://localhost:%d", listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("Visit the newly generated graph of timestamps at URL: %s\n", localURL)
	go func() {
		time.Sleep(time.Millisecond * 500)
		openbrowser(localURL)
	}()
	err = http.Serve(listener, mux)
	if err != nil {
		fmt.Printf("error when serving a directory: %q", err.Error())
		os.Exit(2)
	}
}

// read_lines_to_integers attempts to parse each non-empty lines in r as a time
// parsable by parsefunc. Each integer in the output represents a count of
// milliseconds since UNIX epoch.
//
// If a line fails to parse, then an error is returned. Before returning though,
// a hint is printed on stderr to the user indicating what that line's timestamp
// format _should_ have been, via our best guess.
func read_lines_to_integers(r io.Reader, parsefunc func(string) (time.Time, error)) ([]int64, error) {
	tss := []int64{}
	scnr := bufio.NewScanner(r)
	var i int = 0
	for scnr.Scan() {
		i += 1
		line := scnr.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ts int64
		var err error
		if parsefunc == nil {
			ts, err = strconv.ParseInt(line, 10, 64)
			if err != nil {
				fmt.Fprint(os.Stderr, timeformat.GuessTimestampFormat(line))
				return nil, fmt.Errorf("cannot parse integer on line %d of stdin: %w\n", i, err)
			}
		} else {
			t, err := parsefunc(line)
			if err != nil {
				fmt.Fprint(os.Stderr, timeformat.GuessTimestampFormat(line))
				return nil, fmt.Errorf("cannot parse line %d of stdin to date: %w", i, err)
			}
			ts = t.UnixMilli()
		}
		tss = append(tss, ts)
	}
	return tss, nil
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("cannot launch web browser: %v\n", err)
	}

}
