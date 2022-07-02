package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/cjgajard/persistent-cookiejar"
	"github.com/cjgajard/potoman/core"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "WARNING", err)
	}

	opt, err := NewOptions()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	parser := potoman.NewParser(opt.Source)
	parser.SetOnEvaluate(evaluateWithBinBash)
	if !failFlag {
		parser.SetOnUnknown(promptUnknown)
	}

	des, err := parser.Read()
	opt.Source.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	req, err := potoman.NewRequest(des)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	jarOk := true
	jar, err := cookiejar.New(&cookiejar.Options{SaveAll: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, "WARNING", err)
		jarOk = false
	}

	if jarOk {
		cookies := jar.Cookies(req.URL)
		for _, c := range cookies {
			cookieCopy := http.Cookie{Name: c.Name, Value: c.Value}
			req.Header.Add("Cookie", cookieCopy.String())
		}
	}

	if verboseFlag {
		fmt.Fprintln(os.Stderr, ">", req.URL)
		c := req.Header.Get("Content-Type")
		if strings.Contains(c, "application/json") {
			var pretty bytes.Buffer
			json.Indent(&pretty, des.Body, "> ", "    ")
			fmt.Fprintln(os.Stderr, ">", pretty.String())
		} else {
			fmt.Fprintln(os.Stderr, ">", string(des.Body))
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	c := res.Header.Get("Content-Type")
	if !rawFlag && strings.Contains(c, "application/json") {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, buf, "", "    "); err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, buf)
			os.Exit(1)
		}
		fmt.Println(pretty.String())
	} else {
		fmt.Println(string(buf))
	}

	if jarOk {
		cookies := res.Cookies()
		jar.SetCookies(req.URL, cookies)
		if err := jar.Save(); err != nil {
			fmt.Fprintln(os.Stderr, "WARNING", err)
		}
	}
}

func NewOptions() (Options, error) {
	opt := Options{}

	args := flag.Args()
	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	if path == "" || path == "-" {
		opt.Source = pipe{File: os.Stdin, keep: true}
	} else {
		var f *os.File
		f, err := os.Open(path)
		if err != nil {
			return opt, err
		}
		opt.Source = pipe{File: f, keep: false}
	}
	return opt, nil
}

var failFlag bool
var verboseFlag bool
var rawFlag bool

func init() {
	flag.BoolVar(
		&failFlag, "fail", false,
		"if not set, potoman asks for missing values through stdin",
	)
	flag.BoolVar(
		&verboseFlag, "v", false,
		"if set, also prints the request body",
	)
	flag.BoolVar(
		&rawFlag, "raw", false,
		"if set, it does not parse nor indent the response",
	)
	flag.Parse()
}

func promptUnknown(key []byte) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "Please enter value of %s, parameter is empty: ", key)
	r := bufio.NewReader(os.Stdin)
	s, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return s[:len(s)-1], nil
}

func evaluateWithBinBash(key []byte) (out []byte, err error) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("echo %s", key))
	cmd.Env = os.Environ()

	out, err = cmd.Output()
	out = out[:len(out)-1]
	return
}

type Options struct {
	Source pipe
}

type pipe struct {
	*os.File
	keep bool
}

func (p pipe) Close() error {
	if p.keep {
		return nil
	}
	return p.File.Close()
}
