package internal

import (
	"errors"
	"net/url"
)

func GetUrlHost(rawURL string) (string, error) {
	newURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", err
	}

	if newURL.OmitHost {
		return "", errors.New("missing or invalid host")
	}

	return newURL.Scheme + "://" + newURL.Host, nil
}
