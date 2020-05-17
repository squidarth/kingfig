package github

import (
	"context"
	"fmt"
	"os"

	structs "github.com/fatih/structs"
	diff "github.com/r3labs/diff"
	githubv4 "github.com/shurcooL/githubv4"
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

func GetRepoFromRemote(ownerName string, name string) Repository {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
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
	}
}

func (r Repository) GetDiff() []diff.Change {
	var remoteRepo = GetRepoFromRemote(r.Owner, r.Name)

	var changelog, _ = diff.Diff(remoteRepo, r)

	return changelog
}

func (r Repository) GetExistingFromRemote() map[string]interface{} {
	/* Sample Repository, until
	 * we have API connection set up
	 */

	var newRepo = Repository{
		Description:        "A repo",
		HasIssuesEnabled:   true,
		HasProjectsEnabled: true,
		HasWikiEnabled:     true,
		Name:               "repo-repo",
		HomepageUrl:        "www.google.com",
		Id:                 "123456",
		Template:           false,
	}

	return structs.New(newRepo).Map()
}

func (r Repository) Update() bool {
	return true
}
