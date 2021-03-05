package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "cmd" {
		fmt.Println("export PS1='`ps1 $? \\j`'")
		return
	}

	defer fmt.Print(colorRst())

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
		fmt.Print(title(gitInfo.repositoryName))
	} else {
		fmt.Print(title(prettyPath.string()))
	}

	parts := []string{}

	parts = append(parts, jobs()...)

	if gitInfo.isGit {
		parts = append(parts, gitInfo.infos()...)
		parts = append(parts, rubyVersion(gitInfo.repositoryRoot)...)
		parts = append(parts, goVersion(gitInfo.repositoryRoot)...)
	}

	parts = append(parts, prettyPath.print()...)
	parts = append(parts, ssh()...)
	parts = append(parts, prompt()...)

	fmt.Print(strings.Join(parts, " ") + " ")
}
