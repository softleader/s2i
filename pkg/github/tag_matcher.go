package github

import (
	"fmt"
	"github.com/blang/semver"
	"regexp"
	"strings"
)

// TagMatcherStrategy 用來方便判斷是哪個 matcher
type TagMatcherStrategy struct {
	Regex  bool `yaml:"regex"`
	SemVer bool `yaml:"semver"`
}

// TagMatcher 判斷是否有 match tags
type TagMatcher interface {
	// Matches 傳入的 string 是否匹配
	Matches(s string) bool
}

// NewRegexMatcher 建立 RegexMatcher 物件
func NewRegexMatcher(exprs []string) (*RegexMatcher, error) {
	m := &RegexMatcher{}
	for _, expr := range exprs {
		rr, err := regexp.Compile(expr)
		if err != nil {
			return nil, fmt.Errorf("requires a valid regexp expression: %s", err)
		}
		m.regexps = append(m.regexps, rr)
	}
	return m, nil
}

// RegexMatcher regex 判斷
type RegexMatcher struct {
	regexps []*regexp.Regexp
}

// Matches 判斷傳入 tag 是否匹配
func (m *RegexMatcher) Matches(s string) bool {
	for _, r := range m.regexps {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

// NewSemVerMatcher 建立 SemVerMatcher 物件
func NewSemVerMatcher(ranges []string) (*SemVerMatcher, error) {
	m := &SemVerMatcher{}
	for _, r := range ranges {
		rr, err := semver.ParseRange(r)
		if err != nil {
			return nil, fmt.Errorf("requires a valid semver2 ranges: %s", err)
		}
		if m.r == nil {
			m.r = rr
		} else {
			m.r = m.r.OR(rr)
		}
	}
	return m, nil
}

// SemVerMatcher 以 Semantic Versioning 2.0.0 判斷
type SemVerMatcher struct {
	r semver.Range
}

// Matches 判斷傳入 tag 是否匹配
func (m *SemVerMatcher) Matches(s string) bool {
	v, err := semver.Parse(strings.TrimPrefix(s, "v"))
	if err != nil {
		return false
	}
	return m.r(v)
}
