package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type configLoader interface {
	Load(uri string) ([]io.ReadCloser, []string, error)
}

type configLoaderFunc func(string) ([]io.ReadCloser, []string, error)

func (c configLoaderFunc) Load(uri string) ([]io.ReadCloser, []string, error) {
	return c(uri)
}

type pollable interface {
	Poll()
}

type globFileLoader struct{}

func (globFileLoader) Load(path string) (data []io.ReadCloser, paths []string, err error) {
	files, err := filepath.Glob(path)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get files from file path (glob) %s, %v", path, err)
	}

	if len(files) == 0 {
		return nil, nil, fmt.Errorf("no files found in path %s", path)
	}

	var configs []io.ReadCloser
	for _, file := range files {

		f, err := os.Open(file)

		if err != nil {
			return nil, nil, fmt.Errorf("could not open %s %v", file, err)
		}

		configs = append(configs, f)
	}

	return configs, files, nil
}

func (globFileLoader) Poll() {}

type urlLoader struct{}

func (urlLoader) Load(uri string) ([]io.ReadCloser, []string, error) {
	res, err := http.Get(uri)

	if err != nil {
		return nil, nil, fmt.Errorf("couldn't load config from %s, %v", uri, err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("did not get 200 from %s, got %d", uri, res.StatusCode)
	}

	return []io.ReadCloser{res.Body}, []string{uri}, nil
}
