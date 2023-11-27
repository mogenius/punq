package structs

import (
	"os/exec"
	"time"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"
)

type Command struct {
	Id                      string  `json:"id"`
	JobId                   string  `json:"jobId"`
	ProjectId               string  `json:"projectId"`
	NamespaceId             *string `json:"namespaceId,omitempty"`
	ServiceId               *string `json:"serviceId,omitempty"`
	Title                   string  `json:"title"`
	Message                 string  `json:"message,omitempty"`
	StartedAt               string  `json:"startedAt"`
	State                   string  `json:"state"`
	DurationMs              int64   `json:"durationMs"`
	MustSucceed             bool    `json:"mustSucceed"`
	ReportToNotificationSvc bool    `json:"reportToNotificationService"`
	IgnoreError             bool    `json:"ignoreError"`
	BuildId                 int     `json:"buildId,omitempty"`
	Started                 time.Time
}

func ExecuteBashCommandSilent(title string, shellCmd string) {
	var err error
	_, err = utils.RunOnLocalShell(shellCmd).Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode := exitErr.ExitCode()
		errorMsg := string(exitErr.Stderr)
		logger.Log.Error(shellCmd)
		logger.Log.Errorf("%d: %s", exitCode, errorMsg)
	} else if err != nil {
		logger.Log.Errorf("ERROR: '%s': %s\n", title, err.Error())
	} else {
		if utils.CONFIG.Misc.Debug {
			logger.Log.Infof("SUCCESS '%s': %s\n", title, shellCmd)
		}
	}
}

func ExecuteBashCommandWithResponse(title string, shellCmd string) string {
	var err error
	var returnStr []byte
	_, err = utils.RunOnLocalShell(shellCmd).Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode := exitErr.ExitCode()
		errorMsg := string(exitErr.Stderr)
		logger.Log.Error(shellCmd)
		logger.Log.Errorf("%d: %s", exitCode, errorMsg)
		return errorMsg
	} else if err != nil {
		logger.Log.Errorf("ERROR: '%s': %s\n", title, err.Error())
		return err.Error()
	} else {
		if utils.CONFIG.Misc.Debug {
			logger.Log.Infof("SUCCESS '%s': %s\n", title, shellCmd)
		}
	}
	return string(returnStr)
}
