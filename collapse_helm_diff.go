package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	COLOR_RESET   = "[0m"
	RESOURCE_DIFF = "[0;33m"
	ELLIPSIS      = "..."

	COLLAPSED_SECTION_END = `\e[0Ksection_end:1664723888:section\r\e[0K`
)

var DIFF_COLOR_PREFIXES = []string{"[0;31m-", "[0;32m+"}

func collapsedSectionStart(header string) string {
	return fmt.Sprintf(`\e[0Ksection_start:1664723888:section[collapsed=true]\r\e[0K%s`, header)
}

func matchesPattern(str string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(str) {
			return true
		}
	}
	return false
}

func coloredWith(line string, prefix string) bool {
	return strings.HasPrefix(line, prefix) && strings.HasSuffix(line, COLOR_RESET)
}

func diffLine(line string) bool {
	for _, color := range DIFF_COLOR_PREFIXES {
		if coloredWith(line, color) {
			return true
		}
	}
	return false
}

func uncolor(line string) string {
	for _, color := range DIFF_COLOR_PREFIXES {
		if coloredWith(line, color) {
			return strings.TrimSuffix(strings.TrimPrefix(line, color), COLOR_RESET)
		}
	}
	return line
}

func echo(w io.Writer, str string) {
	cmd := exec.Command("echo", "-e", str)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func collapse(w io.Writer, diffs []string, patterns []*regexp.Regexp) {
	if len(diffs) == 0 {
		return
	}
	header, body := diffs[0], diffs[1:]
	if matchesPattern(uncolor(header), patterns) {
		header = ELLIPSIS
		body = diffs
	}
	echo(w, collapsedSectionStart(header))
	for _, line := range body {
		echo(w, line)
	}
	echo(w, COLLAPSED_SECTION_END)
}

func processResourceDiff(w io.Writer, diffs []string, patterns []*regexp.Regexp) {
	if len(diffs) == 0 {
		return
	}
	diffHeader, lines := diffs[0], diffs[1:]
	toBeCollapsed := []string{diffHeader}

	for _, line := range lines {
		if diffLine(line) {
			if matchesPattern(uncolor(line), patterns) {
				toBeCollapsed = append(toBeCollapsed, line)
			} else {
				//TODO: Keep more context (one line before?, what if this line matchesPattern)
				//if len(toBeCollapsed) != 0 {
				//	contextLine, rest := toBeCollapsed[len(toBeCollapsed)-1], toBeCollapsed[:len(toBeCollapsed)-1]
				//	collapse(rest)
				//	echo(contextLine)
				//}
				collapse(w, toBeCollapsed, patterns)
				toBeCollapsed = []string{}
				echo(w, line)
			}
		} else {
			toBeCollapsed = append(toBeCollapsed, line)
		}
	}
	collapse(w, toBeCollapsed, patterns)
}

func collapseHelmDiff(input *bufio.Scanner, output io.Writer, args []string) {
	flag := flag.NewFlagSet(COLLAPSE_HELM_DIFF_CMD_NAME, flag.ExitOnError)
	err := flag.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	regexes := flag.Args()
	if len(regexes) == 0 {
		log.Fatal("No regex was provided. Usage: command regex1 regex2 ...")
	}
	patterns := make([]*regexp.Regexp, len(regexes))
	for i, regex := range regexes {
		pattern, err := regexp.Compile(regex)
		if err != nil {
			log.Fatal(err)
		}
		patterns[i] = pattern
	}

	resourceDiff := []string{}
	for input.Scan() {
		rawLine := input.Text()
		if coloredWith(rawLine, RESOURCE_DIFF) {
			processResourceDiff(output, resourceDiff, patterns)
			resourceDiff = []string{}
		}
		resourceDiff = append(resourceDiff, rawLine)
	}
	processResourceDiff(output, resourceDiff, patterns)

	if err := input.Err(); err != nil {
		log.Println(err)
	}
}

func CollapseHelmDiff(args []string) {
	// read stdin
	input := bufio.NewScanner(os.Stdin)
	// write to stdout
	output := os.Stdout
	collapseHelmDiff(input, output, args)
}
