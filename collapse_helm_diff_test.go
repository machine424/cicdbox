package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func absolutePath(name string) string {
	return fmt.Sprintf("files/collapse_helm_diff/%s", name)
}

func openFile(name string) *os.File {
	file, err := os.Open(name)
	check(err)
	return file
}

func echoLines(input *bufio.Scanner, output io.Writer) {
	for input.Scan() {
		echo(output, input.Text())
	}
}

func TestEcho(t *testing.T) {
	input := "foo-bar"
	expect := fmt.Sprintf("%s\n", input)
	var output strings.Builder
	echo(&output, input)
	if output.String() != expect {
		t.Errorf("got %s but expected %s", output.String(), expect)
	}
}

func TestCollapseHelmDiff(t *testing.T) {
	tt := []struct {
		description string
		diff        string
		args        []string
		expect      string
	}{
		{"", absolutePath("diff_1.in"), []string{`bar`}, absolutePath("diff_1.outt")},
		{"", absolutePath("diff_2.in"), []string{`^\s+-?image: docker\.foo\.fr`}, absolutePath("diff_2.outt")},
		{"", absolutePath("diff_3.in"), []string{`^\s+-?image: docker\.foo\.fr`}, absolutePath("diff_3.outt")},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			// input
			diffInFile := openFile(tc.diff)
			defer diffInFile.Close()
			input := bufio.NewScanner(diffInFile)

			// expected output
			rawExpectInFile := openFile(tc.expect)
			defer rawExpectInFile.Close()
			var expect strings.Builder
			echoLines(bufio.NewScanner(rawExpectInFile), &expect)

			// run
			var output strings.Builder
			collapseHelmDiff(input, &output, tc.args)
			if expect.String() != output.String() {
				t.Errorf("got %s but expected %s", output.String(), expect.String())
			}
		})
	}
}
