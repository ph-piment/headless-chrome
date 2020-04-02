# About headless-chrome

This is my Go lang study repository.

Currently, the following functions are provided.

* ./cmd/compare/main.go

  Pass the two URLs to compare and compare the images on the web page.
  (Want to use in regression test)

* ./cmd/parallel_compare/main.go

  Compare the web page image from "config/url_list.toml".
  (Want to use in regression test)

* ./cmd/register/main.go

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
cd deployments
docker-compose up -d --build
docker exec -it workspace /bin/bash
go run ./cmd/compare/main.go https://www.google.com https://www.google.com
go run ./cmd/parallel_compare/main.go
```

### Outputs
* Web image of the first argument.(/go/src/work/outputs/images/compare/source/image.png)
* Web image of the second argument.(/go/src/work/outputs/images/compare/target/image.png)
* Image comparison results.(/go/src/work/outputs/images/compare/result/image.png)

### Wish list
* Add test code

### Resources
* https://github.com/golang-standards/project-layout
* https://github.com/chromedp/chromedp
* https://github.com/chromedp/docker-headless-shell
* https://github.com/orisano/pixelmatch
