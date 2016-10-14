# wtf
A tool that integrates github commit and PR comment threads into your dev loop via an easy to use CLI.

## wtf is wtf
I love github, and I love talking about code with my team in github projects. However, I have noticed that those conversations
seem to get lost in the sands of time. It would be nice to be able to easily see historical comment threads about a project or 
even a specific file. It would be even nicer to leave comments about code or a specific github PR in your terminal as you are
diving through the code. That is the goal of wtf.

## Usage

If you do not set an environmental variable called `GITHUB_ACCESS_TOKEN` to yout github access api access token [See here](https://github.com/blog/1509-personal-api-tokens), then 
you will be making unauthenticated Github API requests. This means you will only be able to make 60 requests per hour.

### List comment threads for a file
To view a list of all of the github threads for a file in your git repo then type:
`wtf ls $filepath` 

For example:
`wtf ls src/main/scala/main.scala`

![Example output](http://i.imgur.com/yx5mVff.png)

Here is a [link](https://github.com/Samangan/owLint/commit/f8ee99c1c0e2c7c9b5ea9635657f2410d2b24e38) to the actual comment thread in github as an example.

## Build

`go build`

`go install`

