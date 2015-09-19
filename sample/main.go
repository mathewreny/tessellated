package main

import (
	"flag"
	"fmt"
	"github.com/mathewreny/tessellated"
	"net/http"
	"strconv"
	"time"
)

type Printer struct{}

func (pr Printer) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

var port int
var Width, Height float64

func init() {
	flag.IntVar(&port, "http", -1, "Specify the port number. Example '--http 8080'")
	flag.Float64Var(&Width, "width", 1000, "The width of the svg.")
	flag.Float64Var(&Height, "height", 1000, "The height of the svg.")
}

func main() {
	flag.Parse()

	if -1 == port {
		var p Printer
		tessellated.Triangle(tessellated.Rect{Width, Height}, p)
	}

	if 65535 < port || port < 0 {
		fmt.Println("Invalid Port!")
		return
	}

	http.HandleFunc("/", root)
	http.HandleFunc("/triangle.svg", triangle)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Println(err)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!doctype html>
<html>
	<head></head>
	<body>
 		<h1>Tessellation server is running</h1>
		<a href="/triangle.svg?width=2000&height=2000">click here for a 2000x2000 triangle svg</a>
	 </body>
</svg>`)
}

func triangle(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	width, err := strconv.Atoi(r.FormValue("width"))
	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if width <= 0 || height <= 0 {
		http.Error(w, "Size not valid", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	tessellated.Triangle(tessellated.Rect{float64(width), float64(height)}, w)
	fmt.Printf("Triangle background of size %d, %d took %v\n", width, height, time.Since(t))
}
