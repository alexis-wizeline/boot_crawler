package scrapper

import (
	"errors"
	"net/url"
	"path"
)

func NormalizeURL(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	if parsed.Hostname() == "" && parsed.EscapedFragment() == "" {
		return "", errors.New("Invalid url")
	}

	parsed.Path = path.Clean(parsed.Path)
	return parsed.Hostname() + parsed.EscapedPath(), nil
}
