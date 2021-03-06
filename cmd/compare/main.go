package main

import (
	"flag"
	"log"

	"work/config"
	"work/internal/browser"
	"work/internal/image"
)

func main() {
	conf, err := config.NewConfig("app")
	if err != nil {
		log.Fatal(err.Error())
	}

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

	srcImagePath := conf.PATH.OutputCompareDir + "source/image.png"
	tgtImagePath := conf.PATH.OutputCompareDir + "target/image.png"
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

	resImagePath := conf.PATH.OutputCompareDir + "result/image.png"
	if err := image.CompareImage(srcImage, tgtImage, resImagePath); err != nil {
		log.Fatal(err)
	}
}
