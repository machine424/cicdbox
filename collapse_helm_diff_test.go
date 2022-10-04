package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func absoluteTestPath(name string) string {
	return fmt.Sprintf("files/%s", name)
}

func openFile(name string) *os.File {
	file, err := os.Open(name)
	panicOnError(err)
	return file
}

func echoLines(input *bufio.Scanner, output io.Writer) {
	for input.Scan() {
		echo(output, input.Text())
	}
}

func loadTestFiles(inputFilePath, expectFilePath string) (*bufio.Scanner, strings.Builder) {
	// input (we d not close it explicitly)
	inputFile := openFile(inputFilePath)
	input := bufio.NewScanner(inputFile)

	// expected output
	expectFile := openFile(expectFilePath)
	defer expectFile.Close()
	var expect strings.Builder
	echoLines(bufio.NewScanner(expectFile), &expect)

	return input, expect
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

func assertEqual(t *testing.T, expect, output strings.Builder) {
	if expect.String() != output.String() {
		t.Errorf("got %s but expected %s", output.String(), expect.String())
	}
}

func TestCollapseHelmDiff(t *testing.T) {
	tt := []struct {
		description string
		diff        string
		args        []string
		expect      string
	}{
		{"no diff", absoluteTestPath("diff_1.in"), []string{`bar`}, absoluteTestPath("diff_1.outt")},
		{"diff to callapse", absoluteTestPath("diff_2.in"), []string{`^\s+-?image: docker\.foo\.fr`}, absoluteTestPath("diff_2.outt")},
		{"diff to keep", absoluteTestPath("diff_3.in"), []string{`^\s+-?image: docker\.foo\.fr`}, absoluteTestPath("diff_3.outt")},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			input, expect := loadTestFiles(tc.diff, tc.expect)

			var output strings.Builder
			collapseHelmDiff(input, &output, tc.args)
			assertEqual(t, expect, output)
		})
	}
}
