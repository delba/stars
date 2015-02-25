package main

import (
	"fmt"
	"os"

	"github.com/delba/stars/github"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	username := os.Args[1]

	repositories, err := github.GetFollowingStarred(username)
	handle(err)

	fmt.Println(repositories)
}
