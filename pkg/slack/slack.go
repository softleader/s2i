package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/docker"
	"github.com/softleader/s2i/pkg/github"
	"github.com/softleader/s2i/pkg/release"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const (
	// apiFile 儲存 context 的檔案名稱
	apiFile = "api.yaml"
	// envMountVolume key to specify root-dir to store apiFile
	envMountVolume = "SL_PLUGIN_MOUNT"
)

var (
	// ErrMountVolumeNotExist 代表沒有發現 root-dir
	ErrMountVolumeNotExist = errors.New(`mount volume not found
It looks like you are running the command outside slctl (https://github.com/softleader/slctl)
Please set SL_PLUGIN_MOUNT variable to manually specify the location for the command to store data 
For more details: https://github.com/softleader/slctl/wiki/Plugins-Guide#mount-volume`)
	// ErrMissingSlackWebhookURL 代表沒有傳過 --slack-webhook-url
	ErrMissingSlackWebhookURL = errors.New(`missing slack webhook URL
You must specify '--slack-webhook-url' at the very first time using slack webhook
Or use '--skip-slack' to skip hooking slack`)
)

// API collect slack api information
type API struct {
	log        *logrus.Logger
	path       string
	WebhookURL string
}

// Post post webhook to url with text
func Post(log *logrus.Logger, metadata *release.Metadata, release *github.Release, url string, image *docker.SoftleaderHubImage) error {
	api, err := loadAPI(log)
	if err != nil {
		return err
	}
	if len(url) > 0 {
		api.WebhookURL = url
		defer func(api *API) {
			if err := api.save(); err != nil {
				log.Debug(err)
			}
		}(api)
	}
	if len(api.WebhookURL) <= 0 {
		return ErrMissingSlackWebhookURL
	}
	payload := &slack.WebhookMessage{
		Text: fmt.Sprintf("SIT %s@%s 過版", image.Name, image.Tag),
	}
	if release != nil {
		attachment := newAttachment(release, metadata)
		log.Debugf("appending attachment: %+v", attachment)
		payload.Attachments = append(payload.Attachments, attachment)
	}
	fmt.Printf("%+v", payload)
	return slack.PostWebhook(api.WebhookURL, payload)
}

func newAttachment(release *github.Release, metadata *release.Metadata) slack.Attachment {
	return slack.Attachment{
		Title:      release.TagName,
		AuthorName: release.Author.GetLogin(),
		AuthorLink: release.Author.GetHTMLURL(),
		AuthorIcon: release.Author.GetAvatarURL(),
		TitleLink:  release.HTMLURL,
		Footer:     fmt.Sprintf("s2i@%v", metadata.String()),
		Ts:         json.Number(strconv.FormatInt(release.PublishedAt.Unix(), 10)),
	}
}

func loadAPI(log *logrus.Logger) (*API, error) {
	mount, found := os.LookupEnv(envMountVolume)
	if !found {
		return nil, ErrMountVolumeNotExist
	}
	mount, err := homedir.Expand(mount)
	if err != nil {
		return nil, err
	}
	api := &API{
		log:  log,
		path: filepath.Join(mount, apiFile),
	}
	return api, api.load()
}

func (api *API) load() error {
	api.log.Debugf("loading slack api from: %s\n", api.path)
	data, err := ioutil.ReadFile(api.path)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if os.IsNotExist(err) {
		return nil
	}
	return json.Unmarshal(data, api)
}

func (api *API) save() error {
	data, err := json.Marshal(api)
	if err != nil {
		return err
	}
	api.log.Debugf("refreshing slack api information to: %s\n", string(data))
	return ioutil.WriteFile(api.path, data, 0644)
}
