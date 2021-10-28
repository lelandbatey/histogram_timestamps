package main

import (
	"fmt"
	"io"
	"os"

	_ "github.com/itchyny/timefmt-go"
	isatty "github.com/mattn/go-isatty"
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
	fmt.Printf("hello world")
}
