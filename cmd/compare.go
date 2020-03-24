// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"flag"
	"log"

	"work/internal/browser"
	"work/internal/image"
)

const compareOutputDir = "/go/src/work/outputs/images/compare/"

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage of pixelmatch [flags] image1 image2 :")
	}
	sourceURL := args[0]
	targetURL := args[1]

	ctx, allocCxl, ctxCxl := browser.GetContext()
	defer allocCxl()
	defer ctxCxl()

	sourceScshoByte, err := browser.GetFullScreenshotByteByURL(ctx, sourceURL)
	if err != nil {
		log.Fatal(err)
	}
	targetScshoByte, err := browser.GetFullScreenshotByteByURL(ctx, targetURL)
	if err != nil {
		log.Fatal(err)
	}

	sourceImagePath := compareOutputDir + "source/image.png"
	targetImagePath := compareOutputDir + "target/image.png"
	if err := image.WriteImageByByte(sourceScshoByte, sourceImagePath); err != nil {
		log.Fatal(err)
	}
	if err := image.WriteImageByByte(targetScshoByte, targetImagePath); err != nil {
		log.Fatal(err)
	}
	sourceImage, err := image.ReadImageByPath(sourceImagePath)
	if err != nil {
		log.Fatal(err)
	}
	targetImage, err := image.ReadImageByPath(targetImagePath)
	if err != nil {
		log.Fatal(err)
	}

	resultImagePath := compareOutputDir + "result/image.png"
	if err := image.CompareImage(sourceImage, targetImage, resultImagePath); err != nil {
		log.Fatal(err)
	}
}
