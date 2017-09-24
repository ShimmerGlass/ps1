package main

func main() {
	defer pcolorRst()

	cwd := getCwd()
	gitInfo := gitInfo(cwd)

	ptime()

	if gitInfo.isGit {
		gitInfo.pinfos()
	}

	var cwdBase string
	if gitInfo.isGit {
		cwdBase = gitInfo.repositoryRoot
	} else {
		cwdBase = home()
	}

	pcwd(cwd, cwdBase)
	pprompt()
}
