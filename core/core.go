package core

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/containous/kutteri/filter"
	"github.com/google/go-github/v27/github"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
)

// SlackConfig Slack configuration.
type SlackConfig struct {
	BotName string
	BotIcon string
	Channel string
	DryRun  bool
}

// SearchCriterion GitHub search criterion
type SearchCriterion struct {
	Owner    string
	RepoName string
	Words    string
	Date     string
}

type messageModel struct {
	Issues []github.Issue
	Title  string
}

// Bot The Bot
type Bot struct {
	ghClient    *github.Client
	slackClient *slack.Client
	SlackConfig SlackConfig
}

// NewBot Create a new Bot
func NewBot(ctx context.Context, ghToken string, slackToken string, slackConfig SlackConfig) *Bot {
	return &Bot{
		ghClient:    NewGitHubClient(ctx, ghToken),
		slackClient: slack.New(slackToken),
		SlackConfig: slackConfig,
	}
}

// ProcessAll Execute bot on issues and pull request
func (b *Bot) ProcessAll(ctx context.Context, searchConfig SearchCriterion) error {
	err := b.ProcessIssues(ctx, searchConfig)
	if err != nil {
		return err
	}

	return b.ProcessPullRequest(ctx, searchConfig)
}

// ProcessIssues Execute bot on issues
func (b *Bot) ProcessIssues(ctx context.Context, searchConfig SearchCriterion) error {
	var models []messageModel

	// issue:new
	queryIssueNew := filter.Build(
		filter.Repo(searchConfig.Owner, searchConfig.RepoName),
		filter.Issue, filter.Open,
		filter.CreatedAfter(searchConfig.Date),
		filter.InTitle, filter.Content("traefik"))

	issuesNew, err := search(ctx, b.ghClient, queryIssueNew)
	if err != nil {
		return err
	}

	if len(issuesNew) > 0 {
		models = append(models, messageModel{Title: ":new: New Issues:", Issues: issuesNew})
	}

	// issue:updated
	queryIssueUpdated := filter.Build(
		filter.Repo(searchConfig.Owner, searchConfig.RepoName),
		filter.Issue, filter.Open,
		filter.UpdatedAfter(searchConfig.Date),
		filter.InTitle, filter.Content("traefik"))

	issuesUpdated, err := searchUpdated(ctx, b.ghClient, queryIssueUpdated, issuesNew)
	if err != nil {
		return err
	}

	if len(issuesUpdated) > 0 {
		models = append(models, messageModel{Title: ":up: Updated Issues:", Issues: issuesUpdated})
	}

	for _, model := range models {
		err = sendMessageToSlack(b.slackClient, b.SlackConfig, model)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProcessPullRequest Execute bot on pull request
func (b *Bot) ProcessPullRequest(ctx context.Context, searchConfig SearchCriterion) error {
	var models []messageModel

	// pr:new
	queryPRNew := filter.Build(
		filter.Repo(searchConfig.Owner, searchConfig.RepoName),
		filter.PullRequest, filter.Open,
		filter.CreatedAfter(searchConfig.Date),
		filter.InTitle, filter.Content("traefik"))

	prNew, err := search(ctx, b.ghClient, queryPRNew)
	if err != nil {
		return err
	}

	if len(prNew) > 0 {
		models = append(models, messageModel{Title: ":new: New Pull Request:", Issues: prNew})
	}

	// pr:updated
	queryPRUpdated := filter.Build(
		filter.Repo(searchConfig.Owner, searchConfig.RepoName),
		filter.PullRequest, filter.Open,
		filter.UpdatedAfter(searchConfig.Date),
		filter.InTitle, filter.Content("traefik"))

	prUpdated, err := searchUpdated(ctx, b.ghClient, queryPRUpdated, prNew)
	if err != nil {
		return err
	}

	if len(prUpdated) > 0 {
		models = append(models, messageModel{Title: ":up: Updated Pull Request:", Issues: prUpdated})
	}

	// pr:merged
	queryPRMerged := filter.Build(
		filter.Repo(searchConfig.Owner, searchConfig.RepoName),
		filter.PullRequest,
		filter.MergedAfter(searchConfig.Date),
		filter.InTitle, filter.Content("traefik"))

	prMerged, err := search(ctx, b.ghClient, queryPRMerged)
	if err != nil {
		return err
	}

	if len(prMerged) > 0 {
		models = append(models, messageModel{Title: ":facepunch: Merged Pull Request:", Issues: prMerged})
	}

	for _, model := range models {
		err = sendMessageToSlack(b.slackClient, b.SlackConfig, model)
		if err != nil {
			return err
		}
	}

	return nil
}

func searchUpdated(ctx context.Context, client *github.Client, query string, newIssues []github.Issue) ([]github.Issue, error) {
	issues, err := search(ctx, client, query)
	if err != nil {
		return nil, err
	}

	var realIssues []github.Issue
	for _, issue := range issues {
		var freshIssue bool
		for _, newIssue := range newIssues {
			if newIssue.GetNumber() == issue.GetNumber() {
				freshIssue = true
				break
			}
		}
		if !freshIssue {
			realIssues = append(realIssues, issue)
		}
	}

	return realIssues, nil
}

func search(ctx context.Context, client *github.Client, query string) ([]github.Issue, error) {
	options := &github.SearchOptions{
		Sort:        "updated",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 25},
	}

	result, _, err := client.Search.Issues(ctx, query, options)
	if err != nil {
		return nil, err
	}

	return result.Issues, nil
}

func sendMessageToSlack(client *slack.Client, config SlackConfig, model messageModel) error {
	text, err := createMessageContent(model)
	if err != nil {
		return err
	}

	if config.DryRun {
		fmt.Println(text)
		return nil
	}

	return sendToSlack(client, config, text)
}

func createMessageContent(model messageModel) (string, error) {
	tmpl := template.New("slack message")
	_, err := tmpl.Parse(`
*{{ .Title }}*
{{range $issue := .Issues -}}
:round_pushpin: <{{ $issue.HTMLURL }}|{{ $issue.Number }}>: {{ $issue.Title }}
{{end -}}`)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = tmpl.Execute(b, model)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func sendToSlack(client *slack.Client, config SlackConfig, text string) error {
	params := slack.PostMessageParameters{
		Username:  config.BotName,
		AsUser:    false,
		IconEmoji: config.BotIcon,
		Markdown:  true,
	}

	_, _, err := client.PostMessage(config.Channel, text, params)
	return err
}

// NewGitHubClient create a new GitHub client
func NewGitHubClient(ctx context.Context, token string) *github.Client {
	if len(token) == 0 {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
