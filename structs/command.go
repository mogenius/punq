package structs

import (
	"os/exec"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"
)

func ExecuteShellCommandSilent(title string, shellCmd string) {
	var err error
	output, err := utils.RunOnLocalShell(shellCmd).Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode := exitErr.ExitCode()
		errorMsg := string(exitErr.Stderr)
		logger.Log.Error(shellCmd)
		logger.Log.Errorf("ExitCode: %d - Error: '%s' -> Output: '%s'", exitCode, errorMsg, string(output))
	} else if err != nil {
		logger.Log.Errorf("ERROR: '%s': %s\n", title, err.Error())
	} else {
		if utils.CONFIG.Misc.Debug {
			logger.Log.Infof("SUCCESS '%s': %s\n", title, shellCmd)
		}
	}
}

func ExecuteShellCommandWithResponse(title string, shellCmd string) string {
	var err error
	var returnStr []byte
	returnStr, err = utils.RunOnLocalShell(shellCmd).Output()
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
