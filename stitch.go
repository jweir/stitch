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

	base := createBaseImage()

	counter := 0.0
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == ".png" || ext == ".jpg" {
			counter = counter + 1
			if math.Mod(counter, float64(nth)) == 0 {
				log.Printf("Sitching %s", f.Name())
				img := readImage(filepath.Join(src, f.Name()))
				base = stitch(base, img)
			}
		}
	}

	o, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	jpeg.Encode(o, base, &jpeg.Options{Quality: 100})

	log.Printf("Wrote to %s", out)
}

func stitch(base, new image.Image) image.Image {
	bh := base.Bounds().Dy()
	bw := base.Bounds().Dx()

	nh := new.Bounds().Dy()
	nw := new.Bounds().Dx()

	w := max(nw, bw)

	nr := image.Rect(0, 0, w, bh+nh)
	r := image.NewRGBA(nr)

	draw.Draw(r, base.Bounds(), base, image.Point{0, 0}, draw.Src)

	rc2 := image.Rect(0, bh, nw, bh+nh)
	draw.Draw(r, rc2, new, image.Point{0, 0}, draw.Src)

	return r
}

func createBaseImage() image.Image {
	r := image.Rect(0, 0, 1, 1)
	img := image.NewRGBA(r)
	return img
}

func readImage(path string) image.Image {
	file, err := os.Open(path)

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
