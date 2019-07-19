package main

import (
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/softleader/depl/pkg/deployer"
	"github.com/softleader/depl/pkg/prompt"
	"strings"
)

func prereleaseQuestions(c *prereleaseCmd) error {
	if err := prompt.AskRequired("Name of image to build", c.Image.Name, &c.Image.Name); err != nil {
		return err
	}

	if err := prompt.AskYesNo("Force to delete the tag if it already exists?", "y", &c.Force); err != nil {
		return err
	}

	services, err := deployer.FilterServiceByApp(logrus.StandardLogger(), "depl", metadata.String(), c.Deployer, c.Image.Name)
	if len(services) == 0 || err != nil {
		if err := prompt.Ask("Docker service id to update image (leave blank if you don't need to update)", c.DockerServiceID, &c.DockerServiceID); err != nil {
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
		c.DockerServiceID = services[i].ID
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

	if err := prompt.AskRequired("Designating development stage to build, e.g. 0 for alpha, 1 for beta, 2 for release candidate", c.Stage, &c.Stage); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of the owner (user or org) of the repo to create tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of repo to create to create tag", c.SourceRepo, &c.SourceRepo); err != nil {
		return err
	}

	if !c.SkipTests {
		if err := prompt.AskRequired("Config server host:port to run the test", c.ConfigServer, &c.ConfigServer); err != nil {
			return err
		}

		if err := prompt.Ask("Label of config server to run the test, e.g. sqlServer", c.ConfigLabel, &c.ConfigLabel); err != nil {
			return err
		}
	}

	if err := prompt.AskRequired("Username of docker registry for jib to build", c.Auth.Username, &c.Auth.Username); err != nil {
		return err
	}

	if err := prompt.AskRequired("Password of docker registry for jib to build", c.Auth.Password, &c.Auth.Password); err != nil {
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
	return prereleaseQuestions(c)
}
