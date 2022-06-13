This tool compares the released version of a chart to the available versions.

# Configuration

Currently, configuration is read from a yaml-file, by default `autopatch.yaml` in the current
directory.

The config file has the following schema:

```yaml
---
charts:
  - repo: "http://example.com/charts"
    chart: "my-chart"
    release: "my-release"
    namespace: "my-namespace"
```

Run the check with `autpatch charts`.


# Feature backlog

- App version checks against github
    - How to define "current version"-logic?