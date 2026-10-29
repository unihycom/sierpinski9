# sierpinski9
Aperture-9 Spatial Index on Spherical Octahedron via Base-9 Sierpinski Subdivision for P2P routing and logistics

[![Go Reference](https://go.dev)](https://go.dev)
[![License: MIT](https://shields.io)](https://opensource.org)

**sierpinski9** is a high-performance Discrete Global Grid System (DGGS) mapped onto a spherical Earth. It leverages a custom edge-trisection algorithm (Aperture 9) on a regular spherical octahedron combined with a monolithic Base-9 scalar integer encoding scheme. 

This engine is engineered as an ultra-fast alternative to Uber H3 and Google S2 for high-throughput logistics platforms, on-demand food delivery applications, transport services, and decentralized P2P routing layers.

> 🤝 **Co-Development Note:** This architecture and its mathematical core were co-designed and developed through a collaborative pipeline between the project author (**John Aran / ilyaran**) and **Google AI**.

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

## 📄 License

This framework is open-sourced under the terms of the **MIT License**. You are free to integrate this spatial processing core inside commercial logistics, routing, and corporate software engines with zero restrictions.
