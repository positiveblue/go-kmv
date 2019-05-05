package gokmv

import (
	"math"
	"math/rand"
	"time"

	adaptivetable "github.com/positiveblue/adaptive-table"
	murmur3 "github.com/spaolacci/murmur3"
)

type KMV struct {
	table        adaptivetable.AdaptiveTable
	initialSize  int
	seed         uint32
	totalCounter uint64
}

func NewKMV(size int) *KMV {
	rand.Seed(time.Now().UnixNano())
	return NewKMVWithSeed(size, rand.Uint32())
}

func NewKMVWithSeed(size int, seed uint32) *KMV {
	return &KMV{
		table:        adaptivetable.NewAdaptiveTableComplete(size, math.MaxInt64, size),
		initialSize:  size,
		seed:         seed,
		totalCounter: 0,
	}
}

func (kmv *KMV) ElementsAdded() uint64 {
	return kmv.totalCounter
}

func (kmv *KMV) Size() int {
	return kmv.table.Size()
}

func (kmv *KMV) InsertUint64(hash uint64) {
	kmv.totalCounter++
	kmv.table.Insert(hash)
}

func (kmv *KMV) InsertString(s string) {
	hash := murmur3.Sum64WithSeed([]byte(s), kmv.seed)
	kmv.InsertUint64(hash)
}

func (kmv *KMV) EstimateCardinality() uint64 {
	if kmv.Size() < kmv.initialSize {
		return uint64(kmv.table.Size())
	}

	meanDistance := kmv.table.Max() / uint64(kmv.table.Size())
	return uint64(math.MaxUint64 / meanDistance)
}
