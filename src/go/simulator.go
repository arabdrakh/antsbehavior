package main
import (
	"fmt"
	"strings"
)

type AntState struct {
	ID       int
	Position string
	Done     bool
}

func Simulate(g *Graph, ps PathSet) []string {
	if len(ps.Paths) == 0 || len(ps.AntLoads) == 0 {
		return nil
	}

	antPath := make([]int, g.NumAnts)
	antStep := make([]int, g.NumAnts)
	antDone := make([]bool, g.NumAnts)

	antID := 0
	for pathIdx, load := range ps.AntLoads {
		for k := 0; k < load; k++ {
			antPath[antID] = pathIdx
			antID++
		}
	}

	occupied := make(map[string]bool)

	antRelease := make([]int, g.NumAnts)
	pathNextRelease := make([]int, len(ps.Paths))
	for i := 0; i < len(ps.Paths); i++ {
		pathNextRelease[i] = 1
	}
	for i := 0; i < g.NumAnts; i++ {
		pi := antPath[i]
		antRelease[i] = pathNextRelease[pi]
		pathNextRelease[pi]++
	}

	var turns []string

	for {
		allDone := true
		for i := 0; i < g.NumAnts; i++ {
			if !antDone[i] {
				allDone = false
				break
			}
		}
		if allDone {
			break
		}

		currentTurn := len(turns) + 1
		var moves []string

		for i := 0; i < g.NumAnts; i++ {
			if antDone[i] {
				continue
			}
			if antRelease[i] > currentTurn {
				continue
			}

			path := ps.Paths[antPath[i]].Rooms
			nextStep := antStep[i] + 1
			if nextStep >= len(path) {
				antDone[i] = true
				continue
			}

			nextRoom := path[nextStep]

			if nextRoom != g.End && occupied[nextRoom] {
				continue // blocked this turn, try next turn
			}

			// Free current room if it's an interior room
			currentRoom := path[antStep[i]]
			if currentRoom != g.Start && currentRoom != g.End {
				occupied[currentRoom] = false
			}

			antStep[i] = nextStep
			if nextRoom != g.End {
				occupied[nextRoom] = true
			}
			if nextStep == len(path)-1 {
				antDone[i] = true
			}

			moves = append(moves, fmt.Sprintf("L%d-%s", i+1, nextRoom))
		}

		if len(moves) > 0 {
			turns = append(turns, strings.Join(moves, " "))
		} else {
			break
		}
	}

	return turns
}