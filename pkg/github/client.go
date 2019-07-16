package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	r = regexp.MustCompile(`url = (.+).git`)
)

func newTokenClient(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

// CreateRelease 建立 github 的 release
func CreateRelease(log *logrus.Logger, token, owner, repo, branch, tag string, force bool) error {
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return err
	}

	r := &github.RepositoryRelease{
		TagName:         &tag,
		TargetCommitish: &branch,
	}
	log.Debugf("creating release %s for %s/%s branch: %s", tag, owner, repo, branch)
	release, _, err := client.Repositories.CreateRelease(ctx, owner, repo, r)
	if err != nil {
		githubErr, ok := err.(*github.ErrorResponse)
		if !ok {
			return err
		}
		if force && isTagNameAlreadyExists(githubErr.Errors) {
			log.Debugf("tag name %s already exists, force to delete it..", tag)
			if err := deleteReleaseByName(ctx, client, owner, repo, tag); err != nil {
				return err
			}
		}
		log.Debugf("creating release %s again for %s/%s branch: %s", tag, owner, repo, branch)
		if release, _, err = client.Repositories.CreateRelease(ctx, owner, repo, r); err != nil {
			return err
		}
	}

	log.Printf("Successfully created %s release", release.GetName())
	return nil
}

func deleteReleaseByName(ctx context.Context, client *github.Client, owner, repo, tag string) error {
	rr, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return err
	}
	_, err = client.Repositories.DeleteRelease(ctx, owner, repo, rr.GetID())
	if err != nil {
		return err
	}
	_, err = client.Git.DeleteRef(ctx, owner, repo, fmt.Sprintf("tags/%s", tag))
	return err
}

func isTagNameAlreadyExists(errors []github.Error) bool {
	for _, err := range errors {
		if err.Field == "tag_name" && err.Code == "already_exists" {
			return true
		}
	}
	return false
}

// Remote 回傳預設的 owner and repo
func Remote(log *logrus.Logger, pwd string) (owner, repo string) {
	p := filepath.Join(pwd, ".git", "config")
	log.Debugf("loading git config: %s", p)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return
	}
	config := string(b)
	groups := r.FindStringSubmatch(config)
	log.Debugf("found %d remote url", len(groups)-1)
	if len(groups) < 1 {
		return
	}
	remote := groups[1]
	remote = strings.ReplaceAll(remote, "git@github.com:", "")
	remote = strings.ReplaceAll(remote, "https://github.com/", "")
	log.Debugf("used remote url: %s", remote)
	spited := strings.Split(remote, "/")
	owner = spited[0]
	repo = spited[1]
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
