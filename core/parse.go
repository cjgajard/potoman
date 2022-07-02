package potoman

import (
	"bytes"
	"fmt"
)

type Stage int

const (
	Method Stage = iota
	Endpoint
	Version
	Header
	Body
)

func (p *Parser) Parse() (*Description, error) {
	r := &Description{}
	stage := Method
	isEscaped := false
	isInterpolation := false
	isInterpOpen := false
	isEvaluation := false
	isComment := false

	var content []byte
	var headerKey []byte
	var param []byte

	for token := range p.TokenCh {
		if isComment && !bytes.Equal(token, TokenLineFeed) {
			continue
		}
		if isEscaped {
			// TODO: transform \n, \t, etc.?
			if !bytes.Equal(token, TokenInterp) && !bytes.Equal(token, TokenComment) {
				content = append(content, TokenEscape...)
			}
			content = append(content, token...)
			isEscaped = false
			continue
		}
		if isInterpolation {
			// Only allow opening an interpolation bracket right
			// after an interpolation character
			isInterpolation = false
			param = append(param, token...)
			if bytes.Equal(token, TokenOpen) {
				isInterpOpen = true
				continue
			}
			isEvaluation = true
		}
		if isInterpOpen {
			// If an interpolation bracket was open and a TokenClose
			// appears we evaluate its value, every other token get
			// added to `param` and then we pass to next token
			param = append(param, token...)
			if !bytes.Equal(token, TokenClose) {
				continue
			}
			isInterpOpen = false
			isEvaluation = true
		}
		if isEvaluation {
			value, err := p.evaluate(param)
			if err != nil {
				return nil, err
			}
			content = append(content, value...)
			param = nil
			isEvaluation = false
			isInterpOpen = false
			continue
		}
		if bytes.Equal(token, TokenEscape) {
			isEscaped = true
			continue
		}
		if bytes.Equal(token, TokenComment) {
			isComment = true
			content = removeTrailing(content)
			continue
		}
		if bytes.Equal(token, TokenInterp) {
			isInterpolation = true
			param = append(param, token...)
			continue
		}
		if bytes.Equal(token, TokenHeaderSep) {
			if stage == Header {
				if param != nil {
					return r, fmt.Errorf("incomplete header %s", param)
				}
				headerKey = removeTrailing(content)
				content = nil
				continue
			}
		}
		if bytes.Equal(token, TokenSep) {
			switch stage {
			case Method:
				r.Method = string(content)
				stage++
				content = nil
				continue
			case Endpoint:
				r.URL = string(content)
				stage++
				content = nil
				continue
			case Version:
				continue
			case Header:
				if content == nil {
					continue
				}
			}
		}
		if bytes.Equal(token, TokenLineFeed) {
			isComment = false
			switch stage {
			case Method:
				continue
			case Endpoint:
				r.URL = string(content)
				stage += 2
				content = nil
				continue
			case Version:
				v, err := getHttpVersion(content)
				if err != nil {
					return nil, err
				}
				r.Version = v
				stage++
				content = nil
				continue
			case Header:
				if param != nil {
					return r, fmt.Errorf("incomplete token %s", param)
				}
				if headerKey == nil {
					stage++
					continue
				}
				r.Headers = append(r.Headers, Attr{
					Key:   string(headerKey),
					Value: string(removeTrailing(content)),
				})
				headerKey = nil
				content = nil
				continue
			}
		}
		content = append(content, token...)
	}
	r.Body = append(r.Body, content...)
	return r, nil
}

func (p *Parser) evaluate(key []byte) ([]byte, error) {
	value, err := p.OnEvaluate(key)
	if err != nil {
		return p.OnUnknown(key)
	}
	if len(value) == 0 {
		return p.OnUnknown(key)
	}
	return value, nil
}

func removeTrailing(src []byte) []byte {
	if loc := isTrailingBlank.FindIndex(src); loc != nil {
		src = src[:loc[0]]
	}
	return src
}
