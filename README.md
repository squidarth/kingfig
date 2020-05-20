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

1. Make a folder to store your `kingfig` configs
2. Generate a Github API Token
3. Write a YAML file that contains your Github authorization information at the path ~/.kingfig/auth.yaml

```yaml
# ~/.kingfig/auth.yaml
github_api_token: $YOUR_GITHUB_API_TOKEN
```

4. For an existing Github repo you'd like to manage the settings for, run:

```
$ kingfig new -o sids_repo.yaml --resource-name sids_repo GithubRepository squidarth rubocop-assist
```

In this case, my repo is `squidarth/rubocop-assist`.

5. This will generate a file that looks like this:

```yaml
# sids_repo.yaml
sids_repo:
    description: Webapp to help write rubocop rules well
    owner: squidarth
    has_issues_enabled: true
    has_projects_enabled: true
    has_wiki_enabled: true
    homepage_url: https://github.com/squidarth/rubocop-assist
    name: rubocop-assist
    id: MDEwOlJlcG9zaXRvcnkxMDA4OTcwNjA=
    template: true
```

Go ahead and modify, say, the `description` field in this file.

6. Run `kingfig apply --no-dry-run` to actually apply the changes to your Github repo

## Current Integrations

* Github Repositories 

## Coming up soon

* Have some automated tool for writing the authorization file (there are too many steps to getting set up right now)

## How you can help

The next step for getting this to an alpha stage would be building integrations with tools that people would love this for. Go ahead and make a Github issue is there's a tool you would love to have managed with a config file!

As I continue to build this, I would also love feedback on making
a plugin API that's convenient and nice to use, and requires little
knowledge of how this tool is built. I

Also, help me come up with a better name!!

