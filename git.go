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
	branch         string
	tag            string
	pathFromRoot   string

	wtAdded     int
	wtModified  int
	wtUntracked int

	commitMinus int
	commitPlus  int
}

func (s gitStatus) infos() (res []string) {
	res = append(res,
		color(reposName(s.repositoryName), Text, true),
		color("∙", Neutral, false),
	)

	branchColor := Success
	switch {
	case s.wtModified > 0:
		branchColor = Danger
	case s.wtAdded > 0:
		branchColor = Yellow
	case s.wtUntracked > 0:
		branchColor = Purple
	}

	head := ""
	head += color(s.branch, branchColor, false)
	if s.tag != "" {
		head += color("∙", Neutral, false)
		head += color(s.tag, Text, false)
	}
	if s.commitMinus > 0 || s.commitPlus > 0 {
		head += color("{", Neutral, false)

		if s.commitPlus > 0 {
			head += color("↑", Green, false)
			head += color(strconv.Itoa(s.commitPlus), Green, true)
		}

		if s.commitMinus > 0 {
			head += color("↓", Blue, false)
			head += color(strconv.Itoa(s.commitMinus), Blue, true)
		}

		head += color("}", Neutral, false)
	}

	res = append(res, head)

	tree := ""
	if s.wtAdded > 0 || s.wtModified > 0 || s.wtUntracked > 0 {
		tree += color("⟨", Neutral, false)

		parts := []string{}

		if s.wtAdded > 0 {
			parts = append(parts, color(strconv.Itoa(s.wtAdded), Yellow, true))
		}
		if s.wtModified > 0 {
			parts = append(parts, color(strconv.Itoa(s.wtModified), Danger, true))
		}
		if s.wtUntracked > 0 {
			parts = append(parts, color(strconv.Itoa(s.wtUntracked), Purple, true))
		}

		tree += strings.Join(parts, color(".", Neutral, false))

		tree += color("⟩", Neutral, false)
	}

	if tree != "" {
		res = append(res, tree)
	}

	return
}

func reposName(v string) string {
	return v
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
		return "unknwn"
	}

	res := strings.TrimSpace(string(out))
	if res != "HEAD" {
		return res
	}

	commit, err := exec.Command("git", "log", "--pretty=format:%h", "-n", "1").Output()
	if err != nil {
		return "unknwn"
	}

	msgb, err := exec.Command("git", "log", "--pretty=format:%s", "-n", "1").Output()
	if err != nil {
		return "unknwn"
	}

	msg := strings.TrimSpace(string(msgb))
	if len(msg) > 20 {
		msg = msg[:20]
	}

	return strings.TrimSpace(string(commit)) + "(" + msg + ")"
}

func gitTag() string {
	tagOut, _ := exec.Command("git", "describe", "--exact-match", "--tags").Output()
	return strings.TrimSpace(string(tagOut))
}

func gitRemote(branch string) string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	out, err = exec.Command("git", "remote").Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) == 0 {
		return ""
	}

	return lines[0] + "/" + branch
}

func gitCommitMinus(branch string) int {
	out, err := exec.Command("git", "log", "--oneline", fmt.Sprintf("..%s", branch)).Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(out), "\n")
	return len(lines) - 1
}

func gitCommitPlus(branch string) int {
	out, err := exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..", branch)).Output()
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
		status.branch = gitBranch()
		status.tag = gitTag()
		remoteBranch := gitRemote(status.branch)

		var wg2 sync.WaitGroup

		wg2.Add(1)
		go func() {
			defer wg2.Done()
			status.commitMinus = gitCommitMinus(remoteBranch)

		}()

		wg2.Add(1)
		go func() {
			defer wg2.Done()
			status.commitPlus = gitCommitPlus(remoteBranch)
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
