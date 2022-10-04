package main

import (
	"strings"
	"testing"
)

func TestListenHelmDebug(t *testing.T) {
	tt := []struct {
		description string
		diff        string
		expect      string
	}{
		{"debug logs", absoluteTestPath("helm_debug_1.in"), absoluteTestPath("helm_debug_1.outt")},
		{"no debug logs", absoluteTestPath("helm_debug_2.in"), absoluteTestPath("helm_debug_2.outt")},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			input, expect := loadTestFiles(tc.diff, tc.expect)

			var output strings.Builder
			listenHelmDebug(input, &output, []string{})
			assertEqual(t, expect, output)
		})
	}
}
