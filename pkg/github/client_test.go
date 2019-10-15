package github

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestFindRemoteFromTokenClone(t *testing.T) {
	config := `[core]
	repositoryformatversion = 0
	filemode = false
	bare = false
	logallrefupdates = true
	symlinks = false
	ignorecase = true
[remote "origin"]
	url = https://ec5365ad1a31edd35446b04738aee99dfbf8a7d4@github.com/softleader-product/softleader-pos-policy-rpc.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "develop"]
	remote = origin
	merge = refs/heads/develop
[branch "v1"]
	remote = origin
	merge = refs/heads/v1
[branch "v2_5"]
	remote = origin
	merge = refs/heads/v2_5
[gui]
	wmstate = zoomed
	geometry = 1199x645+38+38 273 293`

	owner, repo := findRemoteOriginURL(logrus.StandardLogger(), config)
	if owner != "softleader-product" {
		t.Fatalf("owner should be softleader-product, but got %q", owner)
	}
	if repo != "softleader-pos-policy-rpc" {
		t.Fatalf("repo should be softleader-product, but got %q", owner)
	}
}

func TestFindRemoteFromSshClone(t *testing.T) {
	config := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
	ignorecase = true
	precomposeunicode = true
[remote "origin"]
	url = git@github.com:softleader/softleader-jasmine.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "develop"]
	remote = origin
	merge = refs/heads/develop
[branch "Robert"]
	remote = origin
	merge = refs/heads/Robert
[branch "v2_5"]
	remote = origin
	merge = refs/heads/v2_5`

	owner, repo := findRemoteOriginURL(logrus.StandardLogger(), config)
	if owner != "softleader" {
		t.Fatalf("owner should be softleader, but got %q", owner)
	}
	if repo != "softleader-jasmine" {
		t.Fatalf("repo should be softleader-jasmine, but got %q", owner)
	}
}

func TestFindRemoteFromHttpsClone(t *testing.T) {
	config := `[core]
    repositoryformatversion = 0
    filemode = false
    bare = false
    logallrefupdates = true
    symlinks = false
    ignorecase = true
[remote "origin"]
    url = https://github.com/softleader/softleader-jasmine.git
    fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
    remote = origin
    merge = refs/heads/master
[branch "develop"]
    remote = origin
    merge = refs/heads/develop
[branch "barch-to-rpc"]
    remote = origin
    merge = refs/heads/barch-to-rpc
[branch "v2_5"]
    remote = origin
    merge = refs/heads/v2_5`

	owner, repo := findRemoteOriginURL(logrus.StandardLogger(), config)
	if owner != "softleader" {
		t.Fatalf("owner should be softleader, but got %q", owner)
	}
	if repo != "softleader-jasmine" {
		t.Fatalf("repo should be softleader-jasmine, but got %q", owner)
	}
}
