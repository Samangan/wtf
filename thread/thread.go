package thread

import (
	"fmt"
	"github.com/google/go-github/github"
	"time"
	//	"os"
)

// TODO: better name
type Object struct {
	// File is the path of the file
	File string

	// Type is either `pull_request` or `commit`
	Type string

	// Threads is the list of threads in this file
	Threads []*Thread
}

type Thread struct {
	// Comments that are contained in this thread
	Comments []*Comment

	// File is the path of the comments for this thread or nil if general comment
	File string

	// Diff is the relevant subset of the file that the thread is about
	Diff string

	// Url is the url of the thread
	Url string

	// Pos is the position in the diff of the thread
	Pos int
}

type Comment struct {
	// Author is the name of the comment author
	Author string

	// Date is the date of the comment
	Date time.Time

	// Body is the actual text of the comment
	Body string

	// Url is the github url to the comment
	Url string

	// CommitId is the git commit SHA
	CommitId string

	// Pos is the position in the diff of the thread
	Pos int
}

func GetThreads(fileName string) ([]*Thread, error) {
	owner, repo, err := parseRepoName()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Collate commit and PR comments for this file:
	threadMap, err := getAllThreadsForRepo(owner, repo)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Return only those threads that are for this specific filename
	fileThreads := threadMap[fileName]

	threads := []*Thread{}

	for _, comments := range fileThreads {
		diff, _ := getDiffFromCommit(owner, repo, comments[0].CommitId, fileName)
		url := comments[0].Url
		pos := comments[0].Pos

		threads = append(threads, &Thread{
			Comments: comments,
			File:     fileName,
			Diff:     diff,
			Url:      url,
			Pos:      pos,
		})
	}

	return threads, nil
}

func parseRepoName() (string, string, error) {
	// TODO: Parse the repo name from `git remote -v`
	return "Samangan", "owLint", nil
}

func getDiffFromCommit(owner string, repo string, commitId string, filename string) (string, error) {
	// TODO: DRY up this github client shit:
	client := github.NewClient(nil) // TODO: Implement authentication (see github readme)

	commit, _, err := client.Repositories.GetCommit(owner, repo, commitId)

	if err != nil {
		return "", err
	}

	for _, file := range commit.Files {
		if *file.Filename == filename {
			return *file.Patch, nil
		}
	}

	return "", err
}

func getAllThreadsForRepo(owner string, repo string) (map[string]map[int][]*Comment, error) {
	// Get all commit comments for this repo.
	client := github.NewClient(nil) // TODO: Implement authentication (see github readme)

	opt := &github.ListOptions{PerPage: 50} // TODO: using pagination grab all of the results from the github api. Not just 50.

	comments, _, err := client.Repositories.ListComments("Samangan", "owLint", opt)

	if err != nil {
		return nil, err
	}

	// TODO: Get all PR comments for this repo.

	// Create 'threads' from the list of all comments.
	// A 'thread' is a list of comments all on the same `position` of the diff.
	return combineCommentsIntoThreads(comments)
}

func combineCommentsIntoThreads(comments []*github.RepositoryComment) (map[string]map[int][]*Comment, error) {
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

	return threadMap, nil
}

// TODO: func GetAllThreads()
