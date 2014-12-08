package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("You must pass a path to the site. No less, no more - just that\n")
		return
	}

	basePath := os.Args[1]
	postsPath := path.Join(basePath, "posts")
	allPostFileInfos, err := ioutil.ReadDir(postsPath)
	if err != nil {
		fmt.Printf("Can't read directory %s\n", postsPath)
		return
	}

	if len(allPostFileInfos) == 0 {
		fmt.Printf("No posts found!\n")
		return
	}

	fmt.Printf("Didn't encounter any problems, but not going to generate the site either. Come back later when this thing is updated to do somethign useful.\n")
}
