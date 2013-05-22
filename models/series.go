package models

import (
	stats "github.com/GaryBoone/GoStats/stats"
	"math"
	"sort"
)

type Series []float64

func NewSeries(n []float64) Series {
	series := n
	sort.Float64s(series)
	return series
}

func (series Series) Range() pair {
	return pair{series[0], series[len(series)-1]}
}

func (series Series) Domain() pair {
	return pair{0, float64(len(series)) - 1}
}

func (series Series) Variance() float64 {
	return stats.StatsSampleVariance(series)
}

// Scott, D. W. (1992) Multivariate Density Estimation: Theory, Practice, and
// Visualization. Wiley. (via science.js, via R)
func (series Series) KernelBandwidth() float64 {
	h := series.InterQuartileRange() / 1.34

	sigma := math.Sqrt(series.Variance())
	if sigma < h {
		h = sigma
	}

	return 1.06 * h * math.Pow(series.length(), -1.0/5.0)
}

func (series Series) InterQuartileRange() float64 {
	return series.Quantile(.75) - series.Quantile(.25)
}

// R's quantile algorithm type=7 (via science.js)
func (series Series) Quantile(position float64) float64 {

	if position <= 0 {
		return series[0]
	}
	if position >= 1 {
		return series.last()
	}

	index := 1 + position*float64(len(series)-1)
	lo := math.Floor(index)
	h := index - lo
	a := series.at(lo - 1)
	b := series.at(lo)

	return a + h*(b-a)
}

func (series Series) last() float64 {
	return series[len(series)-1]
}

func (series Series) length() float64 {
	return float64(len(series))
}

func (series Series) at(i float64) float64 {
	return series[int(i)]
}
