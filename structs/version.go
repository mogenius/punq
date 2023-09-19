package structs

type Version struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	OperatorImage  string `json:"operatorImage"`
	Branch         string `json:"branch"`
	BuildTimestamp string `json:"buildTimestamp"`
	GitCommitHash  string `json:"gitCommitHash"`
}

func VersionFrom(name string, version string, branch string, buildTimestamp string, gitCommitHash string, operatorImage string) Version {
	return Version{
		Name:           name,
		Version:        version,
		OperatorImage:  operatorImage,
		Branch:         branch,
		BuildTimestamp: buildTimestamp,
		GitCommitHash:  gitCommitHash,
	}
}
