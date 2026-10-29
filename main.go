package main

import (
	"encoding/json"
	"fmt"
	"math"
)

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

type TriangleA9 struct {
	ID         uint64
	V0, V1, V2 Vector3
}

func InterpSphere(a, b Vector3, weight float64) Vector3 {
	return Vector3{
		X: a.X*(1.0-weight) + b.X*weight,
		Y: a.Y*(1.0-weight) + b.Y*weight,
		Z: a.Z*(1.0-weight) + b.Z*weight,
	}.Normalize()
}

// SubdivideA9 генерирует 9 дочерних ячеек (используем точную геометрию из прошлых шагов)
func (t TriangleA9) SubdivideA9() []TriangleA9 {
	p01_1 := InterpSphere(t.V0, t.V1, 1.0/3.0)
	p01_2 := InterpSphere(t.V0, t.V1, 2.0/3.0)
	p12_1 := InterpSphere(t.V1, t.V2, 1.0/3.0)
	p12_2 := InterpSphere(t.V1, t.V2, 2.0/3.0)
	p20_1 := InterpSphere(t.V2, t.V0, 1.0/3.0)
	p20_2 := InterpSphere(t.V2, t.V0, 2.0/3.0)

	c0 := InterpSphere(p01_1, p12_2, 0.5)
	c1 := InterpSphere(p01_2, p20_1, 0.5)
	c2 := InterpSphere(p12_1, p20_2, 0.5)

	baseID := t.ID

	return []TriangleA9{
		{ID: baseID*9 + 0, V0: t.V0, V1: p01_1, V2: p20_2},
		{ID: baseID*9 + 1, V0: p01_1, V1: p01_2, V2: c1},
		{ID: baseID*9 + 2, V0: p01_2, V1: t.V1, V2: p12_1},
		{ID: baseID*9 + 3, V0: p01_1, V1: c1, V2: p20_2},
		{ID: baseID*9 + 4, V0: p20_2, V1: c1, V2: c0},
		{ID: baseID*9 + 5, V0: p01_2, V1: p12_1, V2: c2},
		{ID: baseID*9 + 6, V0: c1, V1: p12_1, V2: c2},
		{ID: baseID*9 + 7, V0: p20_2, V1: c0, V2: p20_1},
		{ID: baseID*9 + 8, V0: p20_1, V1: p12_2, V2: t.V2},
	}
}

// Инициализация базовой грани октаэдра (0-7)
func getBaseOctantTriangle(octantID uint64) TriangleA9 {
	p0 := Vector3{0, 0, 1}  // Северный полюс
	p1 := Vector3{1, 0, 0}  // 0°
	p2 := Vector3{0, 1, 0}  // 90°
	p3 := Vector3{-1, 0, 0} // 180°
	p4 := Vector3{0, -1, 0} // 270°
	p5 := Vector3{0, 0, -1} // Южный полюс

	switch octantID {
	case 0:
		return TriangleA9{ID: 0, V0: p0, V1: p1, V2: p2}
	case 1:
		return TriangleA9{ID: 1, V0: p0, V1: p2, V2: p3}
	case 2:
		return TriangleA9{ID: 2, V0: p0, V1: p3, V2: p4}
	case 3:
		return TriangleA9{ID: 3, V0: p0, V1: p4, V2: p1}
	case 4:
		return TriangleA9{ID: 4, V0: p5, V1: p1, V2: p2}
	case 5:
		return TriangleA9{ID: 5, V0: p5, V1: p2, V2: p3}
	case 6:
		return TriangleA9{ID: 6, V0: p5, V1: p3, V2: p4}
	default:
		return TriangleA9{ID: 7, V0: p5, V1: p4, V2: p1}
	}
}

// CellToCoordinates декодирует Base-9 ID и возвращает точную геометрию треугольника
func CellToCoordinates(base9ID uint64) TriangleA9 {
	if base9ID <= 7 {
		return getBaseOctantTriangle(base9ID)
	}

	// 1. Разбираем число на шаги фрактала (от конца к началу)
	var steps []uint64
	val := base9ID
	for val > 7 {
		steps = append(steps, val%9)
		val /= 9
	}
	baseOctant := val // Оставшееся число от 0 до 7 — наш стартовый октант

	// 2. Инициализируем корневую грань октаэдра
	currentTriangle := getBaseOctantTriangle(baseOctant)

	// 3. Воспроизводим спуск по иерархии (разворачиваем массив шагов в правильном порядке)
	for i := len(steps) - 1; i >= 0; i-- {
		targetChildIdx := steps[i]
		children := currentTriangle.SubdivideA9()
		currentTriangle = children[targetChildIdx]
	}

	return currentTriangle
}

// ToGeoJSON конвертирует геометрию треугольника в валидную строку GeoJSON (Polygon)
func ToGeoJSON(t TriangleA9) string {
	// Переводим вершины из 3D в Lat/Lng градусы
	lat0, lon0 := t.V0.ToLatLng()
	lat1, lon1 := t.V1.ToLatLng()
	lat2, lon2 := t.V2.ToLatLng()

	// В стандарте GeoJSON координаты пишутся в формате [Долгота, Широта] (Lng, Lat)!
	// Полигон должен быть замкнут, поэтому первая и последняя точки совпадают
	coordinates := [][][]float64{
		{
			{lon0, lat0},
			{lon1, lat1},
			{lon2, lat2},
			{lon0, lat0},
		},
	}

	// Структура для сериализации в JSON
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

func main() {
	// Возьмем тестовый Base-9 ID (например, сгенерированный геокодером на прошлом шаге)
	// Октант 3 -> Шаг 1 -> Шаг 4 -> Шаг 6 -> Шаг 7
	var testID uint64 = (((3*9+1)*9+4)*9+6)*9 + 7

	fmt.Printf("Исходный лаконичный Base-9 ID: %d\n", testID)
	fmt.Println("Запуск декодера CellToCoordinates...")

	// 1. Декодируем ID обратно в сферические вершины треугольника
	decodedTriangle := CellToCoordinates(testID)

	// 2. Генерируем GeoJSON для вывода на карту
	geoJSONStr := ToGeoJSON(decodedTriangle)

	fmt.Println("\n=== СГЕНЕРИРОВАННЫЙ GEOJSON ПОЛИГОН ===")
	fmt.Println(geoJSONStr)
}
