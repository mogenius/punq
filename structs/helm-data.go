package structs

import "time"

type HelmData struct {
	APIVersion string                 `yaml:"apiVersion"`
	Entries    map[string][]HelmEntry `yaml:"entries"`
}

type HelmEntry struct {
	APIVersion   string           `yaml:"apiVersion"`
	AppVersion   string           `yaml:"appVersion"`
	Dependencies []HelmDependency `yaml:"dependencies"`
	Created      time.Time        `yaml:"created"`
	Description  string           `yaml:"description"`
	Digest       string           `yaml:"digest"`
	Name         string           `yaml:"name"`
	Urls         []string         `yaml:"urls"`
	Version      string           `yaml:"version"`
}

type HelmDependency struct {
	Name       string `yaml:"name"`
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}
