package github

import (
	"github.com/coreos/go-semver/semver"
	"regexp"
	"strings"
)

// TagMatcher 判斷是否有 match tags
type TagMatcher interface {
	// Matches 傳入的 tag 是否匹配
	Matches(tag string) bool
}

// NewRegexMatcher 建立 RegexMatcher 物件
func NewRegexMatcher(expressions []string) *RegexMatcher {
	m := &RegexMatcher{}
	for _, expression := range expressions {
		m.regexps = append(m.regexps, regexp.MustCompile(expression))
	}
	return m
}

// RegexMatcher regex 判斷
type RegexMatcher struct {
	regexps []*regexp.Regexp
}

// Matches 判斷傳入 tag 是否匹配
func (m *RegexMatcher) Matches(tag string) bool {
	for _, r := range m.regexps {
		if r.MatchString(tag) {
			return true
		}
	}
	return false
}

// NewSemVerMatcher 建立 SemVerMatcher 物件
func NewSemVerMatcher(versions []string) *SemVerMatcher {
	m := &SemVerMatcher{}
	for _, version := range versions {
		m.versions = append(m.versions, semver.Must(semver.NewVersion(version)))
	}
	return m
}

// SemVerMatcher 以 Semantic Versioning 2.0.0 判斷
type SemVerMatcher struct {
	versions []*semver.Version
}

// Matches 判斷傳入 tag 是否匹配
func (m *SemVerMatcher) Matches(tag string) bool {
	for _, v := range m.versions {
		other := semver.Must(semver.NewVersion(strings.TrimPrefix(tag, "v")))
		if v.Equal(*other) {
			return true
		}
	}
	return false
}
