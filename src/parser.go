package lemin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Graph holds the entire ant farm
type Graph struct {
	NumAnts   int
	Rooms     map[string]*Room
	RoomOrder []string // preserves input order for output
	Links     []Link
	Start     string
	End       string
}

// Room represents a node in the ant farm
type Room struct {
	Name string
	X, Y int
}

// Link represents a tunnel between two rooms
type Link struct {
	From, To string
}

// ParseFile reads and validates the input file, returns a Graph or an error
func ParseFile(path string) (*Graph, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, cannot open file")
	}
	defer f.Close()

	g := &Graph{
		Rooms: make(map[string]*Room),
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, file read error")
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, empty file")
	}

	// --- Parse number of ants ---
	numAnts, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil || numAnts <= 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
	}
	g.NumAnts = numAnts

	// --- Parse rooms and links ---
	nextIsStart := false
	nextIsEnd := false
	parsingLinks := false

	seenRooms := make(map[string]bool)
	seenLinks := make(map[string]bool)

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle commands
		if line == "##start" {
			nextIsStart = true
			nextIsEnd = false
			continue
		}
		if line == "##end" {
			nextIsEnd = true
			nextIsStart = false
			continue
		}
		// Skip comments (but ##start and ##end are handled above)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this line is a link (contains "-" but not a room definition)
		if isLink(line) {
			parsingLinks = true
			parts := strings.SplitN(line, "-", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("ERROR: invalid data format, malformed link: %s", line)
			}
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])

			// Validate rooms exist (we allow forward refs: validate at end)
			if from == to {
				return nil, fmt.Errorf("ERROR: invalid data format, room links to itself: %s", line)
			}

			// Deduplicate links
			key1 := from + "-" + to
			key2 := to + "-" + from
			if seenLinks[key1] || seenLinks[key2] {
				return nil, fmt.Errorf("ERROR: invalid data format, duplicate link: %s", line)
			}
			seenLinks[key1] = true
			g.Links = append(g.Links, Link{From: from, To: to})
			continue
		}

		// If we're already parsing links, a room-like line is an error
		if parsingLinks {
			// unknown commands are ignored per spec
			continue
		}

		// Parse room: "name x y"
		room, err := parseRoom(line)
		if err != nil {
			return nil, fmt.Errorf("ERROR: invalid data format, %v", err)
		}

		// Room name must not start with 'L' or '#'
		if strings.HasPrefix(room.Name, "L") || strings.HasPrefix(room.Name, "#") {
			return nil, fmt.Errorf("ERROR: invalid data format, room name cannot start with L or #: %s", room.Name)
		}

		// No duplicate rooms
		if seenRooms[room.Name] {
			return nil, fmt.Errorf("ERROR: invalid data format, duplicate room: %s", room.Name)
		}
		seenRooms[room.Name] = true

		g.Rooms[room.Name] = room
		g.RoomOrder = append(g.RoomOrder, room.Name)

		if nextIsStart {
			if g.Start != "" {
				return nil, fmt.Errorf("ERROR: invalid data format, multiple start rooms")
			}
			g.Start = room.Name
			nextIsStart = false
		} else if nextIsEnd {
			if g.End != "" {
				return nil, fmt.Errorf("ERROR: invalid data format, multiple end rooms")
			}
			g.End = room.Name
			nextIsEnd = false
		}
	}

	// --- Final validations ---
	if g.Start == "" {
		return nil, fmt.Errorf("ERROR: invalid data format, no start room found")
	}
	if g.End == "" {
		return nil, fmt.Errorf("ERROR: invalid data format, no end room found")
	}

	// Validate all link endpoints exist
	for _, lk := range g.Links {
		if _, ok := g.Rooms[lk.From]; !ok {
			return nil, fmt.Errorf("ERROR: invalid data format, link references unknown room: %s", lk.From)
		}
		if _, ok := g.Rooms[lk.To]; !ok {
			return nil, fmt.Errorf("ERROR: invalid data format, link references unknown room: %s", lk.To)
		}
	}

	return g, nil
}

// parseRoom parses a "name x y" line into a Room
func parseRoom(line string) (*Room, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid room definition: %s", line)
	}
	name := parts[0]
	if strings.ContainsAny(name, " \t") {
		return nil, fmt.Errorf("room name contains spaces: %s", name)
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("room x coordinate not an integer: %s", parts[1])
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("room y coordinate not an integer: %s", parts[2])
	}
	return &Room{Name: name, X: x, Y: y}, nil
}

// isLink returns true if the line looks like "roomA-roomB"
// It must contain exactly one "-" and no spaces
func isLink(line string) bool {
	if strings.Contains(line, " ") {
		return false
	}
	idx := strings.Index(line, "-")
	if idx <= 0 {
		return false
	}
	// both sides must be non-empty
	left := line[:idx]
	right := line[idx+1:]
	return len(left) > 0 && len(right) > 0
}

// PrintFarm prints the ant farm section (number + rooms + links) as required by the spec
func PrintFarm(g *Graph, originalLines []string) {
	for _, line := range originalLines {
		fmt.Println(line)
	}
}

// ReadLines returns the raw lines of a file (for echoing back in output)
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}