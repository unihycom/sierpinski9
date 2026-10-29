package sierpinski9

import (
	"encoding/json"
	"math"
)

// Base9ID инкапсулирует девятеричный гео-индекс
type Base9ID uint64

type Direction int

const (
	Edge0 Direction = iota
	Edge1
	Edge2
)

// IsInscribedCircle возвращает true, если ячейка не является крайней вершиной родителя
func (id Base9ID) IsInscribedCircle() bool {
	if uint64(id) <= 7 {
		return true
	}
	localIndex := uint64(id) % 9
	// В правильной барицентрической сетке углами родителя являются индексы 0, 3 и 6
	return localIndex != 0 && localIndex != 3 && localIndex != 6
}

// PushStep добавляет уровень фрактала
func (id Base9ID) PushStep(childIndex uint64) Base9ID {
	return Base9ID(uint64(id)*9 + childIndex)
}

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) Normalize() Vector3 {
	len := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if len == 0 {
		return Vector3{0, 0, 0}
	}
	return Vector3{v.X / len, v.Y / len, v.Z / len}
}

// ToLatLng переводит 3D-вектор единичной сферы обратно в географические градусы (Широта, Долгота)
func (v Vector3) ToLatLng() (float64, float64) {
	latRad := math.Asin(v.Z)
	lonRad := math.Atan2(v.Y, v.X)

	latDeg := latRad * 180.0 / math.Pi
	lonDeg := lonRad * 180.0 / math.Pi

	return latDeg, lonDeg
}

func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

type TriangleA9 struct {
	ID         uint64
	Depth      uint64
	V0, V1, V2 Vector3
}

// Вспомогательная функция для точного вычисления точек по барицентрическим весам
func barycentric(v0, v1, v2 Vector3, w0, w1, w2 float64) Vector3 {
	return Vector3{
		X: (v0.X*w0 + v1.X*w1 + v2.X*w2) / 3.0,
		Y: (v0.Y*w0 + v1.Y*w1 + v2.Y*w2) / 3.0,
		Z: (v0.Z*w0 + v1.Z*w1 + v2.Z*w2) / 3.0,
	}.Normalize()
}

// SubdivideA9 генерирует 9 дочерних ячеек через барицентрическое разбиение
func (t TriangleA9) SubdivideA9() []TriangleA9 {
	nextDepth := t.Depth + 1
	baseID := t.ID

	// Вычисляем 10 стабильных точек трисекции
	v300 := t.V0
	v210 := barycentric(t.V0, t.V1, t.V2, 2, 1, 0)
	v120 := barycentric(t.V0, t.V1, t.V2, 1, 2, 0)
	v030 := t.V1
	v021 := barycentric(t.V0, t.V1, t.V2, 0, 2, 1)
	v012 := barycentric(t.V0, t.V1, t.V2, 0, 1, 2)
	v003 := t.V2
	v102 := barycentric(t.V0, t.V1, t.V2, 1, 0, 2)
	v201 := barycentric(t.V0, t.V1, t.V2, 2, 0, 1)
	v111 := barycentric(t.V0, t.V1, t.V2, 1, 1, 1) // Единый центр-центроид

	// Формируем 9 смежных треугольников без щелей и наложений
	return []TriangleA9{
		{ID: baseID*9 + 0, Depth: nextDepth, V0: v300, V1: v210, V2: v201}, // Угол V0
		{ID: baseID*9 + 1, Depth: nextDepth, V0: v210, V1: v120, V2: v111},
		{ID: baseID*9 + 2, Depth: nextDepth, V0: v210, V1: v201, V2: v111},
		{ID: baseID*9 + 3, Depth: nextDepth, V0: v120, V1: v030, V2: v021}, // Угол V1
		{ID: baseID*9 + 4, Depth: nextDepth, V0: v120, V1: v021, V2: v111},
		{ID: baseID*9 + 5, Depth: nextDepth, V0: v021, V1: v012, V2: v111},
		{ID: baseID*9 + 6, Depth: nextDepth, V0: v012, V1: v003, V2: v102}, // Угол V2
		{ID: baseID*9 + 7, Depth: nextDepth, V0: v012, V1: v102, V2: v111},
		{ID: baseID*9 + 8, Depth: nextDepth, V0: v201, V1: v102, V2: v111},
	}
}

func getBaseOctantTriangle(octantID uint64) TriangleA9 {
	p0 := Vector3{0, 0, 1}
	p1 := Vector3{1, 0, 0}
	p2 := Vector3{0, 1, 0}
	p3 := Vector3{-1, 0, 0}
	p4 := Vector3{0, -1, 0}
	p5 := Vector3{0, 0, -1}

	switch octantID {
	case 0:
		return TriangleA9{ID: 0, Depth: 0, V0: p0, V1: p1, V2: p2}
	case 1:
		return TriangleA9{ID: 1, Depth: 0, V0: p0, V1: p2, V2: p3}
	case 2:
		return TriangleA9{ID: 2, Depth: 0, V0: p0, V1: p3, V2: p4}
	case 3:
		return TriangleA9{ID: 3, Depth: 0, V0: p0, V1: p4, V2: p1}
	case 4:
		return TriangleA9{ID: 4, Depth: 0, V0: p5, V1: p1, V2: p2}
	case 5:
		return TriangleA9{ID: 5, Depth: 0, V0: p5, V1: p2, V2: p3}
	case 6:
		return TriangleA9{ID: 6, Depth: 0, V0: p5, V1: p3, V2: p4}
	default:
		return TriangleA9{ID: 7, Depth: 0, V0: p5, V1: p4, V2: p1}
	}
}

func CellToCoordinates(base9ID uint64, depth uint64) TriangleA9 {
	var steps []uint64
	val := base9ID
	for i := uint64(0); i < depth; i++ {
		steps = append(steps, val%9)
		val /= 9
	}
	baseOctant := val

	currentTriangle := getBaseOctantTriangle(baseOctant)

	for i := len(steps) - 1; i >= 0; i-- {
		targetChildIdx := steps[i]
		children := currentTriangle.SubdivideA9()
		currentTriangle = children[targetChildIdx]
	}

	return currentTriangle
}

func (v Vector3) DistanceTo(other Vector3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func ToGeoJSON(t TriangleA9) string {
	lat0, lon0 := t.V0.ToLatLng()
	lat1, lon1 := t.V1.ToLatLng()
	lat2, lon2 := t.V2.ToLatLng()

	coordinates := [][][]float64{
		{
			{lon0, lat0},
			{lon1, lat1},
			{lon2, lat2},
			{lon0, lat0},
		},
	}

	geoJSONMap := map[string]interface{}{
		"type": "Feature",
		"properties": map[string]interface{}{
			"cell_id": t.ID,
		},
		"geometry": map[string]interface{}{
			"type":        "Polygon",
			"coordinates": coordinates,
		},
	}

	jsonBytes, err := json.MarshalIndent(geoJSONMap, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func GetBaseOctant(v Vector3) (uint64, TriangleA9) {
	p0 := Vector3{0, 0, 1}
	p1 := Vector3{1, 0, 0}
	p2 := Vector3{0, 1, 0}
	p3 := Vector3{-1, 0, 0}
	p4 := Vector3{0, -1, 0}
	p5 := Vector3{0, 0, -1}

	if v.Z >= 0 {
		if v.X >= 0 && v.Y >= 0 {
			return 0, TriangleA9{ID: 0, Depth: 0, V0: p0, V1: p1, V2: p2}
		}
		if v.X < 0 && v.Y >= 0 {
			return 1, TriangleA9{ID: 1, V0: p0, V1: p2, V2: p3}
		}
		if v.X < 0 && v.Y < 0 {
			return 2, TriangleA9{ID: 2, V0: p0, V1: p3, V2: p4}
		}
		return 3, TriangleA9{ID: 3, V0: p0, V1: p4, V2: p1}
	} else {
		if v.X >= 0 && v.Y >= 0 {
			return 4, TriangleA9{ID: 4, Depth: 0, V0: p5, V1: p1, V2: p2}
		}
		if v.X < 0 && v.Y >= 0 {
			return 5, TriangleA9{ID: 5, Depth: 0, V0: p5, V1: p2, V2: p3}
		}
		if v.X < 0 && v.Y < 0 {
			return 6, TriangleA9{ID: 6, Depth: 0, V0: p5, V1: p3, V2: p4}
		}
		return 7, TriangleA9{ID: 7, Depth: 0, V0: p5, V1: p4, V2: p1}
	}
}

func (id Base9ID) GetNeighbor(dir Direction) Base9ID {
	localIndex := uint64(id) % 9
	parentID := uint64(id) / 9
	switch dir {
	case Edge0:
		if localIndex == 1 {
			return Base9ID(parentID*9 + 0)
		}
		if localIndex == 4 {
			return Base9ID(parentID*9 + 3)
		}
		return Base9ID(uint64(id) + 1)
	case Edge1:
		if localIndex == 1 {
			return Base9ID(parentID*9 + 3)
		}
		return Base9ID(uint64(id) ^ 1)
	case Edge2:
		if localIndex == 1 {
			return Base9ID(parentID*9 + 5)
		}
		return Base9ID(uint64(id) ^ 2)
	}
	return id
}

func RangeQueryBFS(startCell Base9ID, k int) []Base9ID {
	visited := make(map[Base9ID]bool)
	distance := make(map[Base9ID]int)
	queue := []Base9ID{startCell}

	visited[startCell] = true
	distance[startCell] = 0

	var result []Base9ID

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentDist := distance[current]

		if current.IsInscribedCircle() {
			result = append(result, current)
		}

		if currentDist >= k {
			continue
		}

		edges := []Direction{Edge0, Edge1, Edge2}
		for _, edge := range edges {
			neighbor := current.GetNeighbor(edge)
			if !visited[neighbor] {
				visited[neighbor] = true
				distance[neighbor] = currentDist + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

func (t TriangleA9) Center() Vector3 {
	return t.V0.Add(t.V1).Add(t.V2).Normalize()
}

func Geocode(lat, lon float64, targetDepth uint64) Base9ID {
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	targetPos := Vector3{
		X: math.Cos(latRad) * math.Cos(lonRad),
		Y: math.Cos(latRad) * math.Sin(lonRad),
		Z: math.Sin(latRad),
	}

	octantID, currentTriangle := GetBaseOctant(targetPos)
	currentID := Base9ID(octantID)

	for d := uint64(1); d <= targetDepth; d++ {
		children := currentTriangle.SubdivideA9()
		bestChildIdx := 0
		minDist := math.MaxFloat64

		for idx, child := range children {
			dist := targetPos.DistanceTo(child.Center())
			if dist < minDist {
				minDist = dist
				bestChildIdx = idx
			}
		}

		currentID = currentID.PushStep(uint64(bestChildIdx))
		currentTriangle = children[bestChildIdx]
	}

	return currentID
}

type SphericalPolygon []Vector3

func (poly SphericalPolygon) IsPointInPolygon(p Vector3) bool {
	inside := false
	for i := 0; i < len(poly); i++ {
		j := (i + 1) % len(poly)
		if ((poly[i].Y > p.Y) != (poly[j].Y > p.Y)) &&
			(p.X < (poly[j].X-poly[i].X)*(p.Y-poly[i].Y)/(poly[j].Y-poly[i].Y)+poly[i].X) {
			inside = !inside
		}
	}
	return inside
}

func (poly SphericalPolygon) Polyfill(targetDepth uint64) []Base9ID {
	var result []Base9ID
	rootFaces := []TriangleA9{
		getBaseOctantTriangle(0), getBaseOctantTriangle(1),
		getBaseOctantTriangle(2), getBaseOctantTriangle(3),
		getBaseOctantTriangle(4), getBaseOctantTriangle(5),
		getBaseOctantTriangle(6), getBaseOctantTriangle(7),
	}
	for _, root := range rootFaces {
		result = append(result, poly.fillRecursive(root, targetDepth)...)
	}
	return result
}

func (poly SphericalPolygon) fillRecursive(t TriangleA9, targetDepth uint64) []Base9ID {
	cIn := poly.IsPointInPolygon(t.Center())
	v0In := poly.IsPointInPolygon(t.V0)

	if !cIn && !v0In {
		return nil
	}

	if t.Depth == targetDepth {
		id := Base9ID(t.ID)
		if id.IsInscribedCircle() {
			return []Base9ID{id}
		}
		return nil
	}

	children := t.SubdivideA9()
	var results []Base9ID
	for _, child := range children {
		results = append(results, poly.fillRecursive(child, targetDepth)...)
	}
	return results
}
