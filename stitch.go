package main

import (
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
)

func main() {
	var out string
	var src string
	var nth int

	flag.StringVar(&out, "out", "", "filename to output the images")
	flag.StringVar(&src, "src", "", "directory containing the images to stitch")
	flag.IntVar(&nth, "nth", 1, "only stitch every nth image")

	flag.Parse()

	if out == "" || src == "" {
		flag.PrintDefaults()
		return
	}

	files, err := ioutil.ReadDir(src)

	if err != nil {
		log.Fatal(err)
	}

	selected := filter(files, nth)
	domain := image.Rect(0, 0, 0, 0)

	for _, f := range selected {
		img := readImage(filepath.Join(src, f.Name()))
		domain = precalcSize(domain, img)
	}

	base := createBaseImage(domain)

	bh := 0

	log.Printf("Creating an image %d x %d pixels", base.Bounds().Dx(), base.Bounds().Dy())

	for _, f := range selected {
		log.Printf("Sitching %s", f.Name())
		img := readImage(filepath.Join(src, f.Name()))
		bh, base = stitch(bh, base, img)
	}

	o, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	jpeg.Encode(o, base, &jpeg.Options{Quality: 100})

	log.Printf("Wrote to %s", out)
}

func filter(files []os.FileInfo, nth int) []os.FileInfo {
	out := []os.FileInfo{}

	counter := 0.0
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == ".png" || ext == ".jpg" {
			counter = counter + 1
			if math.Mod(counter, float64(nth)) == 0 {
				out = append(out, f)
			}
		}
	}

	return out
}

func precalcSize(r image.Rectangle, new image.Image) image.Rectangle {
	bh := r.Dy()
	bw := r.Dx()

	nh := new.Bounds().Dy()
	nw := new.Bounds().Dx()

	w := max(nw, bw)

	return image.Rect(0, 0, w, bh+nh)
}

func stitch(bh int, base *image.RGBA, new image.Image) (int, *image.RGBA) {
	nh := new.Bounds().Dy()
	nw := new.Bounds().Dx()

	rc2 := image.Rect(0, bh, nw, bh+nh)
	draw.Draw(base, rc2, new, image.Point{0, 0}, draw.Src)

	return bh + nh, base
}

func createBaseImage(domain image.Rectangle) *image.RGBA {
	img := image.NewRGBA(domain)
	return img
}

func readImage(path string) image.Image {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(file)

	if err != nil {
		log.Fatal(err)
	}

	return img
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
