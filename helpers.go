package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
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

func sameContent(filePath1, filePath2 string, w io.Writer) error {
	content1, errRead1 := readFile(filePath1)
	if errRead1 != nil {
		return errRead1
	}

	c1, errPar1 := parseContentToLower(content1)
	if errPar1 != nil {
		return errPar1
	}

	content2, errRead2 := readFile(filePath2)
	if errRead2 != nil {
		return errRead2
	}

	c2, errPar2 := parseContentToLower(content2)
	if errPar2 != nil {
		return errPar2
	}

	if reflect.DeepEqual(c1, c2) {
		return nil
	}

	msg1 := fmt.Sprintf("\nAnalysis - Diffs in file: '%s' versus file: '%s':\n", filePath1, filePath2)
	w.Write([]byte(msg1))
	diff(c1, c2, w)

	msg2 := fmt.Sprintf("\nAnalysis - Diffs in file: '%s' versus file: '%s':\n", filePath2, filePath1)
	w.Write([]byte(msg2))
	diff(c2, c1, w)

	return errors.New("file contents do NOT match")
}

func diff(d1, d2 map[string]int, w io.Writer) {
	var diffs []string

	for k := range d1 {
		if k == " " || k == "\n" {
			continue
		}

		if _, exists := d2[k]; !exists {
			diffs = append(diffs, "missing: "+k)
			continue
		}

		if d1[k] != d2[k] {
			diffs = append(diffs, strings.ReplaceAll(k, "\n", "")+" -- "+strconv.Itoa(d1[k])+" versus "+strconv.Itoa(d2[k])+"\n")
		}
	}

	if len(diffs) == 0 {
		w.Write([]byte("no differences in content.\n"))
		return
	}

	sort.Strings(diffs)
	w.Write([]byte(strings.Join(diffs, "")))
}

func parseContentToLower(c []string) (map[string]int, error) {
	if len(c) == 0 {
		return nil, errors.New("empty content")
	}

	res := make(map[string]int)

	for _, row := range c {
		buf := strings.ToLower(row)

		if _, exists := res[buf]; exists {
			res[buf]++
			continue
		}

		res[buf] = 1
	}

	return res, nil
}

func sortAlerts(data []Alert) []Alert {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	return data
}
