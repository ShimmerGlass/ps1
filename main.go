package main

import (
	"fmt"
	"strings"

	"flag"
)

var why bool

func main() {
	flag.Var(&Accent, "accent-color", "")
	flag.Var(&Text, "text-color", "")
	flag.Var(&Neutral, "neutral-color", "")
	flag.Var(&Danger, "danger-color", "")
	flag.Var(&Warning, "warning-color", "")
	flag.Var(&Success, "success-color", "")
	flag.BoolVar(&why, "why", false, "Debug performance")
	flag.Parse()

	if !why {
		defer fmt.Print(colorRst())
	}

	cwd := getCwd()
	gitInfo := gitInfo(cwd)

	var cwdBase string
	if gitInfo.isGit {
		cwdBase = gitInfo.repos[0].root
	} else {
		cwdBase = home()
	}

	prettyPath := newPrettyPath(cwd, cwdBase)

	if !why {
		if gitInfo.isGit {
			fmt.Print(title(gitInfo.repos[0].name))
		} else {
			fmt.Print(title(prettyPath.string()))
		}
	}

	parts := []string{}

	parts = append(parts, jobs()...)

	if gitInfo.isGit {
		lastRepos := gitInfo.repos[len(gitInfo.repos)-1]
		parts = append(parts, gitInfo.infos()...)
		parts = append(parts, versions(lastRepos.root)...)
	}

	parts = append(parts, ssh()...)
	parts = append(parts, errorCount()...)

	if !why {
		fmt.Print(strings.Join(parts, color("⋮ ", Neutral, false)) + " ")
		fmt.Print(strings.Join(prettyPath.print(), " ") + " ")
		fmt.Print(strings.Join(prompt(), " ") + " ")
	} else {
		printDebugTimes()
		printErrors()
	}
}
