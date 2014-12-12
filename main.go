package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

type Post struct {
	Title, Name, Content string

	Tags []string

	Date time.Time
}

type Index struct {
	Posts []*Post
}

func makePostFromContent(name string, rawContent []byte) (*Post, error) {
	lines := bytes.SplitN(rawContent, []byte("\n"), 5)
	if len(lines) < 5 {
		return nil, errors.New("Incorrect header format. Some field is missing.")
	}

	post := new(Post)
	post.Name = name
	post.Title = string(lines[0][:])

	date, err := time.Parse("2006-Jan-02", string(lines[1][:]))
	if err != nil {
		return nil, err
	}
	post.Date = date

	post.Tags = strings.Split(string(lines[2][:]), " ")

	if len(string(lines[3][:])) != 0 {
		return nil, errors.New("Header must be separated from content with an empty line.")
	}

	post.Content = string(blackfriday.MarkdownBasic(lines[4][:])[:])
	return post, nil
}

func getSitePaths(basePath string) (postsPath string, genPath string, assetsPath string, tempaltesPath string) {
	return path.Join(basePath, "posts"), path.Join(basePath, "gen"), path.Join(basePath, "assets"), path.Join(basePath, "templates")
}

func checkStructureAndReturnListOfPosts(basePath string, postsPath string, genPath string) ([]os.FileInfo, error) {
	// Note: . and .. are not included.
	allPostFileInfos, err := ioutil.ReadDir(postsPath)
	if err != nil {
		return nil, fmt.Errorf("Can't read directory %s: `%v'", postsPath, err)
	}

	for _, fileInfo := range allPostFileInfos {
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("Found dir %s in %s. There must be only regular files.", fileInfo.Name(), postsPath)
		}
	}

	if len(allPostFileInfos) == 0 {
		return nil, fmt.Errorf("No posts found!")
	}

	err = os.RemoveAll(genPath)
	// Note: if the path doesn't exist, err will be nil.
	if err != nil {
		return nil, fmt.Errorf("Can't remove path %s: `%v'", genPath, err)
	}

	err = os.Mkdir(genPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("Can't create path %s: `%v'", genPath, err)
	}

	return allPostFileInfos, nil
}

func makePostsSlice(files []os.FileInfo, postsPath string) ([]*Post, error) {
	posts := make([]*Post, 0)
	for _, v := range files {
		content, err := ioutil.ReadFile(path.Join(postsPath, v.Name()))
		if err != nil {
			return nil, fmt.Errorf("Failed to read file %s: `%v'", v.Name(), err)
		}

		post, err := makePostFromContent(v.Name(), content)
		if err != nil {
			return nil, fmt.Errorf("Failed to create post %s: `%v'", v.Name(), err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func createPostsFiles(posts []*Post, genPath string, templatesPath string) error {
	postTemplate := template.New("post")
	postTemplate, err := postTemplate.ParseFiles(path.Join(templatesPath, "post"))
	if err != nil {
		return err
	}

	indexTemplate := template.New("index")
	indexTemplate, err = indexTemplate.ParseFiles(path.Join(templatesPath, "index"))
	if err != nil {
		return err
	}
	var indexContent bytes.Buffer
	err = indexTemplate.Execute(&indexContent, Index{posts})

	for _, post := range posts {
		if post.Name == "index.html" {
			return errors.New("You must not have a post with name index.html.")
		}
		var postContent bytes.Buffer
		err = postTemplate.Execute(&postContent, post)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(
			path.Join(genPath, post.Name),
			postContent.Bytes(),
			// Read-only.
			0444)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(
		path.Join(genPath, "index.html"),
		indexContent.Bytes(),
		// Read-only.
		0444)
	if err != nil {
		return err
	}

	return nil
}

func copyAssets(assetsPath string, genPath string) error {
	// Note: . and .. are not included.
	allPostFileInfos, err := ioutil.ReadDir(assetsPath)
	if err != nil {
		return fmt.Errorf("Can't read directory %s: `%v'", assetsPath, err)
	}

	for _, fileInfo := range allPostFileInfos {
		if fileInfo.IsDir() {
			return fmt.Errorf("Found dir %s in %s. There must be only regular files.", fileInfo.Name(), assetsPath)
		}

		pathToGeneratedFile := path.Join(genPath, fileInfo.Name())

		// Based on
		// https://www.socketloop.com/tutorials/golang-copy-file
		//
		err := func(src string, dest string) error {
			r, err := os.Open(src)
			if err != nil {
				return err
			}
			defer r.Close()
			w, err := os.Create(dest)
			if err != nil {
				return err
			}
			defer w.Close()
			_, err = io.Copy(w, r)
			return err
		}(path.Join(assetsPath, fileInfo.Name()),
			pathToGeneratedFile)

		if err != nil {
			return fmt.Errorf("Failed to copy asset: %v", fileInfo.Name())
		}

		err = os.Chmod(pathToGeneratedFile,
			// Read-only.
			0444)
		if err != nil {
			return fmt.Errorf("Failed to change mode on the create assets file: %v",
				pathToGeneratedFile)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("You must pass a path to the site. No less, no more - just that.\n")
		return
	}
	basePath := os.Args[1]
	postsPath, genPath, assetsPath, templatesPath := getSitePaths(basePath)
	files, err := checkStructureAndReturnListOfPosts(basePath, postsPath, genPath)
	if err != nil {
		fmt.Printf("Failure: `%v'\n", err)
		return
	}

	posts, err := makePostsSlice(files, postsPath)
	if err != nil {
		fmt.Printf("Failure: `%v'\n", err)
		return
	}

	err = createPostsFiles(posts, genPath, templatesPath)
	if err != nil {
		fmt.Printf("Failure: `%v'\n", err)
	}

	err = copyAssets(assetsPath, genPath)
	if err != nil {
		fmt.Printf("Failure: `%v'\n", err)
	}
}
