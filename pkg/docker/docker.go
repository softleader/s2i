package docker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

// Push to exec 'docker push' command
func Push(log *logrus.Logger, image *SoftleaderHubImage) error {
	cmd := exec.Command("docker", "push", image.String())
	if log.IsLevelEnabled(logrus.DebugLevel) {
		log.Out.Write([]byte(fmt.Sprintln(strings.Join(cmd.Args, " "))))
	}
	cmd.Stdout = log.Out
	cmd.Stderr = log.Out
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// Rmi to exec 'docker rmi' command
func Rmi(log *logrus.Logger, image *SoftleaderHubImage) error {
	cmd := exec.Command("docker", "rmi", image.String())
	if log.IsLevelEnabled(logrus.DebugLevel) {
		log.Out.Write([]byte(fmt.Sprintln(strings.Join(cmd.Args, " "))))
	}
	cmd.Stdout = log.Out
	cmd.Stderr = log.Out
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
