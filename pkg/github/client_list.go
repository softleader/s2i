package github

import (
	"context"
	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
)

// ListReleaseByMatcher 依照指定 matcher 列出符合的 release 資訊
func ListReleaseByMatcher(log *logrus.Logger, token, owner, repo string, matcher TagMatcher) error {
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
		releases, resp, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
		if err != nil {
			return err
		}
		for _, release := range releases {
			if matcher.Matches(release.GetName()) {
				log.Infof("%s\t%s\t%s", release.GetName(), release.GetPublishedAt(), release.GetAuthor().GetLogin())
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
