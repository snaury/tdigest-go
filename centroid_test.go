package tdigest

import (
	"sort"
	"testing"
)

type sortByMean []Centroid

func (s sortByMean) Len() int           { return len(s) }
func (s sortByMean) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByMean) Less(i, j int) bool { return s[i].Mean < s[j].Mean }

func TestInsertionSort(t *testing.T) {
	Centroids1 := []Centroid{
		Centroid{Mean: 10, Count: 1},
		Centroid{Mean: 11, Count: 2},
		Centroid{Mean: 12, Count: 3},
		Centroid{Mean: 10, Count: 4},
		Centroid{Mean: 16, Count: 5},
		Centroid{Mean: 88, Count: 6},
		Centroid{Mean: 17, Count: 7},
		Centroid{Mean: 19, Count: 8},
		Centroid{Mean: 99, Count: 9},
	}
	N := len(Centroids1)
	Centroids2 := make([]Centroid, N)
	copy(Centroids2, Centroids1)

	insertionSort(Centroids1, len(Centroids1))
	sort.Stable(sortByMean(Centroids2))

	// Check that we get the same sorting
	for i := 0; i < N; i++ {
		if Centroids1[i] != Centroids2[i] {
			t.Error("Mismatch at index: ", i)
		}
	}
}
