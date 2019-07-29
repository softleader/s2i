package main

import (
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/deployer"
	"github.com/softleader/s2i/pkg/prompt"
	"strings"
)

func releaseQuestions(c *releaseCmd) error {
	if err := prompt.AskRequired("Name of image to build", c.Image.Name, &c.Image.Name); err != nil {
		return err
	}

	if err := prompt.AskRequired("Tag of image to build", c.Image.Tag, &c.Image.Tag); err != nil {
		return err
	}

	services, err := deployer.FilterServiceByApp(logrus.StandardLogger(), "s2i", metadata.String(), c.Deployer, c.Image.Name)
	if len(services) == 0 || err != nil {
		if err := prompt.Ask("Service id to update image (leave blank if you don't need to update)", c.ServiceID, &c.ServiceID); err != nil {
			return err
		}
	} else {
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
		c.ServiceID = services[i].ID
	}

	ok, err := prompt.Confirm(logrus.StandardLogger(), c)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	if err := prompt.AskRequired("Name of branch to create to create tag", c.SourceBranch, &c.SourceBranch); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of the owner (user or org) of the repo to create tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of repo to create to create tag", c.SourceRepo, &c.SourceRepo); err != nil {
		return err
	}

	if err := prompt.AskRequired("Jenkins to run the pipeline", c.Jenkins, &c.Jenkins); err != nil {
		return err
	}

	ok, err = prompt.Confirm(logrus.StandardLogger(), c)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	logrus.Println("That's try again!")
	return releaseQuestions(c)
}
