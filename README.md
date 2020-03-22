# About headless-chrome

This is my Go lang study repository.

Currently, the following functions are provided.

* compare.go

  Pass the two URLs to compare and compare the images on the web page.
  (Want to use in regression test)

* register.go -> not started yet...

  Register web page information in the database.
  (Want to use in scraping)

## About compare.go
Pass the two URLs to compare and compare the images on the web page.
(Want to use in regression test)

The processing order is as follows.
1. Get URL from parameters
1. Request URL to headless browser
1. Get images from WEB screen
1. Compare images

### Usage
```sh
git clone https://github.com/ph-piment/headless-chrome.git
cd docker
docker-compose up -d --build
docker exec -it workspace /bin/bash
go run ./cmd/compare.go https://www.google.com https://www.google.com
```

### Outputs
* Web image of the first argument.(/go/src/work/outputs/images/compare/source/image.png)
* Web image of the second argument.(/go/src/work/outputs/images/compare/target/image.png)
* Image comparison results.(/go/src/work/outputs/images/compare/result/image.png)

### Wish list
* Refactor
* Add test code
* Parallel processing
* Exception handling organization
* Add circle ci

### Resources
* https://github.com/golang-standards/project-layout
* https://github.com/chromedp/chromedp
* https://github.com/chromedp/docker-headless-shell
* https://github.com/orisano/pixelmatch
