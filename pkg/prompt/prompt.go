package prompt

import (
	"errors"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"strings"
)

// Confirm 確認問題
func Confirm(log *logrus.Logger, c interface{}) (ok bool, err error) {
	b, err := yaml.Marshal(c)
	if err != nil {
		log.Printf("%#v", c)
	} else {
		log.Println(string(b))
	}
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
