package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func readFile(path string) ([]string, error) {
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, errOpen
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var res []string

	for scanner.Scan() {
		res = append(res, scanner.Text()+"\n")
	}
	if errScan := scanner.Err(); errScan != nil {
		return nil, errScan
	}

	return res, nil
}

func chopLastRow(path, verify string) error {
	contents, errRead := readFile(path)
	if errRead != nil {
		return errRead
	}

	if len(contents) == 0 {
		return errors.New("empty file")
	}

	if contents[len(contents)-1] != verify {
		fmt.Println("last row:", contents[len(contents)-1])
		return errors.New("last row does not match verification string")
	}

	return ioutil.WriteFile(path, []byte(strings.Join(contents[:len(contents)-1], "")), 0644)
}

func startPos(s string) int {
	if len(s) == 0 {
		return -1
	}

	i := 0
	for i < len(s)-1 {
		if s[i:i+1] != " " {
			break
		}

		i++
	}

	return i
}
