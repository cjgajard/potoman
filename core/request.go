package potoman

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
)

var (
	hasProtocol = regexp.MustCompile(`^[a-z]{2,}://`)
)

func NewRequest(d *Description) (*http.Request, error) {
	url := d.URL
	if !hasProtocol.MatchString(url) {
		// TODO: improve default protocol setting. current implementation
		// may fail with weird input like `a://www.google.com`
		url = "https://" + url
	}

	req, err := http.NewRequest(d.Method, url, bytes.NewReader(d.Body))
	if err != nil {
		return nil, err
	}

	for _, h := range d.Headers {
		req.Header.Add(h.Key, h.Value)
	}

	return req, nil
}

type Description struct {
	Method  string
	URL     string
	Version Attr
	Headers []Attr
	Body    json.RawMessage
}

type Attr struct {
	Key   string
	Value string
}
