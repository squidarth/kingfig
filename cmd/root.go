package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/r3labs/diff"
	"github.com/spf13/cobra"
	"github.com/squidarth/kingfig/auth"
	gh "github.com/squidarth/kingfig/github"
	plug "github.com/squidarth/kingfig/plugin"
	"gopkg.in/yaml.v3"
)

func getAuthSettings() (*auth.AuthSettings, error) {
	var authFileBytes, err = ioutil.ReadFile(AuthorizationDetailsFile)

	if err != nil {

		return nil, err
	}

	var authSettings auth.AuthSettings
	err = yaml.Unmarshal(authFileBytes, &authSettings)

	if err != nil {
		fmt.Println("Error reading authorization file")

		return nil, err
	}
	return &authSettings, nil
}

func getConfigurationFiles(rootDir string) []string {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func readFigObjectsFromFile(yamlFile string) map[string]plug.FigObject {

	var fileBytes, _ = ioutil.ReadFile(yamlFile)
	var objectMap map[string]plug.FigObject
	if err := yaml.Unmarshal(fileBytes, &objectMap); err != nil {
		fmt.Println(err.Error())
	}

	return objectMap
}

func displayableChangelog(resourceName string, changelog []diff.Change) string {
	var finalString string
	finalString = ""

	for _, change := range changelog {

		switch change.Type {
		case "update":
			finalString += "\033[1;32m" + "+" + strings.Join(change.Path, ".") + ": " + fmt.Sprintf("%v", change.To) + "\n" + "\033[1;0m"

			finalString += "\033[1;31m" + "-" + strings.Join(change.Path, ".") + ": " + fmt.Sprintf("%v", change.From)

		case "delete":
			finalString += "-" + strings.Join(change.Path, ".") + ": " + fmt.Sprintf("%v", change.From)
		case "create":
			finalString += "+" + strings.Join(change.Path, ".") + ": " + fmt.Sprintf("%v", change.To)
		}
		finalString += "\n"
	}

	return finalString
}

var (
	// Used for flags.
	OutputFilePath           string
	ResourceName             string
	ConfigurationDir         string
	NoDryRun                 bool
	AuthorizationDetailsFile string
	rootCmd                  = &cobra.Command{
		Use:   "kingfig",
		Short: "CLI for applying setting changes.",
		Long:  `KingFig is a program for declaratively specifying your settings.`,
	}

	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Generates a new kingfig yaml file from an existing remote resource",
		Long:  "Generates a new kingfig yaml file from an existing remote resource",
		Run: func(cmd *cobra.Command, args []string) {
			resourceType := args[0]

			var authSettings, _ = getAuthSettings()
			if resourceType == "GithubRepository" {
				fmt.Println("Generating new Github Repository config...")

				repoOwner := args[1]
				repoName := args[2]

				var repo = gh.GetRepoFromRemote(repoOwner, repoName, *authSettings)

				var newResourceMap = make(map[string]gh.Repository)
				newResourceMap[ResourceName] = repo

				var bytes, _ = yaml.Marshal(newResourceMap)

				ioutil.WriteFile(OutputFilePath, bytes, 0644)

				fmt.Println("New config written to: " + OutputFilePath)
			} else {
				fmt.Print("Resource type unknown")
			}
		},
	}

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: `Apply your local changes to the remote server.`,
		Long:  `Apply your local changes to the remote server.`,
		Run: func(cmd *cobra.Command, args []string) {
			var files = getConfigurationFiles(ConfigurationDir)

			var authSettings, _ = getAuthSettings()

			var fullConfiguration = make(map[string]plug.FigObject)

			for _, file := range files {
				for k, v := range readFigObjectsFromFile(file) {
					fullConfiguration[k] = v
				}
			}

			if NoDryRun {
				fmt.Println("Applied the following changes:")
			} else {
				fmt.Println("Will apply the following changes:")
			}
			for resourceName, config := range fullConfiguration {
				var changeLogDisplay = displayableChangelog(resourceName, config.GetDiff(*authSettings))

				if changeLogDisplay != "" {
					fmt.Println(resourceName + ":")
					fmt.Println(changeLogDisplay)
				}
			}

			if NoDryRun {
				for _, config := range fullConfiguration {
					err := config.ApplyConfig(*authSettings)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		},
	}
)

func init() {
	var homeDir, _ = os.UserHomeDir()

	applyCmd.Flags().BoolVarP(&NoDryRun, "no-dry-run", "d", false, "Applies the configuration in production")

	applyCmd.Flags().StringVarP(&ConfigurationDir, "configuration-dir", "c", ".", "Directory to look for config files")

	newCmd.Flags().StringVarP(&OutputFilePath, "output-file-path", "o", "", "File where you'd like the new configuration to go.")
	newCmd.MarkFlagRequired("output-file-path")

	newCmd.Flags().StringVarP(&ResourceName, "resource-name", "r", "", "Key under which resource should be stored.")

	newCmd.MarkFlagRequired("resource-name")

	rootCmd.PersistentFlags().StringVarP(&AuthorizationDetailsFile, "authorization", "a", fmt.Sprintf("%s/.kingfig/auth.yaml", homeDir), "Location of authorization details file.")

	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(newCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
