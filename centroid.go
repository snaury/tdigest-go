package tdigest

// Centroid summarizes Count samples around Mean
type Centroid struct {
	Mean  float64
	Count int64
}

// Add adds a sample x with weight w to the centroid
func (c *Centroid) Add(x float64, w int64) {
	c.Merge(Centroid{
		Mean:  x,
		Count: w,
	})
}

// Merge merges another centroid into this one
func (c *Centroid) Merge(o Centroid) {
	c.Count += o.Count
	c.Mean += float64(o.Count) * (o.Mean - c.Mean) / float64(c.Count)
}

func insertionSort(A []Centroid, N int) {
	// Simple insertion sort shamelessly stolen from wikipedia:
	// https://en.wikipedia.org/wiki/Insertion_sort
	// Also found in:
	// Introduction to Algorithms by Cormen et al.
	for i := 1; i < N; i++ {
		k := A[i]
		j := i - 1
		for j >= 0 && A[j].Mean > k.Mean {
			A[j+1] = A[j]
			j = j - 1
		}
		A[j+1] = k
	}
}
