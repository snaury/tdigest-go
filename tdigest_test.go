package tdigest

import (
	"math/rand"
	"testing"
)

func TestDigest10(t *testing.T) {
	digest := New(100)
	for i := 1; i <= 10; i++ {
		digest.Add(float64(i), 1)
	}
	digest.Compress()
	t.Logf("digest.summary = %v", digest.summary)
	if digest.count != 10 {
		t.Errorf("digest.count = %d", digest.count)
	}
	cases := []struct {
		q float64
		v float64
	}{
		{0.0, 1.0},
		{0.1, 1.0},
		{0.5, 5.0},
		{0.9, 9.0},
		{1.0, 10.0},
	}
	for _, c := range cases {
		v := digest.Quantile(c.q)
		if v != c.v {
			t.Errorf("digest.Quantile(%f) = %f (expected %f)", c.q, v, c.v)
		}
	}
}

func TestDigest1000000(t *testing.T) {
	digest := New(100)
	for i := 1; i <= 1000000; i++ {
		digest.Add(float64(i), 1)
	}
	digest.Compress()
	t.Logf("digest.summary = %v", digest.summary)
	if digest.count != 1000000 {
		t.Errorf("digest.count = %d", digest.count)
	}
	cases := []struct {
		q float64
		v float64
	}{
		{0.0, 1.0},
		{0.1, 100000.0},
		{0.5, 500000.0},
		{0.9, 900000.0},
		{1.0, 1000000.0},
	}
	for _, c := range cases {
		v := digest.Quantile(c.q)
		if v != c.v {
			t.Errorf("digest.Quantile(%f) = %f (expected %f)", c.q, v, c.v)
		}
	}
}

func BenchmarkDigest(b *testing.B) {
	values := make([]float64, b.N)
	for i := range values {
		values[i] = rand.Float64()
	}
	digest := New(100)
	b.ResetTimer()
	for _, value := range values {
		digest.Add(value, 1)
	}
	digest.Compress()
}
