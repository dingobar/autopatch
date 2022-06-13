package autopatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type HelmClient struct {
}

type Release struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Chart       string `json:"chart"`
	App_version string `json:"app_version"`
}

func (release *Release) Version() string {
	version_parts := strings.Split(release.Chart, "-")
	return strings.Join(version_parts[1:], "-")
}

type Chart struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	App_version string `json:"app_version"`
}
type ChartConfig struct {
	Repo      string `mapstructure:"repo"`
	Chart     string `mapstructure:"chart"`
	Release   string `mapstructure:"release"`
	Namespace string `mapstructure:"namespace"`
}

func (client *HelmClient) run(cmd ...string) (output []byte, err error) {
	out, err := exec.Command("helm", cmd...).Output()
	if err != nil {
		logrus.Fatalf("Command %s failed with error %s", cmd, err)
	}
	return out, err
}

func (client *HelmClient) ListReleases(namespace string) []Release {
	args := []string{"list", "-o", "json"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	releases, _ := client.run(args...)
	var out []Release
	json.Unmarshal(releases, &out)
	return out
}

func (client *HelmClient) ListReleasesAll() []Release {
	return client.ListReleases("")
}

func (client *HelmClient) GetReleaseByName(name, namespace string) Release {
	releases := client.ListReleases(namespace)
	for _, release := range releases {
		if release.Name == name {
			return release
		}
	}
	logrus.Warningf("Release with name %s was not found in %s", name, namespace)
	return Release{}
}

func (client *HelmClient) AddAndUpdateRepo(name string, uri string) {
	client.run("repo", "add", name, uri)
	client.run("repo", "update", name)
}

func (client *HelmClient) DeleteRepo(name string) {
	client.run("repo", "remove", name)
}

func (client *HelmClient) GetLatestChart(chart, uri string) Chart {
	repoTempName := "temp_" + chart
	client.AddAndUpdateRepo(repoTempName, uri)
	defer client.DeleteRepo(repoTempName)
	return client.searchChartRepo(repoTempName, chart)
}

func (client *HelmClient) searchChartRepo(repo, chart string) Chart {
	repos, _ := client.run("search", "repo", repo+"/"+chart, "-o", "json")
	var out []Chart
	json.Unmarshal(repos, &out)
	return out[0]
}

func LoopChartsAndCheck(charts []ChartConfig) []error {
	client := HelmClient{}
	var versionErrors []error
	for _, chart := range charts {
		// Get desired version
		desired := client.GetLatestChart(chart.Chart, chart.Repo)
		// Get actual version
		actual := client.GetReleaseByName(chart.Release, chart.Namespace)
		if actual.Name == "" {
			continue
		}
		// Compare
		if desired.Version != actual.Version() {
			versionErrors = append(versionErrors, errors.New(fmt.Sprintf("Release %s in %s is version %s, but %s is available in %s", actual.Name, actual.Namespace, actual.Version(), desired.Version, chart.Repo)))
		}
	}

	return versionErrors
}
