package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path"
)

func getPostsPathAndGenPath(basePath string) (postsPath string, genPath string) {
	return path.Join(basePath, "posts"), path.Join(basePath, "gen")
}

func checkStructureAndReturnListOfPosts(basePath string, postsPath string, genPath string) (files []os.FileInfo, err error) {
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

func processPosts(files []os.FileInfo, postsPath string, genPath string) (err error) {
	for _, v := range files {
		content, err := ioutil.ReadFile(path.Join(postsPath, v.Name()))
		if err != nil {
			return fmt.Errorf("Failed to read file %s: `%v'", v.Name(), err)
		}

		converted := blackfriday.MarkdownBasic(content)

		err = ioutil.WriteFile(
			path.Join(genPath, v.Name()),
			converted,
			// Read-only.
			0444)
		if err != nil {
			return fmt.Errorf("Failed to write file %s: `%v'", v.Name(), err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("You must pass a path to the site. No less, no more - just that\n")
		return
	}
	basePath := os.Args[1]
	postsPath, genPath := getPostsPathAndGenPath(basePath)
	files, err := checkStructureAndReturnListOfPosts(basePath, postsPath, genPath)
	if err != nil {
		fmt.Printf("Failure: `%v'\n", err)
		return
	}

	err = processPosts(files, postsPath, genPath)
	if err != nil {
		fmt.Printf("Failure: `%v'\n", err)
		return
	}
}
