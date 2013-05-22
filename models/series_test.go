package models

import "testing"

func TestRange(t *testing.T) {

	s := Series{1,2,3,4}
	p := pair{1,4}

	if s.Range() != p {
		t.Error("Range was wrong")
	}
}

func TestQuantile(t *testing.T) {
	s := Series{1,2,3,4}

	q := map[float64] float64 {
		0.0: 1,
		0.5: 2.5,
		0.6: 2.8,
		1.0: 4,
	}

	for position, quantile := range q {
		if s.Quantile(position) != quantile {
			t.Errorf("s.Quantile(%v): %v != %v", position, s.Quantile(position), quantile);
		}
	}
}

func TestVariance(t *testing.T) {
	s := Series{1,2,3,4,5,6,7,8,9}

	if s.Variance() != 7.5 {
		t.Errorf("s.Variance(): %v != 7.5", s.Variance());
	}
}

func TestBandwidth(t *testing.T) {
	s := Series{1,2,3,4,5,6,7,8,9}

	if s.KernelBandwidth() != 1.8706304309991295 {
		t.Error("Bandwidth was wrong")
	}
}
