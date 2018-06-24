package main

import (
	"flag"

	"clanpj/lisao/rest/github"
	"clanpj/lisao/tasks"
)

var lisaoRepo = flag.String("repo_name", "lisao-bot", "Lisao bot repo name.")
var lisaoOwner = flag.String("repo_owner", "Bubblyworld", "Lisao bot repo owner.")

type State struct {
	githubClient *github.Client

	buildsPool *tasks.Pool
}

func NewState() *State {
	return &State{
		githubClient: github.NewClient(*lisaoOwner, *lisaoRepo),

		buildsPool: tasks.NewPool("builds", 1, tasks.DoBuild),
	}
}
