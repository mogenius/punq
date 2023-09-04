package structs

type Version struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	Branch         string `json:"branch"`
	BuildTimestamp string `json:"buildTimestamp"`
	GitCommitHash  string `json:"gitCommitHash"`
}

func VersionFrom(name string, version string, branch string, buildTimestamp string, gitCommitHash string) Version {
	return Version{
		Name:           name,
		Version:        version,
		Branch:         branch,
		BuildTimestamp: buildTimestamp,
		GitCommitHash:  gitCommitHash,
	}
}
