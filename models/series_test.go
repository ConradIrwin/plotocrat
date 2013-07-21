package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	s := Series{1, 2, 3, 4}
	p := pair{1, 4}

	assert.Equal(t, p, s.Range())
}

func TestQuantile(t *testing.T) {
	s := Series{1, 2, 3, 4}

	q := map[float64]float64{
		0.0: 1,
		0.5: 2.5,
		0.6: 2.8,
		1.0: 4,
	}

	for position, quantile := range q {
		assert.Equal(t, quantile, s.Quantile(position))
	}
}

func TestVariance(t *testing.T) {
	s := Series{1, 2, 3, 4, 5, 6, 7, 8, 9}

	assert.Equal(t, 7.5, s.Variance())
}

func TestBandwidth(t *testing.T) {
	s := Series{1, 2, 3, 4, 5, 6, 7, 8, 9}

	assert.Equal(t, 1.8706304309991295, s.KernelBandwidth())
}
