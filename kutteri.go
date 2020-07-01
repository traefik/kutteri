package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/containous/flaeg"
	"github.com/containous/kutteri/core"
	"github.com/containous/kutteri/locker"
	"github.com/ogier/pflag"
)

// Configuration Command configuration.
type Configuration struct {
	GitHubToken    string `long:"ghtoken" description:"GitHub Token."`
	SlackToken     string `long:"sltoken" description:"Slack Token. [required]"`
	Owner          string `short:"o" description:"Repository owner. [required]"`
	RepositoryName string `long:"repo-name" short:"r" description:"Repository name. [required]"`
	SearchFilter   string `long:"search-filter" short:"f" description:"Search filter. [required]"`
	BotName        string `long:"bot-name" description:"Name of the bot in Slack."`
	BotIcon        string `long:"bot-icon" description:"Icon of the bot in Slack."`
	Channel        string `short:"c" description:"Slack channel. [required]"`
	DryRun         bool   `long:"dry-run" description:"Dry run mode."`
	ServerMode     bool   `long:"server" description:"Server mode."`
	ServerPort     int    `long:"port" description:"Server port."`
}

// NoOption empty struct.
type NoOption struct{}

func main() {
	config := &Configuration{
		BotName:    "Bender",
		BotIcon:    ":bento:",
		ServerPort: 80,
		DryRun:     true,
	}

	rootCmd := &flaeg.Command{
		Name:                  "kuterri",
		Description:           `Chalepoxenus Kutteri: Track a GitHub repository and publish on Slack.`,
		DefaultPointersConfig: &Configuration{},
		Config:                config,
		Run:                   runCommand(config),
	}

	flag := flaeg.New(rootCmd, os.Args[1:])

	// version
	versionCmd := &flaeg.Command{
		Name:                  "version",
		Description:           "Display the version.",
		Config:                &NoOption{},
		DefaultPointersConfig: &NoOption{},
		Run: func() error {
			displayVersion()
			return nil
		},
	}

	flag.AddCommand(versionCmd)
	err := flag.Run()
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		log.Fatalf("Error: %v\n", err)
	}
}

func runCommand(config *Configuration) func() error {
	return func() error {
		if len(config.GitHubToken) == 0 {
			config.GitHubToken = os.Getenv("GITHUB_TOKEN")
		}

		fmt.Println("config.SlackToken", config.SlackToken)
		if len(config.SlackToken) == 0 {
			config.SlackToken = os.Getenv("SLACK_TOKEN")
		}

		err := validateConfig(config)
		if err != nil {
			return err
		}

		return launch(config)
	}
}

func launch(config *Configuration) error {
	if config.ServerMode {
		server := &server{config: config}
		return server.ListenAndServe()
	}
	return execute(*config)
}

func execute(config Configuration) error {
	slackConfig := core.SlackConfig{
		BotName: config.BotName,
		BotIcon: config.BotIcon,
		Channel: config.Channel,
		DryRun:  config.DryRun,
	}

	rec := locker.TimeLocker{
		FilePath: "./kutteri.lock",
		HourBack: 1,
	}

	date, err := rec.GetLastTime()
	if err != nil {
		return err
	}

	criterion := core.SearchCriterion{
		Owner:    config.Owner,
		RepoName: config.RepositoryName,
		Words:    config.SearchFilter,
		Date:     date,
	}

	ctx := context.Background()

	bot := core.NewBot(ctx, config.GitHubToken, config.SlackToken, slackConfig)

	err = bot.ProcessAll(ctx, criterion)
	if err != nil {
		return err
	}

	_, err = rec.SaveLastTime()
	return err
}

func validateConfig(config *Configuration) error {
	err := required(config.Owner, "owner")
	if err != nil {
		return err
	}

	err = required(config.RepositoryName, "repo-name")
	if err != nil {
		return err
	}

	return required(config.SlackToken, "sltoken")
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		return fmt.Errorf("%s is mandatory", fieldName)
	}
	return nil
}

type server struct {
	config *Configuration
}

func (s *server) ListenAndServe() error {
	return http.ListenAndServe(":"+strconv.Itoa(s.config.ServerPort), s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid http method: %s", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := execute(*s.config)
	if err != nil {
		log.Printf("Report error: %v", err)
		http.Error(w, "Report error.", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Chalepoxenus Kutteri: Scheluded.\n")
}
