package autopatch

import (
	"fmt"
	"strings"
	"testing"
)

func TestListReleases(t *testing.T) {
	client := HelmClient{}

	releases := client.ListReleasesAll()
	fmt.Println(releases)
}

func TestAddAndUpdateRepo(t *testing.T) {
	client := HelmClient{}
	const dummyRepo = "dummy_repo_name_for_test"
	defer client.DeleteRepo(dummyRepo)
	client.AddAndUpdateRepo(dummyRepo, "https://airflow.apache.org")
	repos, _ := client.run("repo", "list")
	if !strings.Contains(string(repos), dummyRepo) {
		t.Errorf("Expected %s to be in the repo list", dummyRepo)
	}
}

func TestGetLatestChart(t *testing.T) {
	client := HelmClient{}
	const jupyterhub = "jupyterhub"
	chart := client.GetLatestChart(jupyterhub, "https://jupyterhub.github.io/helm-chart/")
	if chart.Name != jupyterhub {
		t.Errorf("Expected %s chart name, got %s", jupyterhub, chart.Name)
	}
}

func TestLoopAndCompare(t *testing.T) {
	charts := []ChartConfig{
		{
			Repo:      "https://airflow.apache.org",
			Chart:     "airflow",
			Release:   "airflow",
			Namespace: "systems-airflow",
		},
	}
	errors := LoopChartsAndCheck(charts)

	for _, err := range errors {
		fmt.Println(err)
	}
}
