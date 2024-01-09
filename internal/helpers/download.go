package helpers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/xid"
)

type DownloadedFilesPath struct {
	Path  string
	Error error
}

func FilenameFromUrl(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	x, _ := url.QueryUnescape(u.EscapedPath())
	return filepath.Base(x), nil
}

func DownloadFiles(urls []string, dir string) []DownloadedFilesPath {
	wg := sync.WaitGroup{}
	ch := make(chan *DownloadedFilesPath)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			if url == "" {
				ch <- &DownloadedFilesPath{
					Error: fmt.Errorf("url empty"),
				}
				return
			}

			filename, err := FilenameFromUrl(url)

			if err != nil {
				ch <- &DownloadedFilesPath{
					Error: fmt.Errorf("url parsing failed"),
				}
				return
			}

			tempDir := filepath.Join(".", dir)
			err = os.MkdirAll(tempDir, os.ModePerm)

			if err != nil {
				panic("uploadthing: failed to create or access temp dir")
			}

			fullPath := tempDir + xid.New().String() + filename
			file, err := os.Create(fullPath)

			if err != nil {
				panic("uploadthing: failed to create file in temp dir")
			}

			defer file.Close()

			resp, err := http.Get(url)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				ch <- &DownloadedFilesPath{
					Error: fmt.Errorf("failed downloading file: server responded with status -" + resp.Status),
				}
				return
			}

			_, err = io.Copy(file, resp.Body)

			if err != nil {
				ch <- &DownloadedFilesPath{
					Error: fmt.Errorf("failed copy file"),
				}
				return
			}

			ch <- &DownloadedFilesPath{
				Path: fullPath,
			}
		}(url)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	s := []DownloadedFilesPath{}

	for df := range ch {
		s = append(s, *df)
	}

	return s
}
