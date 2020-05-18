package plugin

import (
	diff "github.com/r3labs/diff"
	gh "github.com/squidarth/kingfig/github"
)

type FigObject struct {
	ResourceType     string
	GithubRepository *gh.Repository
}

type Plugin interface {
	GetExistingFromRemote() map[string]interface{}
	Update() bool
}

func (figObject FigObject) ApplyConfig() error {
	if figObject.GithubRepository != nil {
		return figObject.GithubRepository.ApplyConfig()
	}
	return nil
}

func (figObject FigObject) GetDiff() []diff.Change {
	if figObject.GithubRepository != nil {
		return figObject.GithubRepository.GetDiff()
	}
	return nil
}

func (e *FigObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params struct {
		ResourceType string `yaml:"resource_type"`
	}

	if err := unmarshal(&params); err != nil {
		return err
	}

	var repository gh.Repository
	if err := unmarshal(&repository); err != nil {

	}

	e.ResourceType = params.ResourceType
	e.GithubRepository = &repository
	return nil
}
