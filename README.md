# Kingfig

**Warning: This is a proof of concept. Expect very little to work and everything to change.**

Kingfig is a tool to manage your and your team's settings
on different web services via configuration files. Discover a change to your Github repository settings, and wish you could know who on your team made the change, and what the context was?
With Kingfig, since you make changes to these settings with
config files, you can use Version Control to track changes, and
don't have to depend on the products you're using have "history"
features.

In its current form, it only supports updating Github repositories, but the idea is that it will support a plugin system
that makes supporting other web services really easy to do.

![kingfig-demo](https://user-images.githubusercontent.com/850115/82394624-871d6f80-9a17-11ea-9ac1-982cccfb68cc.png)

This was inspired by tools like [Terraform](https://www.terraform.io/) and [Ansible](https://www.ansible.com/). 

## Installation

For now, this is only available via source. 
Please follow the [instructions for installing Go](https://golang.org/doc/install). Once Go is installed, you can run:

```bash
$ git clone git@github.com:squidarth/kingfig.git 
$ cd kingfig/
$ go install .
$ kingfig --help
```

## Example Usage

As mentioned before, `kingfig` currently only supports 
updating Github repositories.

Steps to use this:

1. Make a Github repo to store your `kingfig` configs
2. Write a YAML file to represent your Github repository that looks something like this:

```yaml
# A file called # config.yaml
sids_repo:
  resource_type: "GithubRepository"
  name: "rubocop-assist"
  owner: "squidarth"
  description:  "Webapp to help write rubocop rules"
  homepage_url: "https://github.com/squidarth/rubocop-assist"
  has_projects_enabled: true
  has_wiki_enabled: true
  has_issues_enabled: true
  template: true
  id: "MDEwOlJlcG9zaXRvcnkxMDA4OTcwNjA="
```
3. Generate a Github API Token
4. Write a YAML file that contains your Github authorization information at the path ~/.kingfig/auth.yaml

```yaml
# ~/.kingfig/auth.yaml
github_api_token: $YOUR_GITHUB_API_TOKEN
```

5. In your config directory, run `kingfig apply`

This will print out a list of changes that will be applied.

6. Run `kingfig apply --no-dry-run` to actually apply the changes to your Github repo

You can find the ID of your Github repository using the [Github GraphQL Explorer](https://developer.github.com/v4/explorer/), and a query like:

```graphql
query{
  repository(name:"rubocop-assist", owner:"squidarth"){
    id
  }
}
```

## Current Integrations

* Github Repositories 

## Coming up soon

* A `kingfig new` command that generates a yaml configuration from
an existing web service, to make bootstrapping your `kingfig` configurations easier.
* Have some automated tool for writing the authorization file (there are too many steps to getting set up right now)

## How you can help

The next step for getting this to an alpha stage would be building integrations with tools that people would love this for. Go ahead and make a Github issue is there's a tool you would love to have managed with a config file!

As I continue to build this, I would also love feedback on making
a plugin API that's convenient and nice to use, and requires little
knowledge of how this tool is built. I

Also, help me come up with a better name!!

