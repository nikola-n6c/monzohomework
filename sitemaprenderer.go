package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
)

// Basic interface for rendering
type SiteMapRenderer interface {
	Render(smap *SiteMap) (*bytes.Buffer, error)
}

type JSONSiteMapRenderer struct{}

func (r JSONSiteMapRenderer) Render(smap *SiteMap) (*bytes.Buffer, error) {
	b, err := json.MarshalIndent(smap, "", "  ")
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

// Not very pretty but nice demonstration on lower crawl depths
type SVGSiteMapRenderer struct{}

// Graphics code incoming
func (r SVGSiteMapRenderer) Render(smap *SiteMap) (*bytes.Buffer, error) {
	width := 2000
	height := 2000

	buf := new(bytes.Buffer)
	canvas := svg.New(buf)
	canvas.Start(width, height)

	canvas.Def()
	// Thanks to
	// https://vanseodesign.com/web-design/svg-markers/
	canvas.Marker("arrow", 10, 6, 50, 50, `orient="auto"`)
	canvas.Path("M0,0 L0,12 L9,6 z", "fill:red")
	canvas.MarkerEnd()
	canvas.DefEnd()

	whereIsX := make(map[string]float64)
	whereIsY := make(map[string]float64)

	paddedWidth := width - 300
	paddedHeight := height - 300

	pages := len(smap.smap)
	delta := 2 * math.Pi / float64(pages)

	q := make([]string, 0)
	i := 0

	q = append(q, smap.root)

	// Render text and circles using basic BFS
	for len(q) > 0 {
		node := q[0]
		q = q[1:]

		angle := float64(i) * delta

		x := math.Cos(angle)*float64(paddedWidth/2) + float64(width/2)
		y := math.Sin(angle)*float64(paddedHeight/2) + float64(height/2)

		whereIsX[node] = x
		whereIsY[node] = y

		anchor := "start"
		canvas.Circle(int(x), int(y), 5, "fill:black;stroke:black")
		if angle > math.Pi/2 && angle < 3*math.Pi/2 {
			// Preventing neck injury
			angle += math.Pi
			anchor = "end"
		}
		canvas.Group(fmt.Sprintf(`transform="rotate(%f, %f %f)"`, angle*180.0/math.Pi, x, y))
		canvas.Text(int(x), int(y), node, fmt.Sprintf(`text-anchor:%s;font-size:12px;fill:black;`, anchor))
		canvas.Gend()

		q = append(q, smap.Get(node).AsSlice()...)
		i++
	}

	// Render lines using BFS
	q = append(q, smap.root)
	for len(q) > 0 {
		node := q[0]
		q = q[1:]

		for _, c := range smap.Get(node).AsSlice() {
			canvas.Line(int(whereIsX[node]), int(whereIsY[node]), int(whereIsX[c]), int(whereIsY[c]), "stroke:black", `marker-end="url(#arrow)"`)
		}

		q = append(q, smap.Get(node).AsSlice()...)
	}

	canvas.End()

	return buf, nil
}

// No way I can code this now but would be cool to try
/*
type ASCIISiteMapRenderer struct{}

func (r ASCIISiteMapRenderer) Render(smap *SiteMap) ([]byte, error) {
	// No way I can do quickly enough
	return nil, errors.New("Not implemented")
}
*/
