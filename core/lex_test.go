package potoman

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testLex(t *testing.T, src string, dst []string) {
	p := &Parser{
		OnEvaluate: EvaluateEnv,
		OnUnknown:  RaiseUnknown,
		Source:     strings.NewReader(src),
		TokenCh:    make(chan []byte),
	}

	go p.Lex()

	var result []string
	for v := range p.TokenCh {
		result = append(result, string(v))
	}

	if !cmp.Equal(dst, result) {
		t.Error(cmp.Diff(dst, result))
	}
}

func TestLexRequest(t *testing.T) {
	testLex(t, testRequest, testRequestTokens)
}

func TestLexRequestNoVersion(t *testing.T) {
	testLex(t, testRequestNoVersion, testRequestNoVersionTokens)
}

func TestLexRequestComments(t *testing.T) {
	testLex(t, testRequestComments, testRequestCommentsTokens)
}

func TestLexHeaders(t *testing.T) {
	testLex(t, testHeaders, testHeadersTokens)
}

func TestLexBody(t *testing.T) {
	testLex(t, testBody, testBodyTokens)
}

var testRequest = `POST api.example.com/posts?format=${FORMAT:-json} HTTP/1.1`
var testRequestTokens = []string{
	"POST", " ", "api", ".", "example", ".", "com", "/", "posts", "?",
	"format", "=", "$", "{", "FORMAT", ":", "-", "json", "}", " ", "HTTP",
	"/", "1", ".", "1",
}

var testRequestNoVersion = `POST api.example.com/posts?format=${FORMAT:-json}`
var testRequestNoVersionTokens = []string{
	"POST", " ", "api", ".", "example", ".", "com", "/", "posts", "?",
	"format", "=", "$", "{", "FORMAT", ":", "-", "json", "}",
}

var testRequestComments = `#/usr/bin/env potoman
# Ignore comments please
POST api.example.com/posts?format=${FORMAT:-json} HTTP/1.1`
var testRequestCommentsTokens = []string{
	"#", "/", "usr", "/", "bin", "/", "env", " ", "potoman", "\n",
	"#", " ", "Ignore", " ", "comments", " ", "please", "\n",
	"POST", " ", "api", ".", "example", ".", "com", "/", "posts", "?",
	"format", "=", "$", "{", "FORMAT", ":", "-", "json", "}", " ", "HTTP",
	"/", "1", ".", "1",
}

var testHeaders = `Authorization: Bearer $JWT # unsafe
  Content-Type  : application/json  `
var testHeadersTokens = []string{
	"Authorization", ":", " ", "Bearer", " ", "$", "JWT", " ", "#", " ",
	"unsafe", "\n",

	" ", " ", "Content", "-", "Type", " ", " ", ":", " ", "application",
	"/", "json", " ", " ",
}

var testBody = `
{
  # comment inside body
  "title": "${NAME}Script - A very modern programming language", # line-end comment
	# tab-indentation comment
	"body": "Everything known until now was \$tupid, with \#${NAME}Script you will be  F  A  B  U  L  O  U  S"
}`
var testBodyTokens = []string{
	"\n",

	"{", "\n",

	" ", " ", "#", " ", "comment", " ", "inside", " ", "body", "\n",
	" ", " ", `"`, "title", `"`, ":", " ", `"`, "$", "{", "NAME",
	"}", "Script", " ", "-", " ", "A", " ", "very", " ", "modern",
	" ", "programming", " ", "language", `"`, ",", " ", "#", " ",
	"line", "-", "end", " ", "comment", "\n",

	"\t", "#", " ", "tab", "-", "indentation", " ", "comment", "\n",

	"\t", `"`, "body", `"`, ":", " ", `"`, "Everything", " ", "known", " ",
	"until", " ", "now", " ", "was", " ", `\`, "$", "tupid", ",", " ",
	"with", " ", "\\", "#", "$", "{", "NAME", "}", "Script", " ", "you",
	" ", "will", " ", "be", " ", " ", "F", " ", " ", "A", " ", " ", "B",
	" ", " ", "U", " ", " ", "L", " ", " ", "O", " ", " ", "U", " ", " ",
	"S", `"`, "\n",

	"}",
}
