package main

func main() {
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

	ptime()

	if gitInfo.isGit {
		gitInfo.pinfos()
	}

	prettyPath.print()
	pprompt()
}
