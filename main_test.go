package sierpinski9

import (
	"math"
	"testing"
)

func TestGeocodeAndDecode(t *testing.T) {
	originLat := 43.2389
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

	if latDelta > 0.5 || lonDelta > 0.5 {
		t.Errorf("Coordinates delta too high. Lat: %f, Lon: %f", latDelta, lonDelta)
	}
}

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

func TestRangeQueryHexagonalPatch(t *testing.T) {
	startCell := Geocode(43.2389, 76.8897, 5)
	k1Ring := RangeQueryBFS(startCell, 1)

	expectedMaxCells := 4
	if len(k1Ring) > expectedMaxCells {
		t.Errorf("Boundary filter skipped edge points. Got cells: %d, expected max %d", len(k1Ring), expectedMaxCells)
	}

	for _, cell := range k1Ring {
		localIdx := uint64(cell) % 9
		if localIdx == 0 || localIdx == 3 || localIdx == 6 {
			t.Errorf("Inscribed filter failure: vertex cell %d included in K-ring output", localIdx)
		}
	}
}

func BenchmarkGeocode(b *testing.B) {
	lat := 43.2389
	lon := 76.8897
	depth := uint64(10)

	for b.Loop() {
		_ = Geocode(lat, lon, depth)
	}
}
func TestCompactionA9(t *testing.T) {
	// 1. Инициализируем искусственного родителя на уровне 2
	// Октант 3 -> Шаг 4 -> Шаг 1
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

	// 4. Проверяем результаты
	// 9 дочерних ячеек должны превратиться в 1 родительскую.
	// Плюс 1 посторонняя ячейка. Итого в массиве должно остаться строго 2 элемента.
	if len(compressed) != 2 {
		t.Errorf("Ошибка сжатия: ожидалось 2 ячейки после Compaction, осталось: %d", len(compressed))
	}

	// Проверяем, что родительский ID восстановился корректно
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
