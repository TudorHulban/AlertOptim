package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
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

func sameContent(filePath1, filePath2 string) error {
	content1, errRead1 := readFile(filePath1)
	if errRead1 != nil {
		return errRead1
	}

	c1, errPar1 := parseContent(content1)
	if errPar1 != nil {
		return errPar1
	}

	content2, errRead2 := readFile(filePath2)
	if errRead2 != nil {
		return errRead2
	}

	c2, errPar2 := parseContent(content2)
	if errPar2 != nil {
		return errPar2
	}

	if !reflect.DeepEqual(c1, c2) {
		return errors.New("file contents do NOT match")
	}

	return nil
}

func parseContent(c []string) (map[string]int, error) {
	if len(c) == 0 {
		return nil, errors.New("empty content")
	}

	res := make(map[string]int)

	for _, row := range c {
		if _, exists := res[row]; exists {
			res[row]++
			continue
		}

		res[row] = 1
	}

	return res, nil
}
