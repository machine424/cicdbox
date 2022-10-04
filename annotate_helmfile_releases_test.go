package main

import (
	"reflect"
	"testing"
)

func TestAnnotateHelmfileReleases(t *testing.T) {
	tt := []struct {
		description string
		releases    string
		expect      []namespacedSecret
	}{
		{"multiple releases", `NAME: foo
		LAST DEPLOYED: Fri Sep 30 19:44:07 2022
		NAMESPACE: ns
		STATUS: deployed
		REVISION: 1
		TEST SUITE: None
		
		NAME: bar
		LAST DEPLOYED: Fri Sep 30 19:44:07 2022
		NAMESPACE: ns
		STATUS: deployed
		REVISION: 100
		TEST SUITE: None
		
		`, []namespacedSecret{
			{namespace: "ns", secretName: "sh.helm.release.v1.foo.v1"},
			{namespace: "ns", secretName: "sh.helm.release.v1.bar.v100"},
		}},
		{"single release", `NAME: bar
		LAST DEPLOYED: Fri Sep 30 19:44:07 2022
		NAMESPACE: foo
		STATUS: failed
		REVISION: 5
		TEST SUITE: None
		
		`, []namespacedSecret{{namespace: "foo", secretName: "sh.helm.release.v1.bar.v5"}}},
		{"no release", ``, []namespacedSecret(nil)},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			output := retrieveSecrets(tc.releases)
			if !reflect.DeepEqual(tc.expect, output) {
				t.Errorf("got %#v but expected %#v", output, tc.expect)
			}
		})
	}
}
