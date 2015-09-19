package tessellated

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"time"
)

const (
	base   = 100.0
	height = base * 1.732050808 * 0.5 // base * (sqrt(3)/2)
	radius = base * 0.35              // should be strictly smaller than (base * 0.5)
	circle = 2.0 * math.Pi

	ystart = 0.0 - height
	xo     = 0.0 - base     // Where every odd row starts in the x direction.
	xe     = 0.0 - base*0.5 // Where every even row starts in the x direction.

	maxopacity = 0.15

	svgprefix   = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="%0.0fpx" height="%0.0fpx" viewBox="0 0 %0.0f %0.0f" zoomAndPan="disable"><defs><style type="text/css">polygon { fill: #090909 }</style></defs>`
	svgtriangle = `<polygon points="%v %v %v" opacity="%0.3f" />`
	svgsuffix   = `</svg>`
)

var random *rand.Rand = rand.New(rand.NewSource(time.Now().Unix()))

type Rect struct {
	Width, Height float64
}

type point struct {
	X, Y float64
}

func (p *point) randomize() {
	hypot := radius * random.Float64()
	theta := circle * random.Float64()
	p.X += hypot * math.Cos(theta)
	p.Y += hypot * math.Sin(theta)
}

func (p point) GoString() string {
	return p.String()
}

func (p point) String() string {
	return fmt.Sprintf("%0.2f,%0.2f", p.X, p.Y)
}

type row []point

func reset(x, y float64, r row) {
	n := len(r)
	for i := 0; i < n; x, i = x+base, i+1 {
		p := &r[i]
		p.X = x
		p.Y = y
		p.randomize()
	}
}

func Triangle(size Rect, w io.Writer) {
	fmt.Fprintf(w, svgprefix, size.Width, size.Height, size.Width, size.Height)
	nr := 3 + int(math.Ceil(size.Height/height)) // Number of rows of points
	npo := 3 + int(math.Ceil(size.Width/base))   // Number of points in an odd row
	npe := 2 + int(math.Ceil(size.Width/base))   // Number of points in an even row

	// Allocates even and odd rows only once per function call! Could use a sync.Pool
	o := make(row, npo)
	reset(xo, ystart, o)
	e := make(row, npe)

	// Construct the triangles
	for y, i := 0.0, 1; i < nr; y, i = y+height, i+1 {
		if 0 == i%2 {
			// Odd row
			reset(xo, y, o)
			for j := 0; j < npe; j++ {
				// Pointing up
				opacity := maxopacity * random.Float64()
				fmt.Fprintf(w, svgtriangle, e[j], o[j], o[j+1], opacity)
			}
			for j := 0; j < npe-1; j++ {
				// Pointing down
				opacity := maxopacity * random.Float64()
				fmt.Fprintf(w, svgtriangle, e[j], e[j+1], o[j+1], opacity)
			}
		} else {
			// Even row
			reset(xe, y, e)
			for j := 0; j < npe; j++ {
				// Pointing down
				opacity := maxopacity * random.Float64()
				fmt.Fprintf(w, svgtriangle, o[j], o[j+1], e[j], opacity)
			}
			for j := 0; j < npe-1; j++ {
				// Pointing up
				opacity := maxopacity * random.Float64()
				fmt.Fprintf(w, svgtriangle, o[j+1], e[j], e[j+1], opacity)
			}
		}
	}
	fmt.Fprint(w, svgsuffix)
}
