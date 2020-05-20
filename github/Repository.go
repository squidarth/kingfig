package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	diff "github.com/r3labs/diff"

	githubv4 "github.com/shurcooL/githubv4"
	auth "github.com/squidarth/kingfig/auth"
	"golang.org/x/oauth2"
)

type Repository struct {
	Description        string `yaml:"description" diff:"description"`
	Owner              string `yaml:"owner" diff: "owner"`
	HasIssuesEnabled   bool   `yaml:"has_issues_enabled" diff:"has_issues_enabled"`
	HasProjectsEnabled bool   `yaml:"has_projects_enabled" diff:"has_projects_enabled"`
	HasWikiEnabled     bool   `yaml:"has_wiki_enabled" diff:"has_wiki_enabled"`
	HomepageUrl        string `yaml:"homepage_url" diff:"homepage_url"`
	Name               string `yaml:"name" diff:"name"`
	Id                 string `yaml:"id" diff:"id"`
	Template           bool   `yaml:"template" diff:"template"`
}

func GetRepoFromRemote(ownerName string, name string, authSettings auth.AuthSettings) Repository {
	client := getGHClient(authSettings)

	var q struct {
		Repository struct {
			Description        string
			HasIssuesEnabled   bool
			HasProjectsEnabled bool
			HasWikiEnabled     bool
			Name               string
			HomepageUrl        string
			IsTemplate         bool
			Owner              struct {
				Login string
			}
			Id string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(ownerName),
		"name":  githubv4.String(name),
	}

	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		fmt.Println(err)
		println("ERROR")
	}

	return Repository{
		Description:        q.Repository.Description,
		HasIssuesEnabled:   q.Repository.HasProjectsEnabled,
		HasProjectsEnabled: q.Repository.HasProjectsEnabled,
		HasWikiEnabled:     q.Repository.HasWikiEnabled,
		Name:               q.Repository.Name,
		HomepageUrl:        q.Repository.HomepageUrl,
		Template:           q.Repository.IsTemplate,
		Owner:              q.Repository.Owner.Login,
		Id:                 q.Repository.Id,
	}
}

func getGHClient(authSettings auth.AuthSettings) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authSettings.GithubApiToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	return client
}

func (r Repository) ApplyConfig(authSettings auth.AuthSettings) error {
	var client = getGHClient(authSettings)
	var m struct {
		UpdateRepository struct {
			Repository struct {
				Id string
			}
		} `graphql:"updateRepository(input: $input)"`
	}

	var description = githubv4.String(r.Description)

	var name = githubv4.String(r.Name)
	var hasIssuesEnabled = githubv4.Boolean(r.HasIssuesEnabled)

	var template = githubv4.Boolean(r.Template)

	var hasWikiEnabled = githubv4.Boolean(r.HasProjectsEnabled)
	var hasProjectsEnabled = githubv4.Boolean(r.HasIssuesEnabled)

	homepageURI, err := url.Parse(r.HomepageUrl)

	if err != nil {
		fmt.Println(err)
		return errors.New("Invalid homepage URL")
	}
	var homepageURL = githubv4.URI{URL: homepageURI}

	input := githubv4.UpdateRepositoryInput{
		Description:        &description,
		HasIssuesEnabled:   &hasIssuesEnabled,
		HasWikiEnabled:     &hasWikiEnabled,
		HomepageURL:        &homepageURL,
		HasProjectsEnabled: &hasProjectsEnabled,
		Name:               &name,
		Template:           &template,
		RepositoryID:       r.Id,
	}
	return client.Mutate(context.Background(), &m, input, nil)
}

func (r Repository) GetDiff(authSettings auth.AuthSettings) []diff.Change {
	var remoteRepo = GetRepoFromRemote(r.Owner, r.Name, authSettings)

	var changelog, _ = diff.Diff(remoteRepo, r)

	return changelog
}
