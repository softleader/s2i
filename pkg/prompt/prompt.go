package prompt

import (
	"errors"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/deployer"
	"gopkg.in/yaml.v2"
	"strings"
)

const (
	width = 70
)

// Confirm 確認問題
func Confirm(log *logrus.Logger, c interface{}) (ok bool, err error) {
	log.Printf("%s", strings.Repeat("-", width))
	b, err := yaml.Marshal(c)
	if err != nil {
		log.Printf("%#v", c)
	} else {
		log.Println(string(b))
	}
	log.Printf("%s", strings.Repeat("-", width))
	err = AskYesNo("Is this OK?", "y", &ok)
	return
}

// Ask 問單一問題
func Ask(question, defaultValue string, ref *string) (err error) {
	p := promptui.Prompt{
		Label:   question,
		Default: defaultValue,
	}
	*ref, err = p.Run()
	return
}

// AskRequired 問單一問題, 且必填
func AskRequired(question, defaultValue string, ref *string) (err error) {
	p := promptui.Prompt{
		Label:   question,
		Default: defaultValue,
		Validate: func(s string) error {
			if strings.TrimSpace(s) == "" {
				return errors.New("required")
			}
			return nil
		},
	}
	*ref, err = p.Run()
	return
}

// AskArrayRequired 問單一問題, 但可以回答多個, 必填
func AskArrayRequired(question string, defaultValue []string, ref *[]string, sep string) (err error) {
	p := promptui.Prompt{
		Label:   question,
		Default: strings.Join(defaultValue, sep),
		Validate: func(s string) error {
			if strings.TrimSpace(s) == "" {
				return errors.New("required")
			}
			return nil
		},
	}
	ans, err := p.Run()
	if err != nil {
		return err
	}
	*ref = strings.Split(ans, sep)
	return
}

// AskYesNoBool 問 boolean 問題, 並且以 bool 做為 defaultValue
func AskYesNoBool(question string, defaultValue bool, ref *bool) (err error) {
	if defaultValue {
		return AskYesNo(question, "y", ref)
	}
	return AskYesNo(question, "n", ref)
}

// AskYesNo 問 boolean 問題
func AskYesNo(question, defaultValue string, ref *bool) (err error) {
	var yesNo string
	p := promptui.Prompt{
		Label:   question + " (y/n)",
		Default: defaultValue,
		Validate: func(s string) error {
			ans := strings.ToLower(s)
			if ans == "y" || ans == "yes" || ans == "n" || ans == "no" {
				return nil
			}
			return errors.New("please answer yes(y) or no(n)")
		},
	}
	yesNo, err = p.Run()
	*ref = strings.ToLower(yesNo) == "y" || strings.ToLower(yesNo) == "yes"
	return
}

// AskServiceID 問 docker swarm serviceID 問題
func AskServiceID(log *logrus.Logger, agent, agentVersion, deployerURL, app, defaultValue string, ref *string) (err error) {
	question := "Service id to update image (leave blank if you don't need to update)"

	if defaultValue != "" {
		return Ask(question, defaultValue, ref)
	}

	services, err := deployer.FilterServiceByApp(log, agent, agentVersion, deployerURL, app)
	if len(services) == 0 || err != nil { // 在 deployer 上找不到任何已部署的服務, 或連線發生問題
		return Ask(question, defaultValue, ref)
	}

	services = append([]deployer.DockerService{{
		Name: "I don't need to update",
	}}, services...)
	prompt := promptui.Select{
		Label: "Select docker service to update",
		Items: services,
		Templates: &promptui.SelectTemplates{
			Active:   promptui.IconSelect + " {{ .Name }}\t{{ .Ports }}",
			Inactive: "  {{ .Name }}\t{{ .Ports }}",
			Selected: promptui.IconGood + " {{ .Name }}\t{{ .Ports }}",
		},
		Searcher: func(input string, index int) bool {
			channel := services[index]
			name := strings.Replace(strings.ToLower(channel.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}
	i, _, err := prompt.Run()
	if err != nil {
		return err
	}
	*ref = services[i].ID
	return nil
}
