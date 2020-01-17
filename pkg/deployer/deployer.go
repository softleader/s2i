package deployer

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/docker"
	"github.com/softleader/s2i/pkg/github"
	"gopkg.in/resty.v1"
	"strconv"
)

// SlackHook 代表要 hook slack 的內容
type SlackHook struct {
	Enabled bool
	Release *github.Release
}

// UpdateService 更新 deployer 的 service
func UpdateService(log *logrus.Logger, agent, agentVersion, deployer, dockerServiceID string, image *docker.SoftleaderHubImage, hook SlackHook) error {
	log.Printf("Updating docker service id: %s", dockerServiceID)
	params := make(map[string]string)
	params["image"] = image.String()
	if hook.Enabled {
		params["slack"] = "1"
		params["title"] = hook.Release.TagName
		params["title_link"] = hook.Release.HTMLURL
		params["author_name"] = hook.Release.Author.GetLogin()
		params["author_link"] = hook.Release.Author.GetHTMLURL()
		params["author_icon"] = hook.Release.Author.GetAvatarURL()
		params["ts"] = strconv.FormatInt(hook.Release.PublishedAt.Unix(), 10)
	}
	resty.SetDebug(log.IsLevelEnabled(logrus.DebugLevel))
	_, err := resty.R().
		SetQueryParams(params).
		SetHeader("User-Agent", fmt.Sprintf("%s/%s", agent, agentVersion)).
		Get(fmt.Sprintf("%s/services/update/%s", deployer, dockerServiceID))
	return err
}

// DockerService 包含了 docker service 的資訊
type DockerService struct {
	ID       string
	Image    string
	Mode     string
	Name     string
	Ports    string
	Replicas string
}

// FilterServiceByApp 依照 label=app 查詢 docker service
func FilterServiceByApp(log *logrus.Logger, agent, agentVersion, deployer, app string) ([]DockerService, error) {
	resty.SetDebug(log.IsLevelEnabled(logrus.DebugLevel))
	params := make(map[string]string)
	params["label"] = fmt.Sprintf("app=%s", app)
	return FilterService(log, agent, agentVersion, deployer, params)
}

// FilterService 依照條件查詢 service
func FilterService(log *logrus.Logger, agent, agentVersion, deployer string, params map[string]string) ([]DockerService, error) {
	resp, err := resty.R().
		SetQueryParams(params).
		SetHeader("User-Agent", fmt.Sprintf("%s/%s", agent, agentVersion)).
		Get(fmt.Sprintf("%s/services/filter", deployer))
	if err != nil {
		return nil, err
	}
	var services []DockerService
	if err = json.Unmarshal(resp.Body(), &services); err != nil {
		return nil, err
	}
	return services, nil
}
