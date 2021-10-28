package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/itchyny/timefmt-go"
	isatty "github.com/mattn/go-isatty"

	"github.com/lelandbatey/histogram_timestamps/tbin"
)

func main() {
	fmt.Printf("Is stdin a terminal?: %t\n", isatty.IsTerminal(os.Stdin.Fd()))
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		all, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading stdin: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("Read the following data:\n%s\n", string(all))
	}
	// 1. read stdin for lines of text
	// 2. Attempt to parse lines of text into dates then into epoch_ms formats
	// 3. Bin each timestamp by the interval
	// 4. Assemble the JSON data to be graphed
	// 5. Render the HTML/JS+JSON data to a tmp file
	// 6. Serve the tmp file from a port
	// 7. Launch a web-browser to view the localhost port

	tss, err := tbin.SimpleRandomTimestamps(10000, 12)
	if err != nil {
		fmt.Printf("cannot generate random timestamps: %q", err.Error())
		os.Exit(2)
	}

	bins, err := tbin.BinTimestamps(tss, "hour")
	if err != nil {
		fmt.Printf("cannot divide timestamps into bins: %q", err.Error())
		os.Exit(2)
	}

	ctx, err := tbin.FormatBinDataForChartJS(bins)
	if err != nil {
		fmt.Printf("cannot convert binned timestamp data into ChartJS data: %q", err.Error())
		os.Exit(2)
	}
	ctxjson, err := json.Marshal(ctx)
	if err != nil {
		fmt.Printf("cannot marshal ChartJS data into JSON format: %q", err.Error())
		os.Exit(2)
	}

	// Asterisk tell CreateTemp where to put a random filename component, which
	// we want to avoid collisions.
	tmpfn := fmt.Sprintf("%d_*_histogram_timestamps.html", time.Now().Unix())
	//tmpdir := os.TempDir()
	tmpdir := "./"
	f, err := os.CreateTemp(tmpdir, tmpfn)
	if err != nil {
		fmt.Printf("cannot open temporary file for recording HTML: %q", err.Error())
		os.Exit(2)
	}
	defer f.Close()

	jslib := MustAssetString("bundle.js")
	html_tmplfile := MustAssetString("index.html")

	strings.ReplaceAll(html_tmplfile, "REPLACE_ME_WITH_JS_CONTEXT", string(ctxjson))
	strings.ReplaceAll(html_tmplfile, "REPLACE_ME_WITH_BUNDLEJS", jslib)

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
	fmt.Printf("Visit the newly generated graph of timestamps at URL: http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)
	err = http.Serve(listener, mux)
	if err != nil {
		fmt.Printf("error when serving a directory: %q", err.Error())
		os.Exit(2)
	}
}
