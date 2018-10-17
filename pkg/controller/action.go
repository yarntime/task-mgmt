package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yarntime/task-mgmt/pkg/types"
	"os"
	"os/exec"
	"path"
	"strings"
)

func runCommand(args []string, input []byte) ([]byte, error) {
	if os.Getenv("CB_ENVDIR") == "" {
		return nil, errors.New("CB_ENVDIR is missing")
	}
	binDir := os.Getenv("CB_BINDIR")
	if binDir == "" {
		binDir = path.Join(path.Dir(os.Getenv("CB_ENVDIR")), "bin")
		if _, err := os.Stat(binDir); os.IsNotExist(err) {
			return nil, errors.New(" The CB_BINDIR " + binDir + " does not exist")
		}
	}
	aipPath := path.Join(binDir, "aip")
	if _, err := os.Stat(aipPath); err != nil {
		return nil, err
	}
	cmd := exec.Command(aipPath, args...)
	if len(input) > 0 {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		go func() {
			defer stdin.Close()
			stdin.Write(input)
		}()
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return output, nil
}

// CreateCron submits a cron job to AIP
func CreateCron(spec []byte) error {
	output, err :=
		runCommand([]string{"cronjob", "create", "--overwrite", "-"}, spec)
	if err != nil {
		return err
	}
	fmt.Print(string(output))
	return nil
}

// DeleteCron deletes an existing cron job in AIP
// To delete all cron job, specify "all" as the name
func DeleteCron(name string) error {
	output, err :=
		runCommand([]string{"cronjob", "remove", name}, nil)
	if err != nil {
		return err
	}
	fmt.Print(string(output))
	return nil
}

// GetCron gets the information of submitted cron job
func GetCron(name string) ([]*types.CronInfo, error) {
	var commands []string
	if name == "" {
		commands = []string{"cronjob", "info", "--long"}
	} else {
		commands = []string{"cronjob", "info", "--long", name}
	}
	output, err := runCommand(commands, nil)
	if err != nil {
		return nil, err
	}
	var cronInfoL []*types.CronInfo
	if strings.Contains(string(output), "No cron job is found") {
		return cronInfoL, nil
	}
	if err := json.Unmarshal(output, &cronInfoL); err != nil {
		return nil, err
	}
	return cronInfoL, nil
}
