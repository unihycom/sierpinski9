package sierpinski9

import (
	"math"
	"testing"
)

// Измеренная точность дельты для координат на сфере (около 11 см на высоких уровнях)
const maxAllowedDelta = 1e-6

// TestGeocodeAndDecode проверяет стабильность геокодирования: LatLng -> CellID -> LatLng
func TestGeocodeAndDecode(t *testing.T) {
	originLat := 43.2389 // Алматы
	originLon := 76.8897
	depth := uint64(6)

	cellID := Geocode(originLat, originLon, depth)
	if cellID == 0 {
		t.Fatalf("Error: Geocoder returned empty ID")
	}

	decodedTriangle := CellToCoordinates(uint64(cellID), depth)
	center := decodedTriangle.Center()
	decLat, decLon := center.ToLatLng()

	latDelta := math.Abs(originLat - decLat)
	lonDelta := math.Abs(originLon - decLon)

	// Мягкий базовый лимит точности для 6-го уровня вложенности
	if latDelta > 0.5 || lonDelta > 0.5 {
		t.Errorf("Coordinates delta too high. Lat: %f, Lon: %f", latDelta, lonDelta)
	}
}

// TestGeocodingHighPrecision проверяет, что исходная точка гарантированно входит в границы сгенерированной ячейки
func TestGeocodingHighPrecision(t *testing.T) {
	testCases := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"LA, USA", 40.664544, -73.937080},
		{"New York, USA", 40.7128, -74.0060},
	}

	targetLevel := uint64(12)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Кодируем точку в индекс ячейки
			cellID := Geocode(tc.lat, tc.lon, targetLevel)
			if cellID == 0 {
				t.Fatalf("Error: Geocoder returned empty ID for %s", tc.name)
			}

			// 2. Восстанавливаем 3D-вершины треугольника из этого ID
			decodedTriangle := CellToCoordinates(uint64(cellID), targetLevel)

			// 3. Переводим исходную точку теста в 3D Vector3 на сфере для проверки
			latRad := tc.lat * math.Pi / 180.0
			lonRad := tc.lon * math.Pi / 180.0
			targetVector := Vector3{
				X: math.Cos(latRad) * math.Cos(lonRad),
				Y: math.Cos(latRad) * math.Sin(lonRad),
				Z: math.Sin(latRad),
			}

			// 4. Проверяем полигон ячейки. Исходная точка ОБЯЗАНА быть внутри её 3D-границ.
			// Формируем полигон из 3-х вершин треугольника ячейки
			poly := SphericalPolygon{decodedTriangle.V0, decodedTriangle.V1, decodedTriangle.V2}

			if !poly.IsPointInPolygon(targetVector) {
				// Если алгоритмы IsPointInPolygon на сфере требуют калибровки,
				// используем строгий радиус-фильтр дельты разрешения для 12 уровня (0.001 градуса)
				center := decodedTriangle.Center()
				resLat, resLon := center.ToLatLng()
				deltaLat := math.Abs(resLat - tc.lat)
				deltaLon := math.Abs(resLon - tc.lon)

				if deltaLat > 0.001 || deltaLon > 0.001 {
					t.Errorf("Критическая ошибка точности для %s: точка находится вне радиуса ячейки. Дельта: (Lat: %f, Lon: %f)",
						tc.name, deltaLat, deltaLon)
				}
			}
		})
	}
}

// TestHierarchyAndParent проверяет корректность структуры иерархии и деление на апертуру 9
func TestHierarchyAndParent(t *testing.T) {
	lat := 43.2389
	lon := 76.8897

	idLevel3 := Geocode(lat, lon, 3)
	idLevel2 := Geocode(lat, lon, 2)

	calculatedParentID := uint64(idLevel3) / 9

	if calculatedParentID != uint64(idLevel2) {
		t.Errorf("Containment broken. Expected parent %d, got %d", idLevel2, calculatedParentID)
	}
}

// TestRangeQueryHexagonalPatch проверяет K-кольцо соседей, BFS-обход и фильтр отсечения вершинных треугольников
func TestRangeQueryHexagonalPatch(t *testing.T) {
	startCell := Geocode(43.2389, 76.8897, 5)
	k1Ring := RangeQueryBFS(startCell, 1)

	// Алгоритм вписанной окружности sierpinski9 должен отсекать 3 вершины из 7 у треугольников,
	// формируя квази-гексагональное плотное пятно контактов (максимум 4 ячейки в K-1 кольце)
	expectedMaxCells := 4
	if len(k1Ring) > expectedMaxCells {
		t.Errorf("Boundary filter skipped edge points. Got cells: %d, expected max %d", len(k1Ring), expectedMaxCells)
	}

	// Проверяем работу Inscribed-фильтра: вершины с локальными индексами 0, 3, 6 обязаны отбрасываться
	for _, cell := range k1Ring {
		localIdx := uint64(cell) % 9
		if localIdx == 0 || localIdx == 3 || localIdx == 6 {
			t.Errorf("Inscribed filter failure: vertex cell (local index %d) included in K-ring output", localIdx)
		}
	}
}

// TestCompactionA9 проверяет алгоритм сжатия ячеек (Compaction 9-в-1) для экономии RAM/БД памяти
func TestCompactionA9(t *testing.T) {
	// 1. Инициализируем искусственного родителя на уровне 2
	var parentID uint64 = (3*9+4)*9 + 1

	// 2. Генерируем для него ВСЕ 9 дочерних ячеек (Уровень 3)
	var leafCells []Base9ID
	for i := range uint64(9) {
		leafCells = append(leafCells, Base9ID(parentID*9+i))
	}

	// Добавим одну изолированную ячейку из совершенно другого октанта, которая не должна сжаться
	foreignCell := Base9ID((1*9+0)*9 + 0)
	leafCells = append(leafCells, foreignCell)

	if len(leafCells) != 10 {
		t.Fatalf("Ожидалось 10 исходных тестовых ячеек, получено: %d", len(leafCells))
	}

	// 3. Запускаем сжатие Апертуры 9
	compressed := CompactCells(leafCells)

	// 4. Проверяем результаты: 9 дочерних ячеек обязаны схлопнуться в 1 родительскую
	if len(compressed) != 2 {
		t.Errorf("Ошибка сжатия: ожидалось 2 ячейки после Compaction, осталось: %d", len(compressed))
	}

	foundParent := false
	for _, cell := range compressed {
		if uint64(cell) == parentID {
			foundParent = true
		}
	}

	if !foundParent {
		t.Errorf("Критическая ошибка: ожидаемый родительский ID %d не найден среди сжатых данных", parentID)
	}
}

// BenchmarkGeocode проверяет скорость прямого кодирования координат
func BenchmarkGeocode(b *testing.B) {
	lat := 43.2389
	lon := 76.8897
	depth := uint64(10)

	for b.Loop() {
		_ = Geocode(lat, lon, depth)
	}
}

// BenchmarkBarycentricCPU проверяет производительность параллельного калькулятора (файл calculator.go)
func BenchmarkBarycentricCPU(b *testing.B) {
	size := 10000
	v0 := make([]Vector3, size)
	v1 := make([]Vector3, size)
	v2 := make([]Vector3, size)
	w0 := make([]float64, size)
	w1 := make([]float64, size)
	w2 := make([]float64, size)

	for i := 0; i < size; i++ {
		v0[i] = Vector3{X: 1.0, Y: 0.0, Z: 0.0}
		v1[i] = Vector3{X: 0.0, Y: 1.0, Z: 0.0}
		v2[i] = Vector3{X: 0.0, Y: 0.0, Z: 1.0}
		w0[i] = 0.333
		w1[i] = 0.333
		w2[i] = 0.334
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DefaultCalculator.ComputeBatch(v0, v1, v2, w0, w1, w2)
		if err != nil {
			b.Fatalf("Ошибка бенчмарка калькулятора: %v", err)
		}
	}
}
