package models

import (
	"bytes"
	svgo "github.com/ajstarks/svgo"
	"io"
	"strconv"
)

type Visualization struct {
	series Series
	width  int
	height int
	margin int
}

type pair [2]float64

type Scaler interface {
	Scale(float64) float64
}

type scale struct {
	Domain pair
	Range  pair
}

func NewVisualization(plot *Plot) *Visualization {
	series := NewSeries(plot.Values())
	return &Visualization{series: series,
		width:  640,
		height: 480,
		margin: 40}
}

func (viz *Visualization) WriteTo(w io.Writer) error {
	svg := svgo.New(w)
	svg.Start(viz.width, viz.height)

	svg.Path(viz.cumulativeProbabilityData(), "fill:none; stroke:red")

	svg.End()
	return nil
}

func (viz *Visualization) cumulativeProbabilityData() string {
	var output bytes.Buffer

	xScale := scale{Domain: viz.series.Range(), Range: viz.Domain()}
	yScale := scale{Domain: viz.series.Domain(), Range: viz.Range()}

	for y, x := range viz.series {
		if y == 0 {
			output.WriteRune('M')
		} else {
			output.WriteRune('L')
		}

		output.WriteString(str(xScale.Scale(x)))
		output.WriteRune(',')
		output.WriteString(str(yScale.Scale(float64(y))))
	}

	return output.String()
}

func str(f float64) string {
	return strconv.FormatFloat(f, 'G', 6, 64)
}

func (viz *Visualization) Domain() pair {
	return pair{float64(viz.margin), float64(viz.width - viz.margin)}
}

func (viz *Visualization) Range() pair {
	return pair{float64(viz.height - viz.margin), float64(viz.margin)}
}

func (scale scale) Scale(x float64) float64 {
	return (x-scale.Domain[0])*scale.Range.Size()/scale.Domain.Size() + scale.Range[0]
}

func (pair pair) Size() float64 {
	return pair[1] - pair[0]
}
