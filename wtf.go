package main

import (
	"fmt"
	"github.com/Samangan/wtf/thread"
	"github.com/mgutz/ansi"
	"github.com/urfave/cli"
	"gopkg.in/kyokomi/emoji.v1"
	"log"
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
			Usage:   "List all threads for `filename`",
			Action: func(c *cli.Context) error {
				// TODO: allow filename to be an absolute path as well as a relative path
				filename := c.Args().First()
				threads, err := thread.GetThreads(filename)

				if err != nil {
					log.Fatal(err)
				}

				// TODO: Put below in a helper function to print out the threads to the console.
				// * I should make a utility class maybe for the formatting and splicing stuff.
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
