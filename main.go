package main

import (
	"bufio"
	"io"
	"os"
)

const configFile = "funds.conf"

var funds []string

func init() {
	f, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		funds = append(funds, string(line))
	}
}

func main() {
	GetFunds(funds)
}
