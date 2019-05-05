package gokmv

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestKMVSize(t *testing.T) {
	kmv := NewKMV(2)

	if kmv.Size() != 0 {
		t.Error("An empty table should have size 0")
	}

	kmv.InsertUint64(1)
	if kmv.Size() != 1 {
		t.Error("We only added one element")
	}

	for i := 0; i < 10; i++ {
		kmv.InsertUint64(uint64(i))
	}

	if kmv.Size() != 2 {
		t.Error("Size should not be bigger than maxSize")
	}
}

func TestKMVElementsAdded(t *testing.T) {
	kmv := NewKMV(2)

	if kmv.ElementsAdded() != 0 {
		t.Error("We did not add any element")
	}

	kmv.InsertUint64(1)
	if kmv.ElementsAdded() != 1 {
		t.Error("We added one element")
	}

	for i := 0; i < 10; i++ {
		kmv.InsertUint64(uint64(i))
	}

	if kmv.ElementsAdded() != 11 {
		t.Error("We added 11 elements")
	}
}

func TestKMVInsertString(t *testing.T) {
	kmv := NewKMV(2)

	kmv.InsertString("Golang")

	if kmv.ElementsAdded() != 1 || kmv.Size() != 1 {
		t.Error("We added element one")
	}
}

func inBounds(relativeError float64, approximation, real int) bool {
	fApprox := float64(approximation)
	fReal := float64(real)

	if fApprox < (1-relativeError)*fReal {
		return false
	}

	if fApprox > (1+relativeError)*fReal {
		return false
	}

	return true
}
func TestKMVEstimateCardinality(t *testing.T) {
	data := make(map[uint64]bool)
	dataSize := 1000000

	rand.Seed(42)
	for len(data) != dataSize {
		n := rand.Uint64()
		data[n] = true
	}

	// We have a sample of `dataSize` random uint64
	// We will estimate the carinality of the sample
	// `iterations` times and check that the `avgEstimation` is
	// not off by more than a factor of `relativeErr`
	avgEstimation := 0
	avgSize := 0
	iterations := 10
	for i := 0; i < iterations; i++ {
		kmv := NewKMV(64)

		for key := range data {
			kmv.InsertString(fmt.Sprint(key))
		}
		avgEstimation += int(kmv.EstimateCardinality())
		avgSize += kmv.Size()
	}
	avgEstimation /= iterations
	avgSize /= iterations

	relativeErr := 0.04

	if !inBounds(relativeErr, avgEstimation, dataSize) {
		errMsg := fmt.Sprintf("The kmv estimation was not in the theoretical bounds: %d out of %d", avgEstimation, dataSize)
		t.Error(errMsg)
	}
}
