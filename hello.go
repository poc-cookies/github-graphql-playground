package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var query struct {
	Viewer struct {
		Login     githubv4.String
		CreatedAt githubv4.DateTime
	}
}

// fetchRepoDescription fetches description of repo with owner and name.
func fetchRepoDescription(client *githubv4.Client, ctx context.Context, owner, name string) (string, error) {
	var q struct {
		Repository struct {
			Description string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	err := client.Query(ctx, &q, variables)
	return q.Repository.Description, err
}

func fetchIssueId(client *githubv4.Client, ctx context.Context, issueNumber int, owner, name string) (string, error) {
	var q struct {
		Repository struct {
			Issue struct {
				Id string
			} `graphql:"issue(number: $issueNumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"name":        githubv4.String(name),
		"issueNumber": githubv4.Int(issueNumber),
	}

	err := client.Query(ctx, &q, variables)
	return q.Repository.Issue.Id, err
}

func main() {
	fmt.Println("Hello, world.")

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	// For GitHub Enterprise client use the following code:
	// client := githubv4.NewEnterpriseClient("https://ghe.telenordigital.com/api/graphql", httpClient)

	// Query viewer
	err1 := client.Query(context.Background(), &query, nil)
	if err1 != nil {
		fmt.Println("ERROR 1:", err1)
	}
	fmt.Println("    Login:", query.Viewer.Login)
	fmt.Println("CreatedAt:", query.Viewer.CreatedAt)

	// Query repo
	repoOwner := "poc-cookies"
	repoName := "friendlyhello"
	repoDescription, err2 := fetchRepoDescription(client, context.Background(), repoOwner, repoName)
	if err2 != nil {
		fmt.Println("ERROR 2:", err2)
	}
	fmt.Println(repoDescription)

	// Issue ID
	issueNumber := 1
	issueId, err3 := fetchIssueId(client, context.Background(), issueNumber, repoOwner, repoName)
	if err3 != nil {
		fmt.Println("ERROR 3:", err3)
	}
	fmt.Println(issueId)

	fmt.Println("The END.")
}
