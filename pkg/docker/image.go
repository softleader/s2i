package docker

import (
	"fmt"
	"github.com/blang/semver"
	"strings"
)

const (
	softleaderHub = "hub.softleader.com.tw"
)

// SoftleaderHubImage 表示該 image 會放在 hub.softleader.com.tw
type SoftleaderHubImage struct {
	Name, Tag string
}

// SetPreRelease 設定 tag 的 pre-release 版號
func (i *SoftleaderHubImage) SetPreRelease(preRelease string) {
	version := strings.TrimPrefix(i.Tag, "v")
	sv, err := semver.Parse(version)
	if err != nil {
		return
	}
	sv.Pre = nil
	prv, err := semver.NewPRVersion(preRelease)
	if err != nil {
		return
	}
	sv.Pre = append(sv.Pre, prv)
	pr := sv.String()
	if strings.HasPrefix(i.Tag, "v") {
		pr = "v" + pr
	}
	i.Tag = pr
}

// String 返回適用於 hub.softleader.com.tw 的 image 全名
func (i *SoftleaderHubImage) String() string {
	return fmt.Sprintf("%s/%s:%s", softleaderHub, i.Name, i.Tag)
}

// CheckValid 檢查 image 資訊是否有效
func (i *SoftleaderHubImage) CheckValid() error {
	if strings.TrimSpace(i.Name) == "" {
		return fmt.Errorf("image name is required")
	}
	if strings.TrimSpace(i.Tag) == "" {
		return fmt.Errorf("tag is required")
	}
	// GitHub 建議我們用 v 開頭, 但 v 開頭不符合 semver, 所以檢查時固定拿掉
	v := strings.TrimPrefix(i.Tag, "v")
	_, err := semver.Parse(v)
	if err != nil {
		return fmt.Errorf("requires valid semver2 tag: %s", err)
	}
	return nil
}
