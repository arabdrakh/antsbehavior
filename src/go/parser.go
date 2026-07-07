package main
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseFile(filename string) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, cannot open file")
	}
	defer file.Close()

	g := NewGraph()
	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, error reading file")
	}

	g.RawLines = lines

	if errMsg := parseSections(g, lines); errMsg != "" {
		return nil, fmt.Errorf("ERROR: invalid data format, %s", errMsg)
	}

	if errMsg := g.Validate(); errMsg != "" {
		return nil, fmt.Errorf("ERROR: invalid data format, %s", errMsg)
	}

	return g, nil
}

func parseSections(g *Graph, lines []string) string {
	const (
		stateAnts  = iota
		stateRooms
		stateLinks
	)

	state := stateAnts
	nextIsStart := false
	nextIsEnd := false

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "##") {
			switch line {
			case "##start":
				if g.Start != "" {
					return "multiple ##start commands"
				}
				nextIsStart = true
			case "##end":
				if g.End != "" {
					return "multiple ##end commands"
				}
				nextIsEnd = true
			}
			if state == stateAnts {
				return "missing number of ants"
			}
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if state == stateAnts {
			n, err := strconv.Atoi(line)
			if err != nil {
				return "invalid number of ants"
			}
			if n <= 0 {
				return "invalid number of ants"
			}
			g.NumAnts = n
			state = stateRooms
			continue
		}

		if isLinkLine(line) {
			state = stateLinks
			parts := strings.SplitN(line, "-", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return "invalid link format"
			}
			if errMsg := g.AddLink(parts[0], parts[1]); errMsg != "" {
				return errMsg
			}
			continue
		}

		if state == stateLinks {
			return "room defined after links"
		}

		errMsg := parseRoomLine(g, line, nextIsStart, nextIsEnd)
		if errMsg != "" {
			return errMsg
		}
		nextIsStart = false
		nextIsEnd = false
	}

	return ""
}

func isLinkLine(line string) bool {
	if strings.Contains(line, " ") {
		return false
	}
	idx := strings.Index(line, "-")
	if idx <= 0 {
		return false
	}
	right := line[idx+1:]
	return right != ""
}

func parseRoomLine(g *Graph, line string, isStart, isEnd bool) string {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return "invalid room format"
	}

	name := parts[0]
	if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
		return "invalid room name"
	}

	x, errX := strconv.Atoi(parts[1])
	y, errY := strconv.Atoi(parts[2])
	if errX != nil || errY != nil {
		return "invalid room coordinates"
	}

	if errMsg := g.AddRoom(name, x, y); errMsg != "" {
		return errMsg
	}

	if isStart {
		g.Start = name
	}
	if isEnd {
		g.End = name
	}

	return ""
}

func PrintInput(g *Graph) {
	for _, line := range g.RawLines {
		fmt.Println(line)
	}
}