package grpchandler

import "net/url"

// CreateLink accepts chars and baseURL for building url.URL.
func CreateLink(chars, baseURL string) (*url.URL, error) {
	URL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	URL.Path = chars

	return URL, nil
}
