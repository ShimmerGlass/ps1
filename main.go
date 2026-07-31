package main

import (
	"fmt"
	"strings"
	"sync"

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

	// fillRepos is a couple of stat() calls, so resolving the repo root up
	// front is cheap and lets git status and version detection (both several
	// milliseconds of subprocess time) run concurrently instead of in series.
	gitInfo := &gitStatus{}
	gitInfo.fillRepos(cwd)

	var wg sync.WaitGroup
	var verParts []string

	if gitInfo.isGit {
		root := gitInfo.repos[len(gitInfo.repos)-1].root

		wg.Add(1)
		go func() {
			defer wg.Done()
			gitInfo.fillGit()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			verParts = versions(root)
		}()
	}

	wg.Wait()

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
		parts = append(parts, gitInfo.infos()...)
		parts = append(parts, verParts...)
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
