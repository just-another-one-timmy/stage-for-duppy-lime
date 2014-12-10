package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func checkStructureAndReturnListOfPosts(basePath string) (files []os.FileInfo, err error) {
	postsPath := path.Join(basePath, "posts")

	// Note: . and .. are not included.
	allPostFileInfos, err := ioutil.ReadDir(postsPath)
	if err != nil {
		return nil, fmt.Errorf("Can't read directory %s: `%v'\n", postsPath, err)
	}

	for _, fileInfo := range allPostFileInfos {
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("Found dir %s in %s. There must be only regular files.\n", fileInfo.Name(), postsPath)
		}
	}

	if len(allPostFileInfos) == 0 {
		return nil, fmt.Errorf("No posts found!\n")
	}

	genPath := path.Join(basePath, "gen")
	err = os.RemoveAll(genPath)
	// Note: if the path doesn't exist, err will be nil.
	if err != nil {
		return nil, fmt.Errorf("Can't remove path %s: `%v'\n", genPath, err)
	}

	err = os.Mkdir(genPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("Can't create path %s: `%v'\n", genPath, err)
	}

	return allPostFileInfos, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("You must pass a path to the site. No less, no more - just that\n")
		return
	}

	_, err := checkStructureAndReturnListOfPosts(os.Args[1])
	if err != nil {
		fmt.Printf("Failure: %v\n", err)
		return
	}
	fmt.Printf("Didn't encounter any problems, but not going to generate the site either. Come back later when this thing is updated to do something useful.\n")
}
