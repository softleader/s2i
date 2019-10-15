package mvn

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

// Test runs mvn test
func Test(log *logrus.Logger, configServer, configLabel string, updateSnapshots bool) error {
	args := []string{"clean", "test", "-Dspring.profiles.active=test", "-Dspring.cloud.config.uri=" + configServer}
	if configLabel != "" {
		args = append(args, "-Dspring.cloud.config.label="+configLabel)
	}
	if updateSnapshots {
		args = append(args, "-U")
	}
	cmd := exec.Command("mvn", args...)
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
