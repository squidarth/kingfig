package plugin

import (
	diff "github.com/r3labs/diff"
	"github.com/squidarth/kingfig/auth"
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

func (figObject FigObject) ApplyConfig(authSettings auth.AuthSettings) error {
	if figObject.GithubRepository != nil {
		return figObject.GithubRepository.ApplyConfig(authSettings)
	}
	return nil
}

func (figObject FigObject) GetDiff(authSettings auth.AuthSettings) []diff.Change {
	if figObject.GithubRepository != nil {
		return figObject.GithubRepository.GetDiff(authSettings)
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
