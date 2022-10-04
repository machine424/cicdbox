package main

import (
	"fmt"
	"io"
)

const (
	COLLAPSED_SECTION_END = `\e[0Ksection_end:1664723888:section\r\e[0K`
)

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func collapsedSectionStart(header string) string {
	return fmt.Sprintf(`\e[0Ksection_start:1664723888:section[collapsed=true]\r\e[0K%s`, header)
}

func echoCollapsed(w io.Writer, header string, body []string) {
	echo(w, collapsedSectionStart(header))
	for _, line := range body {
		echo(w, line)
	}
	echo(w, COLLAPSED_SECTION_END)
}
