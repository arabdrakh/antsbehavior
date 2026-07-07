package main
type Room struct {
	Name string
	X, Y int
}

type Graph struct {
	Rooms     map[string]*Room
	Links     map[string][]string
	Start     string
	End       string
	NumAnts   int
	RoomOrder []string
	RawLines  []string
}

func NewGraph() *Graph {
	return &Graph{
		Rooms: make(map[string]*Room),
		Links: make(map[string][]string),
	}
}

func (g *Graph) AddRoom(name string, x, y int) string {
	if _, exists := g.Rooms[name]; exists {
		return "duplicate room: " + name
	}
	g.Rooms[name] = &Room{Name: name, X: x, Y: y}
	g.RoomOrder = append(g.RoomOrder, name)
	return ""
}

func (g *Graph) AddLink(a, b string) string {
	if a == b {
		return "room links to itself: " + a
	}
	if _, ok := g.Rooms[a]; !ok {
		return "unknown room in link: " + a
	}
	if _, ok := g.Rooms[b]; !ok {
		return "unknown room in link: " + b
	}
	for _, neighbor := range g.Links[a] {
		if neighbor == b {
			return "duplicate tunnel: " + a + "-" + b
		}
	}
	g.Links[a] = append(g.Links[a], b)
	g.Links[b] = append(g.Links[b], a)
	return ""
}

func (g *Graph) Validate() string {
	if g.NumAnts <= 0 {
		return "invalid number of ants"
	}
	if g.Start == "" {
		return "no start room found"
	}
	if g.End == "" {
		return "no end room found"
	}
	if g.Start == g.End {
		return "start and end are the same room"
	}
	if !g.isReachable() {
		return "no path between start and end"
	}
	return ""
}

func (g *Graph) isReachable() bool {
	visited := make(map[string]bool)
	queue := []string{g.Start}
	visited[g.Start] = true
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == g.End {
			return true
		}
		for _, neighbor := range g.Links[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return false
}