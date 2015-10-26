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
