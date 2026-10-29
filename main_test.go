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

	for i := 0; i < b.N; i++ {
		_ = Geocode(lat, lon, depth)
	}
}
