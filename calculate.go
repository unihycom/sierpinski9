package sierpinski9

import (
	"math"
	"sync"
)

// Константа инверсии тройки (1.0 / 3.0) исключает медленную операцию деления на CPU.
const invThree = 1.0 / 3.0

// Vector3 представляет собой трехмерный вектор для сферических координат библиотеки sierpinski9.
type Vector3 struct {
	X, Y, Z float64
}

// Normalize нормализует вектор, приводя его к единичной длине на сфере.
func (v Vector3) Normalize() Vector3 {
	len := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if len == 0 {
		return Vector3{}
	}
	return Vector3{X: v.X / len, Y: v.Y / len, Z: v.Z / len}
}

// ToLatLng переводит 3D-вектор единичной сферы обратно в географические градусы (Широта, Долгота).
func (v Vector3) ToLatLng() (float64, float64) {
	latRad := math.Asin(v.Z)
	lonRad := math.Atan2(v.Y, v.X)

	latDeg := latRad * 180.0 / math.Pi
	lonDeg := lonRad * 180.0 / math.Pi

	return latDeg, lonDeg
}

// Add складывает два трехмерных вектора.
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// DistanceTo вычисляет евклидово расстояние между двумя 3D точками на сфере.
func (v Vector3) DistanceTo(other Vector3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// BarycentricCalculator описывает стратегию параллельного расчета барицентрических координат.
// Разработчики высоконагруженных ГИС, корпоративных кластеров и 3D-движков могут реализовать
// этот интерфейс для внешних векторных мощностей (wgpu-go, CUDA, Metal, OpenCL).
type BarycentricCalculator interface {
	ComputeBatch(v0, v1, v2 []Vector3, w0, w1, w2 []float64) ([]Vector3, error)
}

// CPUMultiThreadCalculator — стандартный высокопроизводительный потоковый вычислитель на CPU.
type CPUMultiThreadCalculator struct {
	WorkerCount int
}

// ComputeBatch параллельно распределяет барицентрические вычисления по ядрам процессора.
func (c *CPUMultiThreadCalculator) ComputeBatch(v0, v1, v2 []Vector3, w0, w1, w2 []float64) ([]Vector3, error) {
	length := len(v0)
	results := make([]Vector3, length)

	workers := c.WorkerCount
	if workers <= 0 {
		workers = 4 // Значение по умолчанию, или runtime.NumCPU()
	}

	chunkSize := (length + workers - 1) / workers
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= length {
			break
		}
		if end > length {
			end = length
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			for j := s; j < e; j++ {
				results[j] = Vector3{
					X: (v0[j].X*w0[j] + v1[j].X*w1[j] + v2[j].X*w2[j]) * invThree,
					Y: (v0[j].Y*w0[j] + v1[j].Y*w1[j] + v2[j].Y*w2[j]) * invThree,
					Z: (v0[j].Z*w0[j] + v1[j].Z*w1[j] + v2[j].Z*w2[j]) * invThree,
				}.Normalize()
			}
		}(start, end)
	}

	wg.Wait()
	return results, nil
}

// DefaultCalculator хранит текущую активную стратегию вычислений библиотеки sierpinski9.
var DefaultCalculator BarycentricCalculator = &CPUMultiThreadCalculator{WorkerCount: 4}

// SetCalculator позволяет динамически подменить вычислительный движок на сторонний GPU/TPU-кластер.
func SetCalculator(calc BarycentricCalculator) {
	DefaultCalculator = calc
}
