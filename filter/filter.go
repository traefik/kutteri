package filter

import (
	"fmt"
	"log"
)

// Build a filter.
func Build(fns ...func() string) string {
	var query string
	for _, fn := range fns {
		query += " " + fn()
	}
	log.Println(query)
	return query
}

// Issue type:issue.
func Issue() string {
	return "type:issue"
}

// PullRequest type:pr.
func PullRequest() string {
	return "type:pr"
}

// Open state:open.
func Open() string {
	return "state:open"
}

// InTitle in:title.
func InTitle() string {
	return "in:title"
}

// Content simple words.
func Content(data string) func() string {
	return func() string {
		return data
	}
}

// CreatedAfter created:>xxx.
func CreatedAfter(date string) func() string {
	return func() string {
		return fmt.Sprintf("created:>%s", date)
	}
}

// UpdatedAfter updated:>xxx.
func UpdatedAfter(date string) func() string {
	return func() string {
		return fmt.Sprintf("updated:>%s", date)
	}
}

// MergedAfter merged:>xxx.
func MergedAfter(date string) func() string {
	return func() string {
		return fmt.Sprintf("merged:>%s", date)
	}
}

// ClosedAfter closed:>xxx.
func ClosedAfter(date string) func() string {
	return func() string {
		return fmt.Sprintf("closed:>%s", date)
	}
}

// Repo repo:xxx/yyy.
func Repo(owner, repoName string) func() string {
	return func() string {
		return fmt.Sprintf("repo:%s/%s", owner, repoName)
	}
}
