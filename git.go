package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type repos struct {
	name    string
	root    string
	gitPath string
}

type gitStatus struct {
	isGit        bool
	repos        []repos
	branch       string
	tag          string
	pathFromRoot string

	wtAdded     int
	wtModified  int
	wtUntracked int
	wtConflict  int

	commitMinus int
	commitPlus  int
}

func (s gitStatus) infos() (res []string) {
	names := []string{}
	for _, r := range s.repos {
		names = append(names, color(r.name, Accent, true))
	}

	res = append(res, strings.Join(names, color("/", Neutral, false)))

	branchColor := Success
	switch {
	case s.wtModified > 0:
		branchColor = Danger
	case s.wtAdded > 0:
		branchColor = Yellow
	case s.wtUntracked > 0:
		branchColor = Purple
	}

	res = append(res, color(s.branch, branchColor, false))
	head := ""
	if s.tag != "" {
		head += color(s.tag, Text, false)
	}
	if s.commitMinus > 0 || s.commitPlus > 0 {
		if s.commitPlus > 0 {
			head += color("↑", Neutral, false)
			head += color(strconv.Itoa(s.commitPlus), Green, true)
		}

		if s.commitMinus > 0 {
			head += color("↓", Neutral, false)
			head += color(strconv.Itoa(s.commitMinus), Blue, true)
		}
	}

	if head != "" {
		res = append(res, head)
	}

	tree := ""
	if s.wtAdded > 0 || s.wtModified > 0 || s.wtUntracked > 0 || s.wtConflict > 0 {
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
		if s.wtConflict > 0 {
			parts = append(parts, color(strconv.Itoa(s.wtConflict), Blue, true))
		}

		tree += strings.Join(parts, color(".", Neutral, false))
	}

	if tree != "" {
		res = append(res, tree)
	}

	return
}

func (s *gitStatus) fillRepos(p string) {
	s.isGit = false

	for {
		gitPath := path.Join(p, ".git")
		stat, err := os.Stat(gitPath)
		if os.IsNotExist(err) {
			parent := path.Dir(p)
			if parent == p {
				break
			}
			p = parent
			continue
		}
		if err != nil {
			errorAdd(err)
			break
		}

		if !stat.IsDir() {
			c, err := os.ReadFile(gitPath)
			if err != nil {
				errorAdd(err)
				break
			}
			gitPath = path.Join(p, strings.TrimSpace(string(c)))
		}

		s.isGit = true
		s.repos = append([]repos{{
			name:    path.Base(p),
			root:    p,
			gitPath: gitPath,
		}}, s.repos...)

		if stat.IsDir() {
			break
		} else {
			p = path.Dir(p)
		}
	}
}

func gitBranch() string {
	defer measure("git branch", time.Now())

	res, err := run("git", "branch", "--show-current")
	if err != nil {
		errorAdd(err)
		return "unknwn"
	}
	if res != "" {
		return res
	}

	commit, err := run("git", "log", "--pretty=format:%h", "-n", "1")
	if err != nil {
		errorAdd(err)
		return "unknwn"
	}

	return commit
}

func gitTag(repPath string) string {
	defer measure("git tag", time.Now())

	f, err := os.Open(filepath.Join(repPath, ".git/refs/tags"))
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		errorAdd(err)
		return ""
	}
	entries, err := f.ReadDir(100)
	if err != nil {
		errorAdd(err)
		return ""
	}
	if len(entries) == 100 {
		return "?"
	}

	tagOut, err := run("git", "describe", "--exact-match", "--tags")
	if err != nil {
		errorAdd(err)
		return ""
	}
	return tagOut
}

func gitRemote(branch string) string {
	defer measure("git remote", time.Now())

	out, err := run("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err == nil {
		return out
	}

	out, err = run("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "HEAD")
	if err != nil {
		return ""
	}

	if out == "HEAD" {
		return ""
	}

	out, err = run("git", "remote")
	if err != nil {
		errorAdd(err)
		return ""
	}

	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return ""
	}

	return lines[0] + "/" + branch
}

func gitCommitMinus(branch string) int {
	defer measure("git commit-", time.Now())

	out, err := run("git", "log", "--oneline", fmt.Sprintf("..%s", branch))
	if err != nil {
		errorAdd(err)
		return 0
	}
	if out == "" {
		return 0
	}

	lines := strings.Split(out, "\n")
	return len(lines)
}

func gitCommitPlus(branch string) int {
	defer measure("git commit+", time.Now())

	out, err := run("git", "log", "--oneline", fmt.Sprintf("%s..", branch))
	if err != nil {
		errorAdd(err)
		return 0
	}
	if out == "" {
		return 0
	}

	lines := strings.Split(out, "\n")
	return len(lines)
}

func gitWtStatus() (added, modified, untracked, conflict int) {
	defer measure("git status", time.Now())

	out, err := run("git", "status", "--porcelain")
	if err != nil {
		errorAdd(err)
		return
	}

	lines := strings.Split(out, "\n")
NextLine:
	for _, l := range lines {
		if len(l) < 2 {
			continue
		}

		switch l[:2] {
		case "UU":
			conflict++
			continue NextLine
		}

		switch l[1] {
		case 'M', 'U', 'R', 'D':
			modified++

		case '?':
			untracked++
			continue NextLine
		}

		switch l[0] {
		case 'A', 'M', 'R', 'D':
			added++
		}
	}

	return
}

func gitInfo(cwd string) *gitStatus {
	defer measure("git", time.Now())

	status := &gitStatus{}
	status.fillRepos(cwd)
	if !status.isGit {
		return status
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		status.branch = gitBranch()
		status.tag = gitTag(status.repos[len(status.repos)-1].gitPath)
		remoteBranch := gitRemote(status.branch)

		if remoteBranch != "" {
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
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		status.wtAdded, status.wtModified, status.wtUntracked, status.wtConflict = gitWtStatus()
	}()

	wg.Wait()

	return status
}
