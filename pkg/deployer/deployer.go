package deployer

import (
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
