package main
import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: go run . [--json out.json] <filename>")
		os.Exit(1)
	}

	jsonPath := ""
	var filename string

	for i := 0; i < len(args); i++ {
		if args[i] == "--json" && i+1 < len(args) {
			jsonPath = args[i+1]
			i++
		} else {
			filename = args[i]
		}
	}

	if filename == "" {
		fmt.Fprintln(os.Stderr, "Usage: go run . [--json out.json] <filename>")
		os.Exit(1)
	}

	g, err := ParseFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	PrintInput(g)
	fmt.Println()

	ps := FindBestPaths(g)
	if len(ps.Paths) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format, no path found from start to end")
		os.Exit(1)
	}

	moves := Simulate(g, ps)
	for _, turn := range moves {
		fmt.Println(turn)
	}

	if jsonPath != "" {
		if err := WriteJSON(g, moves, jsonPath); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: could not write JSON:", err)
		}
	}
}