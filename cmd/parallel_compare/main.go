package main

import (
	"fmt"
	"log"
	"strconv"

	"work/config"
	"work/internal/browser"
	"work/internal/image"
)

func main() {
	sourceURL := ""
	targetURL := ""
	index := 0
	urlList, err := config.NewConfig("url_list")
	if err != nil {
		log.Fatal(err.Error())
	}
	for i, url := range urlList.URLLIST {
		index = i
		sourceURL = url.SourceURL
		targetURL = url.TargetURL
		fmt.Printf("index: %d, source_name: %s, target_name: %s\n", index, sourceURL, targetURL)
	}

	conf, err := config.NewConfig("app")
	if err != nil {
		log.Fatal(err.Error())
	}

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

	fileName := strconv.Itoa(index) + ".png"
	srcImagePath := conf.PATH.OutputCompareDir + "source/" + fileName
	tgtImagePath := conf.PATH.OutputCompareDir + "target/" + fileName
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

	resImagePath := conf.PATH.OutputCompareDir + "result/" + fileName
	if err := image.CompareImage(srcImage, tgtImage, resImagePath); err != nil {
		log.Fatal(err)
	}
}
