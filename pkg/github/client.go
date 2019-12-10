package github

import (
	"context"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	r = regexp.MustCompile(`\[remote "origin"\][\n|\r|\n\r|\t|\s]+url = [https://|git@]+([^@]+)?@?github.com[/:](.+)/(.+).git`)
)

// NewTokenClient 建立跟 github 互動的 client
func newTokenClient(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

// FindNextReleaseVersion 找下一版 revision,  也就是 latest release + 1 版本號
func FindNextReleaseVersion(log *logrus.Logger, token, owner, repo string) (string, error) {
	if token == "" || owner == "" || repo == "" {
		return "", nil
	}
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return "", err
	}
	log.Debugf("fetching latest release of %s/%s", owner, repo)
	rr, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return "", err
	}
	tag := rr.GetTagName()
	log.Debugf("found %s drafted by %s published at %s", tag, rr.GetAuthor().GetLogin(), rr.GetPublishedAt())
	version := strings.TrimPrefix(tag, "v")
	sv, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	sv.BumpPatch()
	next := sv.String()
	if strings.HasPrefix(tag, "v") {
		next = "v" + next
	}
	return next, nil

}

// Remote 回傳從 .git 中找到的 token, owner and repo
func Remote(log *logrus.Logger, pwd string) (token, owner, repo string) {
	p := filepath.Join(pwd, ".git", "config")
	log.Debugf("loading git config: %s", p)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return
	}
	config := string(b)
	return findRemoteOrigin(log, config)
}

func findRemoteOrigin(log *logrus.Logger, config string) (token, owner, repo string) {
	groups := r.FindStringSubmatch(config)
	log.Debugf("found %d remote origin", len(groups)-1)
	if len(groups) < 1 {
		return
	}
	token = groups[1]
	owner = groups[2]
	repo = groups[3]
	log.Debugf("found token %q, owner: %q, repo: %q", token, owner, repo)
	return
}

// Head 回傳當前的 branch
func Head(log *logrus.Logger, pwd string) string {
	p := filepath.Join(pwd, ".git", "HEAD")
	log.Debugf("loading git HEAD: %s", p)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return ""
	}
	head := string(b)
	lines := strings.Split(head, fmt.Sprintln())
	if len(lines) < 1 {
		return ""
	}
	return strings.ReplaceAll(lines[0], "ref: refs/heads/", "")
}
