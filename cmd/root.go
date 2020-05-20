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
	plug "github.com/squidarth/kingfig/plugin"
	"gopkg.in/yaml.v3"
)

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
	ConfigurationDir         string
	NoDryRun                 bool
	AuthorizationDetailsFile string
	rootCmd                  = &cobra.Command{
		Use:   "kingfig",
		Short: "CLI for applying setting changes.",
		Long:  `KingFig is a program for declaratively specifying your settings.`,
	}

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: `Apply your local changes to the remote server.`,
		Long:  `Apply your local changes to the remote server.`,
		Run: func(cmd *cobra.Command, args []string) {
			var files = getConfigurationFiles(ConfigurationDir)
			var authFileBytes, err = ioutil.ReadFile(AuthorizationDetailsFile)

			if err != nil {
				fmt.Println(err)
				fmt.Println("Error opening authorization file")
				return
			}

			var authSettings auth.AuthSettings
			err = yaml.Unmarshal(authFileBytes, &authSettings)

			if err != nil {
				fmt.Println("Error reading authorization file")

				return
			}

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
				var changeLogDisplay = displayableChangelog(resourceName, config.GetDiff(authSettings))

				if changeLogDisplay != "" {
					fmt.Println(resourceName + ":")
					fmt.Println(changeLogDisplay)
				}
			}

			if NoDryRun {
				for _, config := range fullConfiguration {
					err := config.ApplyConfig(authSettings)
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
	rootCmd.PersistentFlags().StringVarP(&AuthorizationDetailsFile, "authorization", "a", fmt.Sprintf("%s/.kingfig/auth.yaml", homeDir), "Location of authorization details file.")

	rootCmd.AddCommand(applyCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
