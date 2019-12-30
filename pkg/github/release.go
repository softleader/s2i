package github

import "github.com/google/go-github/v28/github"

// Release wrap GitHub Repository Release
type Release struct {
	TagName         string
	TargetCommitish string
	Name            string
	Draft           bool
	Prerelease      bool

	PublishedAt github.Timestamp
	HTMLURL     string
	Author      *github.User
}

func newRelease(rr *github.RepositoryRelease) *Release {
	return &Release{
		TagName:         rr.GetTagName(),
		TargetCommitish: rr.GetTargetCommitish(),
		Name:            rr.GetName(),
		Draft:           rr.GetDraft(),
		Prerelease:      rr.GetPrerelease(),
		PublishedAt:     rr.GetPublishedAt(),
		HTMLURL:         rr.GetHTMLURL(),
		Author:          rr.GetAuthor(),
	}
}
