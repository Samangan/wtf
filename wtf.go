package main

import (
	"fmt"
	"github.com/Samangan/wtf/thread"
	"github.com/mgutz/ansi"
	"github.com/urfave/cli"
	"gopkg.in/kyokomi/emoji.v1"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()

	app.Name = "wtf"
	app.Usage = "Github comment threads in your terminal"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"l"},
			Usage:   "List threads",
			Subcommands: []cli.Command{
				{
					Name:    "all",
					Aliases: []string{"a"},
					Usage:   "List all threads in project.",
					Action: func(c *cli.Context) error {
						// TODO: Implement
						return nil
					},
				},
				{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "List all threads for `fileName`",
					Action: func(c *cli.Context) error {
						// TODO: allow filename to be an absolute path as well as a relative path
						filename := c.Args().First()

						// TODO:
						// This wont work...
						// If there are threads in different commits then I need to list all of the
						// commits and PRs that have threads and then ask the user to select which one to see like so:
						// FILE           DATE          THREADS        TYPE
						// blah.sh    12/12/2016          10          commit
						// blah.sh    12/12/2015          5        pull request
						//    ....
						threads, _ := thread.GetThreads(filename)

						for _, thread := range threads {
							fmt.Println("\n@@ " + thread.File + "  @@")

							first, second, _ := spliceDiff(thread.Diff, thread.Comments[0].Pos)
							fmt.Println(first)

							lime := ansi.ColorCode("green+h:black")
							reset := ansi.ColorCode("reset")

							emoji.Println(lime, ":email: "+thread.Comments[0].Author+" [ "+thread.Comments[0].Date.String()+"] : "+thread.Comments[0].Body+"", reset)

							for _, comment := range thread.Comments[1:] {
								emoji.Println(lime, "\\__ :email: "+comment.Author+" [ "+comment.Date.String()+"] : "+comment.Body+"", reset)
							}
							fmt.Println(second)
						}

						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}

// spliceDiff takes in a diff and then returns
// the first half of the diff (ending with the location of the comment thread)
// and also the second half of the diff (the rest of the file that comes after the comment thread)
func spliceDiff(diff string, pos int) (string, string, error) {
	s := strings.SplitN(diff, "\n", -1)
	return strings.Join(s[:pos+1], "\n"), strings.Join(s[pos+1:], "\n"), nil
}
