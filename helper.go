package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	olo "github.com/xorpaul/sigolo"
)

func handleError(response *http.Response, err error, w http.ResponseWriter) {
	if response != nil {

		if response.StatusCode == 401 {
			olo.Info(err.Error())
		} else {
			olo.Error(err.Error())
		}

		for attr, val := range response.Header {
			if strings.ToLower(attr) == "www-authenticate" {
				w.Header().Set(attr, val[0])
			}
		}
		w.WriteHeader(response.StatusCode)
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			olo.Error("Error while reading failed response body")
		}
		w.Write(bodyBytes)
	} else {
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
	}
}

func removeSchemeFromURL(requestedURL string) (string, error) {
	url, err := url.Parse(requestedURL)
	if err != nil {
		return "", fmt.Errorf("unable to remove URL scheme from requested URL '%s'", requestedURL)
	}
	return strings.Replace(requestedURL, url.Scheme+"://", "", 1), nil
}

func ensureDir(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return merr
		}
	}
	return nil
}

func validateCacheURL(cacheURL string) error {
	if strings.Contains(cacheURL, "..") {
		return errors.New(".. is not allowed in request URL")
	}
	if strings.HasSuffix(cacheURL, "/") {
		return errors.New("request URL ending with / is not allowed")
	}
	return nil
}
