package rider

import (
	"testing"
)

func TestCreatePowerZones(t *testing.T) {
	rider := RIDER{Attributes: []RIDER_ATTRIBUTES{{FTP: 100}}}
	createPowerZones(&rider)
	if len(rider.Attributes[0].PowerZones) != 7 {
		t.Error("Wrong number of zones")
	}
	pz := rider.Attributes[0].PowerZones
	expectations := [][]uint32{{0, 20}, {21, 50}, {51, 70}, {71, 85}, {86, 100}, {101, 115}, {116, 2000}}
	for idx, expectation := range expectations {
		if pz[idx].Min != expectation[0] || pz[idx].Max != expectation[1] {
			t.Errorf("Power Zone %d not %d to %d : %d - %d", idx, expectation[0], expectation[1], pz[idx].Min, pz[idx].Max)
		}
	}
}

func TestCreateHRZones(t *testing.T) {
}
