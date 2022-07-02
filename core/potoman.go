package potoman

import (
	"fmt"
	"io"
	"os"
)

type Parser struct {
	OnEvaluate MessageFunc
	OnUnknown  MessageFunc
	Source     io.Reader
	TokenCh    chan []byte
}

func NewParser(src io.Reader) *Parser {
	ch := make(chan []byte)
	p := &Parser{
		OnEvaluate: EvaluateEnv,
		OnUnknown:  RaiseUnknown,
		Source:     src,
		TokenCh:    ch,
	}
	return p
}

func (p *Parser) Read() (*Description, error) {
	go p.Lex()
	return p.Parse()
}

type MessageFunc func(key []byte) ([]byte, error)

func (p *Parser) SetOnUnknown(f MessageFunc) {
	p.OnUnknown = f
}

func (p *Parser) SetOnEvaluate(f MessageFunc) {
	p.OnEvaluate = f
}

func EvaluateEnv(key []byte) (out []byte, err error) {
	return []byte(os.ExpandEnv(string(key))), nil
}

func RaiseUnknown(key []byte) ([]byte, error) {
	return nil, fmt.Errorf("%s is empty", key)
}
