package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	EMPTY                               = ""
	ANNOTATE_HELMFILE_RELEASES_CMD_NAME = "annotate-helmfile-releases"
	COLLAPSE_HELM_DIFF_CMD_NAME         = "collapse-helm-diff"
	LISTEN_HELM_DEBUG_CMD_NAME       = "listen-helm-debug"
)

func commandUsage() string {
	return fmt.Sprintf("Subcommand must be one of: %s, %s or %s", ANNOTATE_HELMFILE_RELEASES_CMD_NAME, COLLAPSE_HELM_DIFF_CMD_NAME, LISTEN_HELM_DEBUG_CMD_NAME)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Please specify a subcommand. %s", commandUsage())
	}
	cmd, args := args[0], args[1:]
	switch cmd {
	case ANNOTATE_HELMFILE_RELEASES_CMD_NAME:
		AnnotateHelmfileReleases(args)
	case COLLAPSE_HELM_DIFF_CMD_NAME:
		CollapseHelmDiff(args)
	case LISTEN_HELM_DEBUG_CMD_NAME:
		ListenHelmDebug(args)
	default:
		log.Fatalf("Unrecognized subcommand %q. %s", cmd, commandUsage())
	}
}
