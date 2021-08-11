package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var addr = flag.String("addr", "localhost:8080", "http service adress")
var fnt = goregular.TTF

func main() {
	flag.Parse()
	log.Printf("Listening on %s", *addr)
	http.Handle("/", http.HandlerFunc(placeHolder))
	log.Fatal("ListenAndServe:", http.ListenAndServe(*addr, nil))
}

func placeHolder(w http.ResponseWriter, r *http.Request) {
	width, height := plSz(r.URL.Path)
	if width == 0 {
		width = 100
	}
	if height == 0 {
		height = 100
	}
	enc := png.Encoder{}
	img := plBg(width, height)
	if err := plText(img, fmt.Sprintf("%vx%v", width, height)); err != nil {
		log.Printf("plText: %v", err)
	}
	if err := enc.Encode(w, img); err != nil {
		log.Printf("Encode: %v", err)
	}
}

func plSz(urlStr string) (w int, h int) {
	wd, tl := ShiftPath(urlStr)
	w = atoi(wd)
	hd, tl := ShiftPath(tl)
	h = atoi(hd)
	return
}

func atoi(a string) int {
	w, err := strconv.Atoi(a)
	if err != nil {
		log.Println("Atoi:", err)
	}
	return w
}
func plBg(w, h int) *image.RGBA {
	rect := image.Rect(0, 0, w, h)
	fill := color.RGBA{255, 255, 125, 255}
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{fill}, image.Point{}, draw.Src)
	return img
}

func plText(img image.Image, txt string) error {
	font, err := truetype.Parse(fnt)
	if err != nil {
		return fmt.Errorf("truetype.Parse(%v): %v ", fnt, err)
	}
	log.Println(font)
	return nil
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
