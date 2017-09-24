package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
)

type gitStatus struct {
	isGit          bool
	repositoryName string
	repositoryRoot string
	branchName     string
	pathFromRoot   string

	wtAdded     int
	wtModified  int
	wtUntracked int

	commitMinus int
	commitPlus  int
}

func (s gitStatus) pinfos() {
	branchColor := Green
	if s.wtModified > 0 {
		branchColor = Red
	} else if s.wtUntracked > 0 {
		branchColor = Purple
	} else if s.wtAdded > 0 {
		branchColor = Yellow
	}

	pcolor(s.branchName, branchColor, true)

	if s.commitMinus > 0 || s.commitPlus > 0 {
		pcolor("(", Black, true)

		if s.commitPlus > 0 {
			pcolor("+", Green, false)
			pcolor(strconv.Itoa(s.commitPlus), Green, true)
		}

		if s.commitMinus > 0 {
			pcolor("-", Blue, false)
			pcolor(strconv.Itoa(s.commitMinus), Blue, true)
		}

		pcolor(")", Black, true)
	}

	if s.wtAdded > 0 || s.wtModified > 0 || s.wtUntracked > 0 {
		pcolor("[", Black, true)

		if s.wtAdded > 0 {
			pcolor(strconv.Itoa(s.wtAdded), Green, true)
		}
		if s.wtModified > 0 {
			pcolor(strconv.Itoa(s.wtModified), Red, true)
		}
		if s.wtUntracked > 0 {
			pcolor(strconv.Itoa(s.wtUntracked), Purple, true)
		}

		pcolor("]", Black, true)
	}

	pcolor(":", Cyan, false)
	pcolor(s.repositoryName+" ", White, true)
}

func isDirGit(p string) (string, bool) {
	if _, err := os.Stat(path.Join(p, ".git")); os.IsNotExist(err) {
		parent := path.Dir(p)
		if parent == p {
			return "", false
		}

		return isDirGit(parent)
	}
	return p, true
}

func gitBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "master"
	}

	return strings.TrimSpace(string(out))
}

func gitRemote() string {
	out, err := exec.Command("git", "remote").Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) == 0 {
		return ""
	}

	return lines[0]
}

func gitCommitMinus(remote, branch string) int {
	out, err := exec.Command("git", "log", "--oneline", fmt.Sprintf("..%s/%s", remote, branch)).Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(out), "\n")
	return len(lines) - 1
}

func gitCommitPlus(remote, branch string) int {
	out, err := exec.Command("git", "log", "--oneline", fmt.Sprintf("%s/%s..", remote, branch)).Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(out), "\n")
	return len(lines) - 1
}

func gitWtStatus() (added, modified, untracked int) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		if len(l) < 2 {
			continue
		}

		switch l[1] {
		case 'M', 'U', 'R', 'D':
			modified++
		case '?':
			untracked++
		}

		switch l[0] {
		case 'A', 'M', 'R', 'D':
			added++
		}
	}

	return
}

func gitInfo(cwd string) gitStatus {
	status := gitStatus{}

	repPath, isGit := isDirGit(cwd)
	if !isGit {
		return status
	}

	status.repositoryName = path.Base(repPath)
	status.repositoryRoot = repPath
	status.isGit = isGit

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		remote := gitRemote()
		status.branchName = gitBranch()

		var wg2 sync.WaitGroup

		wg2.Add(1)
		go func() {
			defer wg2.Done()
			status.commitMinus = gitCommitMinus(remote, status.branchName)

		}()

		wg2.Add(1)
		go func() {
			defer wg2.Done()
			status.commitPlus = gitCommitPlus(remote, status.branchName)
		}()

		wg2.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		status.wtAdded, status.wtModified, status.wtUntracked = gitWtStatus()
	}()

	wg.Wait()

	return status
}
