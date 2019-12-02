package docker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	r = regexp.MustCompile(`COPY\s+--from=`)
)

// ContainsMultiStageBuilds 判斷 Dockerfile 是否包含 multi-stage builds
func ContainsMultiStageBuilds(log *logrus.Logger, pwd string) bool {
	p := filepath.Join(pwd, "Dockerfile")
	log.Debugf("loading Dockerfile: %s", p)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Debugf("error detecting is Dockerfile contains multi-stage builds: %v", err)
		return false
	}
	dockerfile := string(b)
	matches := r.MatchString(dockerfile)
	log.Debugf("detecting %s/Dockerfile contains multi-stage builds: %v", pwd, matches)
	return matches
}

// Build to exec 'docker build' command
func Build(log *logrus.Logger, image *SoftleaderHubImage) error {
	cmd := exec.Command("docker", "build", "-t", image.String(), ".")
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
