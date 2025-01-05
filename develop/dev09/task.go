package main

import (
	"dev09/internal/download"
	"flag"
	"fmt"
	"log"
)

func main() {
	mirror := flag.Bool("m", false, "Mirror the entire website")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: go-wget [-m] <url>")
		return
	}

	startURL := flag.Arg(0)

	if startURL == "" {
		log.Fatal("URL is required")
	}

	fmt.Println(*mirror)
	downloader := download.NewDownloader() // Создаем экземпляр загрузчика

	if *mirror {
		err := downloader.Mirror(startURL)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := downloader.Download(startURL, "")
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Download complete.")
}
