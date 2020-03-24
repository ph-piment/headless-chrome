package image

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/orisano/pixelmatch"

	"github.com/pkg/errors"
)

type colorValue color.RGBA

// OpenImage open by file path
func OpenImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open")
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}
	return img, nil
}

// WriteImageByByte write image by bytes
func WriteImageByByte(buf []byte, imagePath string) error {
	if error := ioutil.WriteFile(imagePath, buf, 0644); error != nil {
		return error
	}
	return nil
}

// DiffImage compare sourceImage and targetImage.
func DiffImage(sourceImage image.Image, targetImage image.Image, imagePath string) {
	threshold := flag.Float64("threshold", 0.1, "threshold")
	aa := flag.Bool("aa", false, "ignore anti alias pixel")
	alpha := flag.Float64("alpha", 0.1, "alpha")
	antiAliased := colorValue(color.RGBA{R: 255, G: 255})
	diffColor := colorValue(color.RGBA{R: 255})
	var out image.Image
	opts := []pixelmatch.MatchOption{
		pixelmatch.Threshold(*threshold),
		pixelmatch.Alpha(*alpha),
		pixelmatch.AntiAliasedColor(color.RGBA(antiAliased)),
		pixelmatch.DiffColor(color.RGBA(diffColor)),
		pixelmatch.WriteTo(&out),
	}
	if *aa {
		opts = append(opts, pixelmatch.IncludeAntiAlias)
	}

	_, err := pixelmatch.MatchPixel(sourceImage, targetImage, opts...)
	if err != nil {
		log.Fatal(err)
	}

	var w io.Writer
	f, err := os.Create(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w = f

	var encErr error
	encErr = png.Encode(w, out)
	if encErr != nil {
		log.Fatal(err)
	}
}
