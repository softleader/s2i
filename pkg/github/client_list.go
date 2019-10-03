package github

import (
	"context"
	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"regexp"
)

// ListReleaseByRegex 依照 regex 列出符合的 release 資訊
func ListReleaseByRegex(log *logrus.Logger, token, owner, repo string, regex []*regexp.Regexp) error {
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return err
	}

	opt := &github.ListOptions{}
	for {
		release, resp, err := client.Repositories.ListReleases(ctx, owner, repo, nil)
		if err != nil {
			return err
		}
		if err := listReleasesByRegex(log, release, regex); err != nil {
			return err
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return nil
}

func listReleasesByRegex(log *logrus.Logger, releases []*github.RepositoryRelease, regex []*regexp.Regexp) error {
	for _, release := range releases {
		if anyMatch(regex, release.GetName()) {
			log.Infof("%s\t%s\t%s", release.GetName(), release.GetPublishedAt(), release.GetAuthor().GetLogin())
		}
	}
	return nil
}

// ListRelease 依照 release 名稱列出符合相關資訊
func ListRelease(log *logrus.Logger, token, owner, repo string, tags []string) error {
	ctx := context.Background()
	client, err := newTokenClient(ctx, token)
	if err != nil {
		return err
	}

	for _, tag := range tags {
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
		log.Infof("%s\t%s\t%s", rr.GetName(), rr.GetPublishedAt(), rr.GetAuthor().GetLogin())
	}
	return nil
}
