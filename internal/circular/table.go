package circular

import (
	"fmt"
	"math"
)

type Table struct {
	ZVals []float64   `json:"z_vals"`
	RVals []float64   `json:"r_vals"`
	Data  [][]float64 `json:"data"`
}

func (t Table) Validate() error {
	if len(t.ZVals) == 0 || len(t.RVals) == 0 {
		return fmt.Errorf("table has no axes")
	}
	if len(t.Data) != len(t.RVals) {
		return fmt.Errorf("table rows %d do not match r axis %d", len(t.Data), len(t.RVals))
	}
	for i, row := range t.Data {
		if len(row) != len(t.ZVals) {
			return fmt.Errorf("row %d length %d does not match z axis %d", i, len(row), len(t.ZVals))
		}
		for _, value := range row {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return fmt.Errorf("table contains non-finite value")
			}
		}
	}
	return nil
}

func (t Table) Influence(zRatio, rRatio float64) (float64, error) {
	if err := t.Validate(); err != nil {
		return 0, err
	}
	zi := indexFor(t.ZVals, zRatio)
	ri := indexFor(t.RVals, rRatio)
	z0, z1 := t.ZVals[zi[0]], t.ZVals[zi[1]]
	r0, r1 := t.RVals[ri[0]], t.RVals[ri[1]]
	dz := 0.0
	if z1 > z0 {
		dz = (zRatio - z0) / (z1 - z0)
	}
	dr := 0.0
	if r1 > r0 {
		dr = (rRatio - r0) / (r1 - r0)
	}
	f00 := t.Data[ri[0]][zi[0]]
	f10 := t.Data[ri[0]][zi[1]]
	f01 := t.Data[ri[1]][zi[0]]
	f11 := t.Data[ri[1]][zi[1]]
	return f00*(1-dz)*(1-dr) + f10*dz*(1-dr) + f01*(1-dz)*dr + f11*dz*dr, nil
}

func indexFor(values []float64, target float64) [2]int {
	if len(values) == 1 {
		return [2]int{0, 0}
	}
	if target <= values[0] {
		return [2]int{0, 1}
	}
	last := len(values) - 1
	if target >= values[last] {
		return [2]int{last - 1, last}
	}
	right := 1
	for right < last && values[right] < target {
		right++
	}
	return [2]int{right - 1, right}
}

func (t Table) Snapshot() map[string]interface{} {
	return map[string]interface{}{
		"z_axis": t.ZVals,
		"r_axis": t.RVals,
		"rows":   len(t.Data),
		"cols": func() int {
			if len(t.Data) > 0 {
				return len(t.Data[0])
			}
			return 0
		}(),
	}
}

func (t Table) Centerline(zRatio float64) (float64, error) {
	return t.Influence(zRatio, 0)
}

func (t Table) Surface(rRatio float64) (float64, error) {
	return t.Influence(0, rRatio)
}

func (t Table) MaxValue() float64 {
	max := 0.0
	for _, row := range t.Data {
		for _, value := range row {
			if value > max {
				max = value
			}
		}
	}
	return max
}

func (t Table) Size() (int, int) {
	return len(t.ZVals), len(t.RVals)
}

func (t Table) At(zIndex, rIndex int) (float64, error) {
	if zIndex < 0 || zIndex >= len(t.ZVals) || rIndex < 0 || rIndex >= len(t.RVals) {
		return 0, fmt.Errorf("index out of range")
	}
	return t.Data[rIndex][zIndex], nil
}
