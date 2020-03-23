// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"flag"
	"log"

	"work/internal/browser"
)

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

	sourceImage := browser.GetImageByURL(ctx, sourceURL, "/source/image.png")
	targetImage := browser.GetImageByURL(ctx, targetURL, "/target/image.png")

	browser.DiffImage(sourceImage, targetImage, "/result/image.png")
}
