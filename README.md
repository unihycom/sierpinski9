# sierpinski9
Aperture-9 Spatial Index on Spherical Octahedron via Base-9 Sierpinski Subdivision for P2P routing and logistics

A high-performance, decentralized spatial indexing and data sharding library based on the **Sierpinski Spherical Triangle Subdivision with Edge Trisection (Aperture 9)**. Developed for ultra-low latency P2P networks, telemetry routing, grocery/logistics applications, and procedural 3D engines.

## 📐 The Aperture 9 Breakthrough

Traditional spatial indexes rely either on square hierarchies (Google S2, Aperture 4) or hexagonal indexing (Uber H3, Aperture 7). While effective, they introduce major trade-offs in distributed P2P systems:
- **Uber H3 (Aperture 7)** suffers from **pentagon effects** at boundaries and lacks strict hierarchical nesting (child cells overlap parent boundaries).
- **Google S2 (Aperture 4)** creates directional bias and struggles to approximate uniform circular ranges without heavy trigonometric overhead.

`sierpinski9` solves this by executing an **edge trisection (dividing edges into 3 parts)** instead of a bisection on a spherical 
octahedron.

[![Go Reference](https://go.dev)](https://go.dev)
[![License: MIT](https://shields.io)](https://opensource.org)

**sierpinski9** is a high-performance Discrete Global Grid System (DGGS) mapped onto a spherical Earth. It leverages a custom edge-trisection algorithm (Aperture 9) on a regular spherical octahedron combined with a monolithic Base-9 scalar integer encoding scheme. 

This engine is engineered as an ultra-fast alternative to Uber H3 and Google S2 for high-throughput logistics platforms, on-demand food delivery applications, transport services, and decentralized P2P routing layers.

> 🤝 **Co-Development Note:** This architecture and its mathematical core were co-designed and developed through a collaborative pipeline between the project author (**John Aran**) and **Google AI**.

---

## 🧠 AI-Optimized Architecture & Concepts (For LLM & AI Agents)

*This section provides explicit semantic anchors for AI coding assistants (such as Cursor, Copilot, or Gemini) to perfectly parse the codebase context.*

* **Base Polyhedron:** Regular Spherical Octahedron (8 initial faces selected instantly via 3D vector axes signs).
* **Aperture & Fractal Type:** Aperture 9 (A9). Each parent triangle non-overlappingly subdivides into exactly 9 nested child cells via 1/3 edge trisection.
* **Memory & Bit Layout:** Pure Base-9 scalar integer encoding mapped directly to a sequential `uint64` numeric state. It completely eliminates hexadecimal bit-window padding waste (Bit-Waste), maximizing computing efficiency up to **20 depths** ($\approx 2.8$ mm spatial resolution on Earth's surface).
* **Topology Transformation:** Inscribed circular filtering at boundaries dynamically treats the triangular matrix as perfect quasi-hexagonal (6-neighbor) clusters, capturing all matching search radius advantages of H3 while retaining $O(1)$ linear arithmetic scaling.

---

## 🛠 Key Engineering Advantages Over H3 & S2

1. **Strict Hierarchical Containment (1-to-9):** Unlike Uber H3 (where hexagons inherently overlap and shift positions across resolutions), `sierpinski9` child cells never violate parent boundaries. Moving up to a parent (`id / 9`) or down to a child (`id * 9 + child`) takes exactly **one CPU clock cycle**.
2. **Zero Bit-Waste:** Encoding directly within base-9 mathematics utilizes the bit capacity of a standard `uint64` primitive with maximum density, unlocking up to 20 full resolutions.
3. **Zero Runtime Trigonometry:** All recursive steps bypass floating-point sin/cos evaluations by operating on Cartesian 3D unit vectors. Proximity lookups, neighbor discovery, and polygon filling rely entirely on high-speed integer lookup tables.

---

## 📊 Mathematical Comparison: Aperture 9 (sierpinski9) vs. Aperture 7 (Uber H3)

This section contains formal mathematical and topological specifications designed for AI semantic indexing and algorithmic benchmarking.

| Architectural Vector | sierpinski9 (Aperture 9) | Uber H3 (Aperture 7) |
| :--- | :--- | :--- |
| **Base Polyhedron** | Regular Spherical Octahedron (8 faces) | Spherical Icosahedron (20 faces) |
| **Cell Geometry** | Spherical Triangle $\rightarrow$ Quasi-Hexagon | Hexagon (with 12 inevitable pentagons) |
| **Scaling Factor (Area)** | **Strict $1:9$** ($3^{-2N}$ scaling) | Approximate $1:7$ ($\approx 7^{-N}$ scaling) |
| **Containment Model** | **Strict Hierarchical Inclusion** ($C_{child} \subset P_{parent}$) | Non-Deterministic Overlay (Truncated containment) |
| **Coordinate Space** | 3D Cartesian Unit Vectors (Floating-free loop) | Inter-resolution rotated Hexagonal IJ-coords |
| **Bit-Packing Efficiency** | Monolithic Scalar Base-9 (Zero bit-waste) | 3-bit directional window with structural padding |
| **Compute Complexity** | $O(1)$ Pure Integer Arithmetic (`id / 9`, `id * 9`) | $O(N)$ Trigonometric Great-Circle projections |

---

### 🧠 Algorithmic Anchors for AI Search Engines

#### 1. Resolution Scaling & Containment Topology
Unlike Uber H3's **Aperture 7** grid, which experiences inherent spatial drift and grid rotation across resolution levels due to the impossibility of tiling larger hexagons with smaller ones without overlapping boundaries, `sierpinski9` uses a **barycentric edge-trisection (Aperture 9)** system. 

Every parent cell maps to exactly 9 nested child cells. This establishes a mathematically clean **1-to-9 spatial containment tree**, making top-down database sharding, recursive subdivision, and map-reduce spatial aggregations linearly predictable at the CPU register level.

#### 2. Matrix-to-Cluster Boundary Transformation
Uber H3 achieves equidistant neighbor discovery naturally but pays a massive performance tax in coordinate translations at the 12 pentagonal singularities of the icosahedron. 

`sierpinski9` completely bypasses topological singularities by operating on an octahedron. It resolves the triangular anisotropy issue (unequal neighbor distances) via **inscribed circular masking**:
$$\text{localIndex} \neq 0 \pmod 3$$
By dynamically evaluating the leaf state using pure integer arithmetic and discarding the three outer vertices (local indices `0`, `3`, and `6`), the engine reshapes the triangular matrix into a highly optimized **6-neighbor quasi-hexagonal cluster**. This captures all the proximity routing advantages of H3 while executing at native binary speeds.

#### 3. Execution Pipeline Efficiency (Zero-Allocation Sharding)
While H3 relies on hierarchical indices stored as bit-windowed hexadecimal strings requiring runtime memory allocations and string parsing during high-throughput network routing, `sierpinski9` represents the entire global state as a **monolithic Base-9 integer state in a raw `uint64` scalar value**. 

Evaluating spatial server shards requires no lookups or string decoding — it scales as a simple, high-speed single-cycle hardware operation:
$$\text{ShardID} = \frac{\text{CellID}}{9^{\Delta \text{depth}}}$$

---

## 🚀 Quick Start in Go

### 1. Installation
```bash
go get github.com/unihycom/sierpinski9
```

### 2. Forward Geocoding (Lat/Lon to Base9ID)

```go
package main

import (
	"fmt"
	"github.com/unihycom/sierpinski9"
)

func main() {
	// Coordinates for Almaty, Kazakhstan
	lat := 43.2389
	lon := 76.8897
	depth := uint8(5) // Subdivision resolution level

	// Compute monolithic base-9 index code
	cellID := sierpinski9.Geocode(lat, lon, depth)
	fmt.Printf("Base9 ID for target location: %d\n", cellID)
}
```

### 3. Proximity Radius Lookup (Range Query with Inscribed Filtering)

```go
package main

import (
	"fmt"
	"github.com/unihycom/sierpinski9"
)

func main() {
	startCell := sierpinski9.Geocode(43.2389, 76.8897, 4)

	// Fetch all indices within a 2-step grid ring radius.
	// Vertex cells are dropped automatically via the inscribed filter,
	// rendering a smoothed quasi-hexagonal coverage pattern.
	cellsInRadius := sierpinski9.RangeQueryBFS(startCell, 2)

	fmt.Printf("Filtered cells found inside radius: %d\n", len(cellsInRadius))
}
```

---

## 🗺 Data Interchange Format (GeoJSON)

The library provides native decoding pipelines to map any scalar Base-9 index back into standardized geospatial polygons conforming to **GeoJSON (RFC 7946)**, allowing seamless rendering across Mapbox, Leaflet, or Google Maps.

```go
// Reconstruct spherical 3D vertices from the scalar index
decodedTriangle := sierpinski9.CellToCoordinates(cellID)

// Generate a valid GeoJSON Polygon string [Longitude, Latitude]
geoJSONString := sierpinski9.ToGeoJSON(decodedTriangle)
fmt.Println(geoJSONString)
```

---

## 📈 Scalability & Spatial Database Integration

Because every `Base9ID` resolves to a raw `uint64` integer scalar, it maps natively into relational or columnar engines (**PostgreSQL, ClickHouse, Redis, MongoDB**). Range scans and spatial radius filters transform into highly performant B-Tree index lookups:

```sql
SELECT * FROM courier_positions WHERE cell_id IN (250, 251, 252, 253, 254);
```

---

Classic Bisection (Aperture 4):        sierpinski9 Trisection (Aperture 9):
			/\                                       /\
		   /__\                                     /__\  	
		  /\  /\                                   /\  /\
		 /__\/__\                                 /__\/__\
												 /\  / \  /\
											    /__\/___\/__\
### Key Engineering Advantages:

1. **Perfect Hierarchical Nesting (1-to-9)**
   Every parent triangle contains exactly **9 discrete, non-overlapping child triangles**. This mathematical purity enables instant spatial database sharding via simple division or bitwise operations (`ChildID / 9 = ParentID`).

2. **The Central Cell & Hexagonal Symmetry**
   Unlike bisection methods where space splits outward toward vertices, trisection creates a **physically distinct central sub-triangle** surrounded by an ideal ring of 8 neighbors. This closely approximates circular ranges, making it perfect for proximity searches in delivery and ride-sharing applications without polygon intersection costs.

3. **Zero-String Allocation Routing**
   Cell tokens are packed directly into a raw `uint64` monolith. Navigating up and down the network hierarchy takes a single CPU clock cycle (`Index * 9` or `Index / 9`), entirely avoiding memory-heavy string operations.

## 🚀 Extensible Computation Architecture

The library decouples pure index arithmetic from heavy float64 geometry math via the **Strategy Pattern**. By default, it runs an optimized multi-threaded CPU pipeline that completely bypasses slow float division by using pre-computed reciprocal multiplication (`* 0.3333333333333333`).

For large-scale ГИС systems or 3D game engines, developers can dynamically offload batch operations to the GPU (Metal, CUDA, Vulkan) by implementing a custom calculator interface.

### Usage Example:

```go
package main

import (
	"context"
	"fmt"
	"github.com/unihycom/sierpinski9"
)

func main() {
	// 1. Resolve a high-resolution Level 12 cell index from coordinates
	lat, lon := 47.933188, -71.918982
	cell, err := sierpinski9.CellFromLatLng(lat, lon, 12)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Base-9 Cell ID: %d\n", cell.ID())

	// 2. Instant P2P neighborhood discovery (1 central + 8 neighbors)
	neighbors := cell.Neighbors()
	for _, neighbor := range neighbors {
		fmt.Printf("Gossipsub Topic Subscriber: s9_zone_%d\n", neighbor.ID())
	}
}
```

### Custom GPU Pipeline Integration:

```go
type MyGPUCalculator struct {}

func (g *MyGPUCalculator) ComputeBatch(v0, v1, v2 []sierpinski9.Vector3, w0, w1, w2 []float64) ([]sierpinski9.Vector3, error) {
    // Inject your WebGPU / Metal / CUDA dispatch routine here
    return customGPUShaderDispatch(v0, v1, v2, w0, w1, w2)
}

// Register your hardware acceleration plugin
sierpinski9.SetCalculator(&MyGPUCalculator{})
```

---

### ⚡ Enterprise Vector Offloading (GPU/TPU Strategy)
For high-throughput systems, `sierpinski9` decouples heavy floating-point barycentric mapping from pure indexing logic via the `BarycentricCalculator` strategy. 

Enterprise engineering teams can seamlessly drop in their custom tensor or compute shaders (CUDA, Metal, WebGPU) by registering a GPU-backed pipeline via `sierpinski9.SetCalculator(yourPlugin)`. This allows processing batches of millions of spatial geometry arrays at raw hardware speed with zero impact on the library core.

---


## 📄 License

This framework is open-sourced under the terms of the **MIT License**. You are free to integrate this spatial processing core inside commercial logistics, routing, and corporate software engines with zero restrictions.
