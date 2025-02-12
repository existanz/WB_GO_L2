package fetch

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// FetchURL запрашивает URL и возвращает ответ.
func FetchURL(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d for URL %s", resp.StatusCode, url)

	}
	return resp, nil
}

// ExtractLinks извлекает ссылки из HTML-кода.
func ExtractLinks(page string) (<-chan string, error) {
	resp, err := http.Get(page)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	linksChan := make(chan string)

	go func() {
		defer close(linksChan)

		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "link") {
				for _, a := range n.Attr {
					if a.Key == "href" {
						linksChan <- a.Val
						break
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
	}()

	return linksChan, nil
}

func IsValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func SanitizeURL(url string) string {
	return strings.ReplaceAll(url, "/", "_")
}
