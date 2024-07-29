package runner

import (
	"bufio"
	"os"
)

type InputProcessor struct {
	urls chan string
}

func NewInputProcessor(url, urlFile string) (*InputProcessor, error) {
	ip := &InputProcessor{
		urls: make(chan string),
	}

	go func() {
		defer close(ip.urls)

		if url != "" {
			ip.urls <- url
		}

		if urlFile != "" {
			file, err := os.Open(urlFile)
			if err != nil {
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				ip.urls <- scanner.Text()
			}
		}
	}()

	return ip, nil
}

func (ip *InputProcessor) URLs() <-chan string {
	return ip.urls
}
