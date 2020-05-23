package fetcher

import (
	"errors"
	"io"
	"net/http"
	"os"
)

type Fetch func(url string) (io.Reader, error)

func FileFetch(url string) (io.Reader, error) {
	file, err := os.Open(url)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func HttpFetch(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, errors.New(url + "bad status code" + string(resp.StatusCode))
	}
	return resp.Body, nil
}

func GeneralFetch(url string) (io.Reader, error) {
	httpText := "http://"
	httpsText := "https://"
	if url[0:len(httpsText)] == httpsText || url[0:len(httpText)] == httpText {
		return HttpFetch(url)
	}
	return FileFetch(url)
}
