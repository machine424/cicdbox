package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
)

var debugLineRegex = regexp.MustCompile(`.+\[debug\]\s*(?P<debug_log>.*)$`)

// TODO: also warnings?
func printIfDebugLine(output io.Writer, line []byte) {
	result := debugLineRegex.FindSubmatch(line)
	if len(result) > 0 {
		debugLog := result[len(result)-1]
		_, err := output.Write(append(debugLog, []byte("\n")...))
		panicOnError(err)
	}
}

func listenHelmDebug(input *bufio.Scanner, output io.Writer, args []string) {
	flag := flag.NewFlagSet(LISTEN_HELM_DEBUG_CMD_NAME, flag.ExitOnError)
	err := flag.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	_, err = output.Write([]byte("\x1b[0;33mcicdbox: Only Helm debug logs will be displayed...\x1b[0m\n"))
	panicOnError(err)

	var outputToCollapse []string
	for input.Scan() {
		rawLine := input.Bytes()
		// Keep for the end
		outputToCollapse = append(outputToCollapse, string(rawLine))
		printIfDebugLine(output, rawLine)
	}
	echoCollapsed(output, "\x1b[0;33mExpand to see the entire output.\x1b[0m", outputToCollapse)

	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
}

func ListenHelmDebug(args []string) {
	// read stdin
	input := bufio.NewScanner(os.Stdin)
	// write to stdout
	output := os.Stdout
	listenHelmDebug(input, output, args)
}
