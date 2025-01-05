package download

import (
	"bytes"
	"dev09/internal/fetch"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Downloader struct {
	baseURL *url.URL
}

func NewDownloader(baseURL string) (*Downloader, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Downloader{parsedURL}, nil
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

func (d *Downloader) mirror(urlStr, prefix string, v visited) error {
	fmt.Println("Mirroring:", urlStr)
	if v.lvl < 0 {
		return nil
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = d.baseURL.Scheme
	}

	if parsedURL.Host == "" {
		parsedURL.Host = d.baseURL.Host
	}

	if parsedURL.Host != d.baseURL.Host {
		return nil
	}

	resp, err := fetch.FetchURL(parsedURL.String())
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	localPath := path.Join("mirror", parsedURL.Host, parsedURL.Path)
	if strings.HasSuffix(parsedURL.Path, "/") || parsedURL.Path == "" {
		localPath = path.Join(localPath, "index.html")
	}

	err = os.MkdirAll(path.Dir(localPath), 0755)
	if err != nil {
		return err
	}

	// Создаем локальный файл
	file, err := os.Create(localPath)
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

	links, err := fetch.ExtractLinks(urlStr)
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
			err := d.mirror(link, filepath.Join(prefix, filepath.Dir(localPath)), v)
			fmt.Println("\033[1m\033[31m", err, "\033[0m")
			return err
		})
	}

	return eg.Wait()
}
