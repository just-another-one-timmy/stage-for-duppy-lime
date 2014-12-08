# Intro
"Stage for Duppy Lime" is the static site generator that I'm going to use for my site - http://iaro.jp
It's a static site generator for people with small brain.
The name stands for *sta*tic site *ge*nerator *for* *du*mb *p*eo*p*l *li*ke *me*.

Jekyll, Octopress and etc are all great.
But it's too hard to setup and use.
I want to have something dumb simple.

# Usage
To use this tool, you need to have a directory with predefined structure.
Then, run
`stage-for-duppy-lime /path/to/site/`
and your site will be generated under the
`/path/to/site/gen/`

If the `gen` directory existed, it will be deleted and recreated.
All existing content will be lost.

# Site structure
Assuming you are in /path/to/site/, the following structure is expected:
`./posts/` contains all blog posts.

There must be a post with name `index` that will be used as a starting page and will contain links to all other blog posts.

## Post structure
Each post is a simple text file.
File with name `file` will correspond to `./gen/posts/file.html` in the generated site.
Each post has a header:
```
date: yyyy mm dd
tags: tag1 tag2 tag3

(post content here)
```

Tags must not have spaces.
No tags is equivalent to tag 'not-tagged'.
After the date and tags list and blank line there comes a post in the markdown format.

# Not supported
What's not supported now:

* Images
* Links between posts (unless you use knowledge about mapping from file names to posts)
* Showing a cut from post and having "See more" button
* Well, anything else, really.

