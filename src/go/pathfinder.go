package main
type Path struct {
	Rooms []string
}

type PathSet struct {
	Paths    []Path
	AntLoads []int
}

func FindBestPaths(g *Graph) PathSet {
	allDisjoint := findDisjointPathsFlow(g)
	if len(allDisjoint) == 0 {
		return PathSet{}
	}

	best := PathSet{}
	bestTurns := -1

	for k := 1; k <= len(allDisjoint); k++ {
		subset := allDisjoint[:k]
		loads := distributeAnts(g.NumAnts, subset)
		turns := calcTurns(subset, loads)
		if bestTurns == -1 || turns < bestTurns {
			bestTurns = turns
			cp := make([]Path, k)
			copy(cp, subset)
			cpLoads := make([]int, k)
			copy(cpLoads, loads)
			best = PathSet{Paths: cp, AntLoads: cpLoads}
		}
	}

	return best
}

const (
	nodeIn  = 0
	nodeOut = 1
)

type flowEdge struct {
	to, rev int
	cap     int
}

type flowGraph struct {
	adj   [][]int    
	edges []flowEdge
}

func newFlowGraph(n int) *flowGraph {
	return &flowGraph{adj: make([][]int, n)}
}

func (fg *flowGraph) addEdge(u, v, cap int) {
	fg.adj[u] = append(fg.adj[u], len(fg.edges))
	fg.edges = append(fg.edges, flowEdge{v, len(fg.edges) + 1, cap})
	fg.adj[v] = append(fg.adj[v], len(fg.edges))
	fg.edges = append(fg.edges, flowEdge{u, len(fg.edges) - 1, 0})
}

func (fg *flowGraph) bfsPath(src, dst int) []int {
	prev := make([]int, len(fg.adj))
	for i := range prev {
		prev[i] = -1
	}
	prev[src] = -2
	queue := []int{src}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, ei := range fg.adj[u] {
			e := fg.edges[ei]
			if e.cap > 0 && prev[e.to] == -1 {
				prev[e.to] = ei
				if e.to == dst {
					return prev
				}
				queue = append(queue, e.to)
			}
		}
	}
	return nil
}

func (fg *flowGraph) maxFlow(src, dst int) int {
	flow := 0
	for {
		prev := fg.bfsPath(src, dst)
		if prev == nil {
			break
		}
		v := dst
		for v != src {
			ei := prev[v]
			fg.edges[ei].cap--
			fg.edges[fg.edges[ei].rev].cap++
			v = fg.edges[fg.edges[ei].rev].to
		}
		flow++
	}
	return flow
}

func findDisjointPathsFlow(g *Graph) []Path {
	roomIdx := make(map[string]int)
	idx := 0
	for _, name := range g.RoomOrder {
		roomIdx[name] = idx
		idx++
	}
	n := idx * 2

	fg := newFlowGraph(n)

	nodeID := func(name string, side int) int {
		return roomIdx[name]*2 + side
	}

	for name := range g.Rooms {
		cap := 1
		if name == g.Start || name == g.End {
			cap = 1000
		}
		fg.addEdge(nodeID(name, nodeIn), nodeID(name, nodeOut), cap)
	}

	seen := make(map[[2]string]bool)
	for a, neighbors := range g.Links {
		for _, b := range neighbors {
			key := [2]string{a, b}
			rev := [2]string{b, a}
			if seen[key] || seen[rev] {
				continue
			}
			seen[key] = true
			fg.addEdge(nodeID(a, nodeOut), nodeID(b, nodeIn), 1)
			fg.addEdge(nodeID(b, nodeOut), nodeID(a, nodeIn), 1)
		}
	}

	src := nodeID(g.Start, nodeIn)
	dst := nodeID(g.End, nodeOut)

	numPaths := fg.maxFlow(src, dst)
	if numPaths == 0 {
		return nil
	}

	paths := extractPaths(fg, g, roomIdx, numPaths)

	sortPathsByLength(paths)
	return paths
}

func extractPaths(fg *flowGraph, g *Graph, roomIdx map[string]int, numPaths int) []Path {
	nodeID := func(name string, side int) int {
		return roomIdx[name]*2 + side
	}

	idxRoom := make(map[int]string)
	for name, i := range roomIdx {
		idxRoom[i*2+nodeIn] = name
		idxRoom[i*2+nodeOut] = name
	}

	var paths []Path

	for p := 0; p < numPaths; p++ {
		var route []string
		current := nodeID(g.Start, nodeOut)
		route = append(route, g.Start)

		visited := make(map[int]bool)
		visited[nodeID(g.Start, nodeIn)] = true
		visited[nodeID(g.Start, nodeOut)] = true

		endIn := nodeID(g.End, nodeIn)

		for current != nodeID(g.End, nodeOut) {
			moved := false
			for _, ei := range fg.adj[current] {
				e := fg.edges[ei]
				// A forward edge that carried flow has rev.cap > 0
				revCap := fg.edges[e.rev].cap
				if revCap > 0 && !visited[e.to] {
					// consume this flow unit
					fg.edges[e.rev].cap--
					fg.edges[ei].cap++

					visited[e.to] = true
					current = e.to
					if current != endIn {
						roomName := idxRoom[current]
						outNode := nodeID(roomName, nodeOut)
						for _, iei := range fg.adj[current] {
							ie := fg.edges[iei]
							if ie.to == outNode {
								revIe := fg.edges[ie.rev]
								_ = revIe
								if fg.edges[ie.rev].cap > 0 {
									fg.edges[ie.rev].cap--
									fg.edges[iei].cap++
									visited[outNode] = true
									current = outNode
									route = append(route, roomName)
									break
								}
							}
						}
					} else {
						route = append(route, g.End)
						current = nodeID(g.End, nodeOut)
					}

					moved = true
					break
				}
			}
			if !moved {
				break
			}
		}

		if len(route) >= 2 && route[len(route)-1] == g.End {
			paths = append(paths, Path{Rooms: route})
		}
	}

	return paths
}


func distributeAnts(numAnts int, paths []Path) []int {
	loads := make([]int, len(paths))
	lengths := make([]int, len(paths))
	for i, p := range paths {
		lengths[i] = len(p.Rooms) - 1
	}
	for ant := 0; ant < numAnts; ant++ {
		best := 0
		bestTurns := finishTurn(lengths[0], loads[0]+1)
		for i := 1; i < len(paths); i++ {
			t := finishTurn(lengths[i], loads[i]+1)
			if t < bestTurns {
				bestTurns = t
				best = i
			}
		}
		loads[best]++
	}
	return loads
}

func finishTurn(length, ants int) int {
	if ants == 0 {
		return 0
	}
	return length + ants - 1
}

func calcTurns(paths []Path, loads []int) int {
	max := 0
	for i, p := range paths {
		length := len(p.Rooms) - 1
		t := finishTurn(length, loads[i])
		if t > max {
			max = t
		}
	}
	return max
}

func sortPathsByLength(paths []Path) {
	n := len(paths)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if len(paths[j].Rooms) < len(paths[i].Rooms) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}
}

func findAllSimplePaths(g *Graph) []Path {
	var results []Path
	visited := make(map[string]bool)
	var dfs func(current string, route []string)
	dfs = func(current string, route []string) {
		if len(results) >= 500 {
			return
		}
		if current == g.End {
			cp := make([]string, len(route))
			copy(cp, route)
			results = append(results, Path{Rooms: cp})
			return
		}
		visited[current] = true
		for _, neighbor := range g.Links[current] {
			if !visited[neighbor] {
				dfs(neighbor, append(route, neighbor))
			}
		}
		visited[current] = false
	}
	dfs(g.Start, []string{g.Start})
	sortPathsByLength(results)
	return results
}

func selectDisjointPaths(g *Graph, candidates []Path) []Path {
	used := make(map[string]bool)
	var chosen []Path
	for _, p := range candidates {
		conflict := false
		for _, r := range p.Rooms[1 : len(p.Rooms)-1] {
			if used[r] {
				conflict = true
				break
			}
		}
		if conflict {
			continue
		}
		for _, r := range p.Rooms[1 : len(p.Rooms)-1] {
			used[r] = true
		}
		chosen = append(chosen, p)
	}
	return chosen
}