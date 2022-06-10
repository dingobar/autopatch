package helmclient

import (
	"encoding/json"
	"log"
	"os/exec"
)

type HelmClient struct {
}

type Release struct {
	name      string
	namespace string
}

func (client *HelmClient) run(cmd ...string) (output string, err error) {
	out, err := exec.Command("helm", cmd...).Output()
	if err != nil {
		log.Fatalf("Command %s failed with error %s", cmd, err)
	}
	outString := string(out)
	return outString, err
}

func (client *HelmClient) ListReleases(namespace string) []map[string]string {
	args := []string{"list", "-o", "json"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	releases, _ := client.run(args...)
	var out []map[string]string
	json.Unmarshal([]byte(releases), &out)
	return out
}

func (client *HelmClient) ListReleasesAll() []map[string]string {
	return client.ListReleases("")
}

func (client *HelmClient) AddAndUpdateRepo(name string, uri string) {
	//Add repo if it didn't already exist
}
