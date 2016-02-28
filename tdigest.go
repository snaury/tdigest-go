package tdigest

import (
	"math"
)

// Default maximum number of unmerged items
const DefaultMaxUnmerged = 512

// MergingDigest amortizes computation by merging in fixed sized batches
type MergingDigest struct {
	merged      []Centroid
	summary     []Centroid
	unmerged    []Centroid
	count       int64
	compression float64
	MaxUnmerged int
}

// New creates a new MergingDigest with the given compression
func New(compression float64) *MergingDigest {
	return &MergingDigest{
		compression: compression,
		MaxUnmerged: DefaultMaxUnmerged,
	}
}

// Add adds a sample x with weight w
func (digest *MergingDigest) Add(x float64, w int64) {
	if w <= 0 {
		panic("Cannot add samples with non-positive weight")
	}
	digest.Merge(Centroid{
		Mean:  x,
		Count: w,
	})
}

// Merge merges a centroid into the digest
func (digest *MergingDigest) Merge(c Centroid) {
	if c.Count <= 0 {
		panic("Cannot merge centroids with non-positive count")
	}
	digest.unmerged = append(digest.unmerged, c)
	digest.count += c.Count
	if len(digest.unmerged) >= digest.MaxUnmerged {
		digest.Compress()
	}
}

func (digest *MergingDigest) collapse(sum int64, m int, b Centroid) (int64, int) {
	if len(digest.merged) == 0 {
		digest.merged = append(digest.merged, b)
		return sum, 0
	}
	a := digest.merged[m]
	qa := (float64(sum) + float64(a.Count-1)*0.5) / float64(digest.count-1)
	err := qa * (1.0 - qa)
	qb := (float64(sum+a.Count) + float64(b.Count-1)*0.5) / float64(digest.count-1)
	err2 := qb * (1.0 - qb)
	if err > err2 {
		err = err2
	}
	k := 4 * float64(digest.count) * err / digest.compression
	if float64(a.Count+b.Count) <= k {
		digest.merged[m].Merge(b)
		return sum, m
	}
	digest.merged = append(digest.merged, b)
	return sum + a.Count, m + 1
}

// Compress merges any unmerged data into the summary
func (digest *MergingDigest) Compress() {
	if len(digest.unmerged) == 0 {
		return
	}
	stableSort(digest.unmerged)
	sum := int64(0)
	m := 0
	i := 0
	j := 0
	digest.merged = digest.merged[:0]
	for i < len(digest.summary) && j < len(digest.unmerged) {
		if digest.summary[i].Mean <= digest.unmerged[j].Mean {
			sum, m = digest.collapse(sum, m, digest.summary[i])
			i++
		} else {
			sum, m = digest.collapse(sum, m, digest.unmerged[j])
			j++
		}
	}
	for i < len(digest.summary) {
		sum, m = digest.collapse(sum, m, digest.summary[i])
		i++
	}
	for j < len(digest.unmerged) {
		sum, m = digest.collapse(sum, m, digest.unmerged[j])
		j++
	}
	digest.merged, digest.summary = digest.summary, digest.merged
	digest.unmerged = digest.unmerged[:0]
}

// Summary returns compressed summary of the digest
func (digest *MergingDigest) Summary() []Centroid {
	digest.Compress()
	summary := make([]Centroid, len(digest.summary))
	copy(summary, digest.summary)
	return summary
}

// Quantile returns an estimate of the value at quantile q
func (digest *MergingDigest) Quantile(q float64) float64 {
	digest.Compress()
	if len(digest.summary) == 0 || q < 0.0 || q > 1.0 {
		return math.NaN()
	} else if len(digest.summary) == 1 {
		return digest.summary[0].Mean
	}
	index := float64(digest.count) * q
	sum := int64(1)
	aMean := digest.summary[0].Mean
	aIndex := float64(0)
	bMean := digest.summary[0].Mean
	bIndex := float64(sum) + float64(digest.summary[0].Count-1)*0.5
	for i := 1; i < len(digest.summary); i++ {
		if index <= bIndex {
			break
		}
		sum += digest.summary[i-1].Count
		aMean = bMean
		aIndex = bIndex
		bMean = digest.summary[i].Mean
		bIndex = float64(sum) + float64(digest.summary[i].Count-1)*0.5
	}
	p := (index - aIndex) / (bIndex - aIndex)
	return aMean*(1.0-p) + bMean*p
}
