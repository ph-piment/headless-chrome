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

	srcScshoByte, err := browser.GetFullScreenshotByteByURL(ctx, sourceURL)
	if err != nil {
		log.Fatal(err)
	}
	tgtScshoByte, err := browser.GetFullScreenshotByteByURL(ctx, targetURL)
	if err != nil {
		log.Fatal(err)
	}

	srcImagePath := compareOutputDir + "source/image.png"
	tgtImagePath := compareOutputDir + "target/image.png"
	if err := image.WriteImageByByte(srcScshoByte, srcImagePath); err != nil {
		log.Fatal(err)
	}
	if err := image.WriteImageByByte(tgtScshoByte, tgtImagePath); err != nil {
		log.Fatal(err)
	}
	srcImage, err := image.ReadImageByPath(srcImagePath)
	if err != nil {
		log.Fatal(err)
	}
	tgtImage, err := image.ReadImageByPath(tgtImagePath)
	if err != nil {
		log.Fatal(err)
	}

	resImagePath := compareOutputDir + "result/image.png"
	if err := image.CompareImage(srcImage, tgtImage, resImagePath); err != nil {
		log.Fatal(err)
	}
}
