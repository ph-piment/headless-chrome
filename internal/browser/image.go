package browser

import (
	"context"
	"errors"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

type colorValue color.RGBA

var compareDir = filepath.Dir("/go/src/work/outputs/images/compare/")

func openImage(path string) (image.Image, error) {
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

func GetImageByURL(ctx context.Context, url string, imagePath string) image.Image {
	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	sourceImagePath := compareDir + imagePath
	if err := ioutil.WriteFile(sourceImagePath, buf, 0644); err != nil {
		log.Fatal(err)
	}
	imgfile, err := openImage(sourceImagePath)
	if err != nil {
		log.Fatal(err)
	}

	return imgfile
}

func DiffImage(sourceImage image.Image, targetImage image.Image, imagePath string) {
	// compare
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
	f, err := os.Create(compareDir + "/result/image.png")
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
