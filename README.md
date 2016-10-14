# wtf
A simple little tool that integrates github commit and PR comment threads into your dev loop via an easy to use CLI.

## Usage

If you do not set an environmental variable called `GITHUB_ACCESS_TOKEN` to your github api access token [See here](https://github.com/blog/1509-personal-api-tokens), then 
you will be making unauthenticated Github API requests. This means you will only be able to make 60 requests per hour, and you will only be able to access comments in public repositories.

### List comment threads for a file
To view a list of all of the github threads for a file in your git repo then type:
`wtf ls $filepath` 

For example:
`wtf ls src/main/scala/main.scala`
w
![Example output](http://i.imgur.com/yx5mVff.png)

Here is a [link](https://github.com/Samangan/owLint/commit/f8ee99c1c0e2c7c9b5ea9635657f2410d2b24e38) to the actual comment thread in github as an example.

## Development

`go build`

`go install`

`go test ./...`

## TODO
 - Allow for relative filepath as well (SEE TODO)
 - If user inserts a filepath that does not exist return error message / code
 - Parse repo and owner name from `git remove -v` (SEE TODO)
 - Implement pagination (SEE TODO)
 - Improve the UI and usability more.
