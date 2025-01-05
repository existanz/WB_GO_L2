package download

import (
	"bytes"
	"dev09/internal/fetch"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(url, path string) error {
	// TODO: create file with correct name
	resp, err := fetch.FetchURL(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	filename := filepath.Join(path, filepath.Base(url))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (d *Downloader) Mirror(url string) error {
	links := make(map[string]struct{})
	fmt.Println("Mirroring:", url)
	return d.mirror(url, "", visited{links: links, mu: &sync.Mutex{}, lvl: 2})
}

type visited struct {
	links map[string]struct{}
	mu    *sync.Mutex
	lvl   int
}

func (v visited) isVisited(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	_, ok := v.links[url]
	if ok {
		return true
	}
	v.links[url] = struct{}{}
	return false
}

func (d *Downloader) mirror(url, prefix string, v visited) error {
	fmt.Println("Mirroring:", url)
	if v.lvl < 0 {
		return nil
	}

	resp, err := fetch.FetchURL(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	filename := filepath.Join(prefix, filepath.Base(url))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	io.TeeReader(resp.Body, buf)

	_, err = io.Copy(file, buf)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	links, err := fetch.ExtractLinks(url)
	if err != nil {
		return fmt.Errorf("failed to extract links: %w", err)
	}
	v.lvl--

	var eg errgroup.Group
	for link := range links {
		fmt.Println(link)
		if !fetch.IsValidURL(link) && v.isVisited(link) {
			continue
		}

		eg.Go(func() error {
			err := d.mirror(link, filepath.Join(prefix, filepath.Dir(url)), v)
			fmt.Println("\033[1m\033[31m", err, "\033[0m")
			return err
		})
	}

	return eg.Wait()
}
