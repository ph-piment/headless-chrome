package main

import (
	"log"
	"strconv"
	"sync"

	"work/config"
	"work/internal/browser"
	"work/internal/image"
)

// ThreadCntUpperLimit is upper limit of thread count.
const ThreadCntUpperLimit = 3

func main() {
	urlList, err := config.NewConfig("url_list")
	if err != nil {
		log.Fatal(err.Error())
	}

	limit := make(chan struct{}, ThreadCntUpperLimit)

	var wg sync.WaitGroup
	for i, url := range urlList.URLLIST {
		wg.Add(1)
		go func(index int, sourceURL string, targetURL string) {
			limit <- struct{}{}
			defer wg.Done()
			outputDiff(sourceURL, targetURL, index)
			<-limit
		}(i, url.SourceURL, url.TargetURL)
	}
	wg.Wait()
}

func outputDiff(sourceURL string, targetURL string, index int) {
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
