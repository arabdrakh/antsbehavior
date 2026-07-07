# *Ants Behavior ( Lem-in)*

## Project Description

The project reads a description of an ant farm (rooms, tunnels, ant count) from a text file and finds the **quickest path** to move all ants from `##start` to `##end` in the fewest number of turns.

## The Algorithm

1. **Parse & Validate** Reads the input format, handles coordinate mapping, and ensures the structural integrity of the graph (no disconnected components, unique room names).
2. **Pathfinding Optimization** Uses Depth-First Search (DFS) combined with a Max-Flow approach to discover all simple paths and select the largest possible subset of room-disjoint paths.
3. **Ant Distribution Pipeline** Distributes ants across available paths using a greedy pipeline formula to minimize the absolute finishing turn
4. **Simulation** Moves ants turn-by-turn. The output matches the standard `Lx-room` format.

---

## Technical Features: 3D Browser Visualizer (Bonus Track)

To eliminate bottleneck stutters when processing heavy high-volume stress tests (1000+ ants), we developed a client-side **3D Web Visualizer** powered by Three.js implementing low-level GPU acceleration techniques.

### Key Architectural Upgrades:
* **InstancedMesh GPU Rendering:** Compresses the drawing overhead of thousands of independent models down to a **single Draw Call**. Instead of binding thousands of nodes on the CPU, spatial transformation matrices are computed concurrently and flushed directly into GPU memory.
* **Zero-Server Portable Runtime:** Operates completely client-side. No local servers or CORS-bypass proxies required. Works fully offline by dropping exported simulation JSON traces straight into the browser interface.
* **Procedural Locomotive Scaling:** Implements a dynamic matrix deformation loop. Moving agents dynamically squeeze and expand along their orientation axes to simulate kinetic muscle biomechanics under high-performance constraints.
* **Cyberpunk HUD & Graph Architecture:** High-contrast matrix backdrop, translucent glassmorphic console docks, custom `JetBrains Mono` text buffers, neon edge-glow emissive borders for node boundaries, and non-blocking white label sprites (`depthTest: false`) that remain strictly visible during camera orbits.
---

## How to Run

### 1. Execute Core Pathfinder
```bash
python -m pip install -r requirements.txt
go run ./src/go example/example00.txt

## How to Run

```bash
python -m pip install -r requirements.txt
go run ./src/go example/example00.txt


# Animated terminal visualizer (2D)
go run ./src/go example/example05.txt | py src/visualizer3d.py
go run ./src/go example/example05.txt | py src/visualizer3d.py --fast
go run ./src/go example/example05.txt | py src/visualizer3d.py --slow

## 3D Visualizer (Bonus)

The project also includes 3D visualizer in the browser,
plus a terminal isometric fallback for environments without a display.

### Browser visualizer (Three.js)

1. Build the project and export the simulation as JSON alongside the normal output:

```bash
go build ./src/go
go run ./src/go --json colony.json example/example05.txt
```

2. Open `src/visualizer3d.html` in a browser:

```bash
# Windows (Git Bash)
start src/visualizer3d.html

# macOS
open src/visualizer3d.html

# Linux
xdg-open src/visualizer3d.html
```

3. Drag `colony.json` onto the page (or click to browse). Rooms render as
glowing 3D spheres connected by tunnels; ants appear as orange spheres that
move turn by turn. Drag to orbit the camera, scroll to zoom, use the
play/pause/reset controls and scrubber at the bottom to control playback.

### Terminal isometric fallback (no browser)

```bash
py src/visualizer3d.py colony.json
py src/visualizer3d.py colony.json --fast
py src/visualizer3d.py colony.json --slow

# or pipe directly without generating JSON first
go run ./src/go example/example05.txt | py src/visualizer3d.py
```

## Error Handling

Invalid input prints to stderr and exits with code 1:

```
ERROR: invalid data format ...
```
but with specific reasons
badexample01.txt:

```bash
go run ./src/go example/badexample01.txt
```
```
ERROR: invalid data format, room links to itself: 3
exit status 1
```

badexample00.txt:

```bash
go run ./src/go example/badexample00.txt
```
ERROR: invalid data format, invalid number of ants
exit status 1
```


## Demo Video

`https://youtu.be/XxJyT1cq0Z8?si=bf_Gtj64op64ScS4`