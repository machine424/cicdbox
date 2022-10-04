package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	HELMFILE_BIN = "helmfile"
	KUBECTL_BIN  = "kubectl"

	NAMESPACE_KEY = "NAMESPACE: "
	NAME_KEY      = "NAME: "
	REVISION_KEY  = "REVISION: "
)

var (
	GITLAB_USER_LOGIN = os.Getenv("GITLAB_USER_LOGIN")
	CI_PIPELINE_URL   = os.Getenv("CI_PIPELINE_URL")
)

type namespacedSecret struct {
	secretName string
	namespace  string
}

// Runs the command name arg... Returns its Stderr if it fails or its stringified Stdout
func RunCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	return stdout.String(), nil
}

// Adds annotations to the Secret
func annotateSecret(secret string, namespace string) {
	args := []string{"annotate", "secret", secret, "--namespace", namespace, "--overwrite=true", fmt.Sprintf("cicdbox/releaser=%s", GITLAB_USER_LOGIN), fmt.Sprintf("cicdbox/gitlab-ci-pipeline-url=%s", CI_PIPELINE_URL)}

	_, err := RunCommand(KUBECTL_BIN, args...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Secret %s/%s was annotated.\n", namespace, secret)
}

// Infers the Helm Secrets corresponding to the revisions of the given releases
func retrieveSecrets(releases string) []namespacedSecret {
	var secretList []namespacedSecret

	var releaseNamespace string
	var releaseName string

	for _, line := range strings.Split(releases, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, NAMESPACE_KEY) {
			releaseNamespace = strings.TrimPrefix(line, NAMESPACE_KEY)
		} else if strings.HasPrefix(line, NAME_KEY) {
			releaseName = strings.TrimPrefix(line, NAME_KEY)
		} else if strings.HasPrefix(line, REVISION_KEY) {
			releaseRevision := strings.TrimPrefix(line, REVISION_KEY)
			// We suppose REVISION_KEY is the last key.
			secretList = append(secretList, namespacedSecret{namespace: releaseNamespace, secretName: fmt.Sprintf("sh.helm.release.v1.%s.v%s", releaseName, releaseRevision)})
		}
	}
	return secretList
}

func AnnotateHelmfileReleases(args []string) {
	flag := flag.NewFlagSet(ANNOTATE_HELMFILE_RELEASES_CMD_NAME, flag.ExitOnError)
	helmfileEnvironment := flag.String("environment", EMPTY, "Helmfile environment (Required)")
	helmfileSelector := flag.String("selector", EMPTY, "Helmfile selector (Required)")
	file := flag.String("file", EMPTY, "Helmfile file (Required)")
	err := flag.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	helmfileArgs := []string{}
	if *helmfileEnvironment != EMPTY {
		helmfileArgs = append(helmfileArgs, "--environment", *helmfileEnvironment)
	}
	if *helmfileSelector != EMPTY {
		helmfileArgs = append(helmfileArgs, "--selector", *helmfileSelector)
	}
	if *file != EMPTY {
		helmfileArgs = append(helmfileArgs, "--file", *file)
	}
	helmfileArgs = append(helmfileArgs, "status")

	// Retrieve releases status
	output, err := RunCommand(HELMFILE_BIN, helmfileArgs...)
	if err != nil {
		log.Fatal(err)
	}

	// Annotate secrets
	for _, secret := range retrieveSecrets(output) {
		annotateSecret(secret.secretName, secret.namespace)
	}

}
