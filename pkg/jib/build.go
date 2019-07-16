package jib

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
	ur = regexp.MustCompile(`-Djib.to.auth.username=(\w+)`)
	pr = regexp.MustCompile(`-Djib.to.auth.password=(\w+)`)
)

// DockerBuild from jib
func DockerBuild(log *logrus.Logger, image, tag string) error {
	cmd := exec.Command("mvn", "compile", "jib:dockerBuild", "-Dbuild.image="+image, "-Dbuild.tag="+tag)
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

// Auth for jib:build
type Auth struct {
	Username string
	Password string
}

// Build from jib
func Build(log *logrus.Logger, image, tag string, auth *Auth) error {
	cmd := exec.Command("mvn", "compile", "jib:build", "-Djib.to.auth.username="+auth.Username, "-Djib.to.auth.password="+auth.Password, "-Dbuild.image="+image, "-Dbuild.tag="+tag)
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

// GetAuth 會試著從 Jenkinsfile 取得帳密, 因為我們通常是放在 Jenkinsfile 中
func GetAuth(log *logrus.Logger, pwd string) (auth *Auth) {
	auth = &Auth{}
	p := filepath.Join(pwd, "Jenkinsfile")
	log.Debugf("loading Jenkinsfile: %s", p)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return
	}
	jenkinsfile := string(b)
	groups := ur.FindStringSubmatch(jenkinsfile)
	if len(groups) < 1 {
		return
	}
	log.Debugf("found jib.to.auth.username: %s", groups[1])
	auth.Username = groups[1]

	groups = pr.FindStringSubmatch(jenkinsfile)
	if len(groups) < 1 {
		return
	}
	log.Debugf("found jib.to.auth.password: %s", groups[1])
	auth.Password = groups[1]
	return
}
