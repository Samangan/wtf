package thread

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Thread struct {
	Comments []*Comment // Comments that are contained in this thread
	File     string     // File is the path of the comments for this thread or nil if general comment
	Diff     string     // Diff is the relevant subset of the file that the thread is about
	Url      string     // Url is the url of the thread
	Pos      int        // Pos is the position in the diff of the thread
}

type Comment struct {
	Author   string    // Author is the name of the comment author
	Date     time.Time // Date is the date of the comment
	Body     string    // Body is the actual text of the comment
	Url      string    // Url is the github url to the comment
	CommitId string    // CommitId is the git commit SHA
	Pos      int       // Pos is the position in the diff of the thread
}

// GetThreads returns a list of all of the PR and Commit comment threads for the given filepath.
func GetThreads(fileName string) ([]*Thread, error) {
	owner, repo := parseRepoName()

	// Collate commit and PR comments for this file:
	threadMap, err := getAllThreadsForRepo(owner, repo)

	if err != nil {
		return nil, err
	}

	// Return only those threads that are for this specific filename:
	fileThreads := threadMap[fileName]

	threads := []*Thread{}

	for _, comments := range fileThreads {
		diff, err := getDiffFromCommit(owner, repo, comments[0].CommitId, fileName)

		if err != nil {
			return nil, err
		}

		threads = append(threads, &Thread{
			Comments: comments,
			File:     fileName,
			Diff:     diff,
			Url:      comments[0].Url,
			Pos:      comments[0].Pos,
		})
	}

	return threads, nil
}

// parseRepoName returns the owner and repository name for the
// repository that is currently in this working directory.
// It uses `git remove -v` to find this information.
func parseRepoName() (string, string) {
	cmd := exec.Command("git", "remote", "-v")
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	outStr := string(out[:])

	// TODO: Below is hacky... Clean up:
	parsedGitOutput := strings.Split(strings.Split(outStr, ":")[1], "/")
	return parsedGitOutput[0], strings.Split(parsedGitOutput[1], ".git")[0]
}

// getGithubAuthenticatedClient retrieves the github access token from the `GITHUB_ACCESS_TOKEN` environmental
// variable and then returns the authenticated client.
// If no env variable exists then it will return nil and the unauthenticated client will still be used.
func getGithubAuthenticatedClient() *http.Client {
	accessToken, tokenExists := os.LookupEnv("GITHUB_ACCESS_TOKEN")

	if !tokenExists {
		fmt.Println("`GITHUB_ACCESS_TOKEN` env variable not found. Using un-authenticated github API!")
		return nil
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(oauth2.NoContext, ts)
}

// getDiffFromCommit returns the github diff for this filename and commitId hash
func getDiffFromCommit(owner string, repo string, commitId string, filename string) (string, error) {
	// TODO: DRY UP the github client usage: (I should only have to init the github client and grab the env var once per command)
	client := github.NewClient(getGithubAuthenticatedClient())
	commit, _, err := client.Repositories.GetCommit(owner, repo, commitId)

	if err != nil {
		return "", err
	}

	for _, file := range commit.Files {
		if *file.Filename == filename {
			return *file.Patch, nil
		}
	}

	return "", errors.New("Diff was not found for this commitId")
}

// getAllThreadsForRepo returns all comment threads for this repository.
func getAllThreadsForRepo(owner string, repo string) (map[string]map[int][]*Comment, error) {
	// Get all commit comments for this repo.
	client := github.NewClient(getGithubAuthenticatedClient())
	opt := &github.ListOptions{PerPage: 50} // TODO: using pagination to grab all of the results from the github api. Not just 50.

	comments, _, err := client.Repositories.ListComments(owner, repo, opt)

	if err != nil {
		return nil, err
	}

	// TODO: Get all PR comments as well.

	// Create 'threads' from the list of all comments.
	// A 'thread' is a list of comments all on the same `position` of the diff.
	return combineCommentsIntoThreads(comments), nil
}

// combineCommentsIntoThreads returns a map of `filename` -> `linenumber` -> List of comments from an array of github comments.
func combineCommentsIntoThreads(comments []*github.RepositoryComment) map[string]map[int][]*Comment {
	// Combine comments that are in the same `path` and same `position`
	threadMap := make(map[string]map[int][]*Comment)

	for _, comment := range comments {
		if comment.Path != nil {
			if _, ok := threadMap[*comment.Path]; !ok {
				threadMap[*comment.Path] = make(map[int][]*Comment)
			}

			threadMap[*comment.Path][*comment.Position] = append(threadMap[*comment.Path][*comment.Position], &Comment{
				Author:   *comment.User.Login,
				Date:     *comment.CreatedAt,
				Body:     *comment.Body,
				CommitId: *comment.CommitID,
				Url:      *comment.URL,
				Pos:      *comment.Position,
			})
		}
	}

	return threadMap
}
