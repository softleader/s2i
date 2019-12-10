package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
)

// DeleteMatchesReleasesAndTags 刪除所有符合的 release 及其 tag
func DeleteMatchesReleasesAndTags(log *logrus.Logger, token, owner, repo string, matcher TagMatcher, dryRun bool) error {
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return err
	}

	opt := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	for {
		log.Debugf("fetching page %v of tags", opt.Page)
		tags, resp, err := client.Repositories.ListTags(ctx, owner, repo, opt)
		if err != nil {
			return err
		}
		for _, tag := range tags {
			if matcher.Matches(tag.GetName()) {
				log.Infof("'%s' matches! start to delete it...", tag.GetName())
				if err := deleteReleaseAndTag(ctx, log, client, owner, repo, tag.GetName(), dryRun); err != nil {
					return err
				}
				log.Infof("'%s' has been deleted from GitHub", tag.GetName())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		log.Debugf("moving to the next page: %v/%v", resp.NextPage, resp.LastPage)
		opt.Page = resp.NextPage
	}

	return nil
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
