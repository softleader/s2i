package github

import (
	"context"
	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
)

// CreateRelease 建立 github 的 release
func CreateRelease(log *logrus.Logger, token, owner, repo, branch, tag string) error {
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
		return err
	}
	log.Printf("Successfully created release: %s", release.GetHTMLURL())
	return nil
}

// CreatePrerelease 建立 github 的 pre-release
func CreatePrerelease(log *logrus.Logger, token, owner, repo, branch, tag string, force bool) error {
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return err
	}
	pre := true
	r := &github.RepositoryRelease{
		TagName:         &tag,
		TargetCommitish: &branch,
		Prerelease:      &pre,
	}
	log.Debugf("creating pre-release %s for %s/%s branch: %s", tag, owner, repo, branch)
	release, _, err := client.Repositories.CreateRelease(ctx, owner, repo, r)
	if err != nil {
		githubErr, ok := err.(*github.ErrorResponse)
		if !ok {
			return err
		}
		if force && isTagNameAlreadyExists(githubErr.Errors) {
			log.Debugf("tag name %s already exists, force to delete it..", tag)
			if err := deleteReleaseAndTag(ctx, log, client, owner, repo, tag, false); err != nil {
				return err
			}
		}
		log.Debugf("creating pre-release %s again for %s/%s branch: %s", tag, owner, repo, branch)
		if release, _, err = client.Repositories.CreateRelease(ctx, owner, repo, r); err != nil {
			return err
		}
	}

	log.Printf("Successfully created pre-release: %s", release.GetHTMLURL())
	return nil
}

func isTagNameAlreadyExists(errors []github.Error) bool {
	for _, err := range errors {
		if err.Field == "tag_name" && err.Code == "already_exists" {
			return true
		}
	}
	return false
}
