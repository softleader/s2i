package main

import (
	"github.com/kr/pretty"
	"github.com/manifoldco/promptui"
)

func (c *preReleaseCmd) prompt() error {
	if err := ask("Name of image to build", c.image, &c.image); err != nil {
		return err
	}

	sp := promptui.Select{
		Label: "Force to delete the tag (" + c.tag + ") if it already exists?",
		Items: []string{"yes", " no"},
	}
	_, selected, err := sp.Run()
	if err != nil {
		return err
	}
	c.force = selected == "yes"

	if err := ask("Docker service id to update image (leave blank if you don't need to update)", c.dockerServiceID, &c.dockerServiceID); err != nil {
		return err
	}

	sp = promptui.Select{
		Label: "Do you want to go through all of the questions?",
		Items: []string{"no", "yes"},
	}
	_, selected, err = sp.Run()
	if err != nil {
		return err
	}
	if selected == "no" {
		pretty.Print(c)

		sp = promptui.Select{
			Label: "Is this OK?",
			Items: []string{"yes", "no"},
		}
		_, selected, err = sp.Run()
		if err != nil {
			return err
		}
		if selected == "yes" {
			return nil
		}
	}

	if err := ask("Name of branch to create to create tag", c.sourceBranch, &c.sourceBranch); err != nil {
		return err
	}

	if err := ask("Designating development stage to build, e.g. 0 for alpha, 1 for beta, 2 for release candidate", c.stage, &c.stage); err != nil {
		return err
	}

	if err := ask("Name of the owner (user or org) of the repo to create tag", c.sourceOwner, &c.sourceOwner); err != nil {
		return err
	}

	if err := ask("Name of repo to create to create tag", c.sourceRepo, &c.sourceRepo); err != nil {
		return err
	}

	if !c.skipTests {
		if err := ask("Config server host:port to run the test", c.configServer, &c.configServer); err != nil {
			return err
		}

		if err := ask("Label of config server to run the test, e.g. sqlServer", c.configLabel, &c.configLabel); err != nil {
			return err
		}
	}

	if err := ask("Deployer host:port to update stack service", c.deployer, &c.deployer); err != nil {
		return err
	}

	if err := ask("Username of docker registry for jib to build", c.auth.Username, &c.auth.Username); err != nil {
		return err
	}

	if err := ask("Password of docker registry for jib to build", c.auth.Password, &c.auth.Password); err != nil {
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
