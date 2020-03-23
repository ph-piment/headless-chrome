// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"flag"
	"log"

	"work/internal/browser"
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

	ctx, allocCancel, ctxtCancel := browser.GetContext()
	defer allocCancel()
	defer ctxtCancel()

	sourceImage := browser.GetImageByURL(ctx, sourceURL, compareOutputDir+"source/image.png")
	targetImage := browser.GetImageByURL(ctx, targetURL, compareOutputDir+"target/image.png")

	browser.DiffImage(sourceImage, targetImage, compareOutputDir+"result/image.png")
}
