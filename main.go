package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "cmd" {
		fmt.Println("export PS1='`ps1 $? \\j`'")
		return
	}

	defer pcolorRst()

	cwd := getCwd()
	gitInfo := gitInfo(cwd)

	var cwdBase string
	if gitInfo.isGit {
		cwdBase = gitInfo.repositoryRoot
	} else {
		cwdBase = home()
	}

	prettyPath := newPrettyPath(cwd, cwdBase)

	if gitInfo.isGit {
		ptitle(gitInfo.repositoryName)
	} else {
		ptitle(prettyPath.string())
	}

	pjobs()

	if gitInfo.isGit {
		gitInfo.pinfos()
	}

	prettyPath.print()
	pprompt()
}
