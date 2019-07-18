package deployer

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// UpdateService 更新 deployer 的 service
func UpdateService(log *logrus.Logger, agent, agentVersion, deployer, dockerServiceID, image, tag string) error {
	log.Printf("updating docker service id: %s", dockerServiceID)
	params := make(map[string]string)
	params["image"] = fmt.Sprintf("hub.softleader.com.tw/%s:%s", image, tag)
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

// FilterServiceByImage 依照 image 查詢 docker service
func FilterServiceByImage(log *logrus.Logger, agent, agentVersion, deployer, image string) ([]DockerService, error) {
	resty.SetDebug(log.IsLevelEnabled(logrus.DebugLevel))
	params := make(map[string]string)
	params["label"] = fmt.Sprintf("com.docker.stack.image=%s", image)
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
