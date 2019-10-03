package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"regexp"
)

// DeleteReleasesAndTagsByRegex 依照 regex 刪除符合的 release 及其 tag
func DeleteReleasesAndTagsByRegex(log *logrus.Logger, token, owner, repo string, regex []*regexp.Regexp, dryRun bool) error {
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return err
	}

	opt := &github.ListOptions{}
	for {
		tags, resp, err := client.Repositories.ListTags(ctx, owner, repo, nil)
		if err != nil {
			return err
		}
		if err := deleteReleasesAndTagsByRegex(ctx, log, client, owner, repo, tags, regex, dryRun); err != nil {
			return err
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return nil
}

func deleteReleasesAndTagsByRegex(ctx context.Context, log *logrus.Logger, client *github.Client, owner, repo string, tags []*github.RepositoryTag, regex []*regexp.Regexp, dryRun bool) error {
	for _, tag := range tags {
		if anyMatch(regex, tag.GetName()) {
			log.Infof("one the of regex matches '%s', start to delete it...", tag.GetName())
			if err := deleteReleaseAndTag(ctx, log, client, owner, repo, tag.GetName(), dryRun); err != nil {
				return err
			}
			log.Infof("'%s' has been deleted from GitHub", tag.GetName())
		}
	}
	return nil
}

func anyMatch(regex []*regexp.Regexp, s string) bool {
	for _, r := range regex {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

// DeleteReleasesAndTags 刪除多筆 release 及其 refs/tag
func DeleteReleasesAndTags(log *logrus.Logger, token, owner, repo string, tags []string, dryRun bool) error {
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if err := deleteReleaseAndTag(ctx, log, client, owner, repo, tag, dryRun); err != nil {
			return err
		}
	}
	return nil
}

// DeleteReleaseAndTag 刪除 release 及其 refs/tag
func deleteReleaseAndTag(ctx context.Context, log *logrus.Logger, client *github.Client, owner, repo, tag string, dryRun bool) error {
	if err := deleteRelease(ctx, log, client, owner, repo, tag, dryRun); err != nil {
		return err
	}
	return deleteTag(ctx, log, client, owner, repo, tag, dryRun)
}

func deleteTag(ctx context.Context, log *logrus.Logger, client *github.Client, owner, repo, tag string, dryRun bool) error {
	log.Debugf("deleting refs/tags %s", tag)
	if !dryRun {
		_, err := client.Git.DeleteRef(ctx, owner, repo, fmt.Sprintf("tags/%s", tag))
		if err != nil {
			githubErr, ok := err.(*github.ErrorResponse)
			if !ok {
				return err
			}
			if githubErr.Response.StatusCode == 422 { // 代表 ref 刪除動作失敗, 直接中斷不丟錯
				return nil
			}
			return err
		}
	}
	return nil
}
func deleteRelease(ctx context.Context, log *logrus.Logger, client *github.Client, owner, repo, tag string, dryRun bool) error {
	log.Debugf("fetching release-id of tag '%s'", tag)
	rr, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		githubErr, ok := err.(*github.ErrorResponse)
		if !ok {
			return err
		}
		if githubErr.Response.StatusCode == 404 { // 代表 release 不存在, 直接中斷不丟錯
			return nil
		}
		return err
	}
	log.Debugf("deleting release %s by release-id %d", tag, rr.GetID())
	if !dryRun {
		_, err = client.Repositories.DeleteRelease(ctx, owner, repo, rr.GetID())
		if err != nil {
			return err
		}
	}
	return nil
}
