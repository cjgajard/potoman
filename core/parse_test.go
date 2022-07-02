package potoman

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testParse(t *testing.T, tokens []string, exp *Description) {
	os.Setenv("NAME", "Trendy")

	p := &Parser{
		OnUnknown: func(key []byte) ([]byte, error) {
			buf := new(bytes.Buffer)
			_, err := fmt.Fprintf(buf, "[%s]", key)
			return buf.Bytes(), err
		},
		OnEvaluate: EvaluateEnv,
		Source:     nil,
		TokenCh:    make(chan []byte),
	}

	go func() {
		tokens = append(tokens, "\n")
		tokens = append(tokens, testHeadersTokens...)
		tokens = append(tokens, "\n")
		tokens = append(tokens, testBodyTokens...)
		for _, buf := range tokens {
			p.TokenCh <- []byte(buf)
		}
		close(p.TokenCh)
	}()

	result, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(exp, result) {
		t.Error(cmp.Diff(exp, result))
	}
}

func TestParse(t *testing.T) {
	testParse(t, testRequestCommentsTokens, &Description{
		Method:  "POST",
		URL:     "api.example.com/posts?format=[${FORMAT:-json}]",
		Version: Attr{"1", "1"},
		Headers: []Attr{
			{Key: "Authorization", Value: "Bearer [$JWT]"},
			{Key: "Content-Type", Value: "application/json"},
		},
		Body: []byte(
			`{
  "title": "TrendyScript - A very modern programming language",
	"body": "Everything known until now was $tupid, with #TrendyScript you will be  F  A  B  U  L  O  U  S"
}`,
		),
	})
}
