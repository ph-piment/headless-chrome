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

	sourceScreenshotByte, _ := browser.GetFullScreenshotByteByURL(ctx, sourceURL)
	targetScreenshotByte, _ := browser.GetFullScreenshotByteByURL(ctx, targetURL)

	sourceImagePath := compareOutputDir + "source/image.png"
	targetImagePath := compareOutputDir + "target/image.png"
	resultImagePath := compareOutputDir + "result/image.png"
	image.WriteImageByByte(sourceScreenshotByte, sourceImagePath)
	image.WriteImageByByte(targetScreenshotByte, targetImagePath)
	sourceImage, _ := image.ReadImageByPath(sourceImagePath)
	targetImage, _ := image.ReadImageByPath(targetImagePath)

	image.CompareImage(sourceImage, targetImage, resultImagePath)
}
