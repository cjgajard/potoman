package potoman

import (
	"bufio"
	"regexp"
)

var (
	TokenEscape    = []byte(`\`)
	TokenInterp    = []byte("$")
	TokenOpen      = []byte("{")
	TokenClose     = []byte("}")
	TokenLineFeed  = []byte("\n")
	TokenSep       = []byte(" ")
	TokenHeaderSep = []byte(":")
	TokenComment   = []byte(`#`)
)

var (
	isSeparator     = regexp.MustCompile(`\W`)
	isBlank         = regexp.MustCompile(`\s+`)
	isTrailingBlank = regexp.MustCompile(`\s+$`)
	// isStartingBlank = regexp.MustCompile(`^\s+`)
)

func ScanTokens(data []byte, atEOF bool) (int, []byte, error) {
	loc := isSeparator.FindIndex(data)
	if len(loc) >= 2 {
		start, end := loc[0], loc[1]
		if start == 0 {
			return end, data[:end], nil
		}
		return start, data[:start], nil
	}
	if atEOF {
		return len(data), nil, nil
	}
	return len(data), data, nil
}

func (p *Parser) Lex() {
	ss := bufio.NewScanner(p.Source)
	ss.Split(ScanTokens)
	for ss.Scan() {
		p.TokenCh <- ss.Bytes()
	}
	close(p.TokenCh)
}
