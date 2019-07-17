package main

import (
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"strings"
)

func (c *releaseCmd) prompt() error {
	var yesNo string

	if err := ask("Name of image to build", c.Image, &c.Image); err != nil {
		return err
	}

	if err := ask("Do you want to go through all of the questions?", "n", &yesNo); err != nil {
		return err
	}
	if strings.ToLower(yesNo) == "n" || strings.ToLower(yesNo) == "no" {
		b, err := yaml.Marshal(c)
		if err != nil {
			logrus.Printf("%#v", c)
		} else {
			logrus.Println(string(b))
		}
		if err := ask("Is this OK?", "y", &yesNo); err != nil {
			return err
		}
		if strings.ToLower(yesNo) == "y" || strings.ToLower(yesNo) == "yes" {
			return nil
		}
	}

	if err := ask("Name of branch to create to create tag", c.SourceBranch, &c.SourceBranch); err != nil {
		return err
	}

	if err := ask("Name of the owner (user or org) of the repo to create tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := ask("Name of repo to create to create tag", c.SourceRepo, &c.SourceRepo); err != nil {
		return err
	}

	if err := ask("Jenkins to run the pipeline", c.Jenkins, &c.Jenkins); err != nil {
		return err
	}

	return nil
}

func ask(question, defaultValue string, ref *string) (err error) {
	p := promptui.Prompt{
		Label:   question,
		Default: defaultValue,
	}
	*ref, err = p.Run()
	return
}
