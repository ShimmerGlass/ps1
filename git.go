package main

import (
	"bufio"
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

	commitMinus string
	commitPlus  string
}

func (s gitStatus) infos() (res []string) {
	names := []string{}
	for _, r := range s.repos {
		names = append(names, color(r.name, Accent, true))
	}

	res = append(res, strings.Join(names, color("/", Neutral, false)))

	branchColor := Success
	switch {
	case s.wtConflict > 0:
		branchColor = Blue
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
	if s.commitMinus == "?" && s.commitPlus == "?" {
		head += color("↕?", Yellow, true)
	} else {
		if s.commitMinus != "" || s.commitPlus != "" {
			if s.commitPlus != "" {
				head += color("↑", Neutral, false)
				head += color(s.commitPlus, Green, true)
			}

			if s.commitMinus != "" {
				head += color("↓", Neutral, false)
				head += color(s.commitMinus, Blue, true)
			}
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
			// A .git file (submodule / linked worktree) points at the real git
			// directory as "gitdir: <path>", where <path> may be relative to p.
			dir := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(c)), "gitdir:"))
			if filepath.IsAbs(dir) {
				gitPath = dir
			} else {
				gitPath = filepath.Join(p, dir)
			}
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

// gitStatusPorcelain gathers branch, upstream tracking (ahead/behind) and
// working-tree status in a single `git status` invocation instead of the five
// separate git subprocesses this used to spawn.
func gitStatusPorcelain(s *gitStatus) {
	defer measure("git status", time.Now())

	out, err := run("git", "status", "--porcelain=v2", "--branch")
	if err != nil {
		errorAdd(err)
		s.branch = "unknwn"
		return
	}

	var oid string
	hasUpstream := false
	hasAB := false

	for _, l := range strings.Split(out, "\n") {
		if l == "" {
			continue
		}

		if l[0] == '#' {
			fields := strings.SplitN(l, " ", 3)
			if len(fields) < 3 {
				continue
			}
			switch fields[1] {
			case "branch.oid":
				oid = fields[2]
			case "branch.head":
				if fields[2] != "(detached)" {
					s.branch = fields[2]
				}
			case "branch.upstream":
				hasUpstream = true
			case "branch.ab":
				hasAB = true
				ab := strings.Fields(fields[2])
				if len(ab) == 2 {
					if ahead := strings.TrimPrefix(ab[0], "+"); ahead != "0" {
						s.commitPlus = ahead
					}
					if behind := strings.TrimPrefix(ab[1], "-"); behind != "0" {
						s.commitMinus = behind
					}
				}
			}
			continue
		}

		s.countStatus(l)
	}

	// Detached HEAD: show the short commit hash, like the old fallback did.
	if s.branch == "" {
		if len(oid) >= 7 {
			s.branch = oid[:7]
		} else {
			s.branch = "unknwn"
		}
	}

	// Upstream is configured but git couldn't compute ahead/behind (e.g. the
	// upstream is gone): surface it as the "↕?" unknown marker, as before.
	if hasUpstream && !hasAB {
		s.commitPlus = "?"
		s.commitMinus = "?"
	}
}

// countStatus tallies a single porcelain v2 status line into the working-tree
// counters: untracked ('?'), conflict (any 'u' line) or, for ordinary/renamed
// changes, staged (X) and worktree (Y) modifications from the XY code.
func (s *gitStatus) countStatus(l string) {
	switch l[0] {
	case '?':
		s.wtUntracked++
		return
	case 'u':
		// Every 'u' line is an unmerged path, i.e. a conflict, regardless of
		// its XY code (UU, AA, DD, AU, UD, UA, DU).
		s.wtConflict++
		return
	case '!':
		return
	case '1', '2':
		// Ordinary or renamed change, handled below.
	default:
		return
	}

	fields := strings.SplitN(l, " ", 3)
	if len(fields) < 2 || len(fields[1]) < 2 {
		return
	}
	xy := fields[1]

	switch xy[1] {
	case 'M', 'R', 'D':
		s.wtModified++
	}
	switch xy[0] {
	case 'A', 'M', 'R', 'D':
		s.wtAdded++
	}
}

func gitTag(gitDir string) string {
	defer measure("git tag", time.Now())

	if !hasTags(gitDir) {
		return ""
	}

	tagOut, err := run("git", "describe", "--exact-match", "--tags")
	if err != nil {
		// A non-zero exit here means HEAD isn't exactly on a tag — the common
		// case, not a failure worth reporting.
		return ""
	}
	return tagOut
}

// hasTags reports whether the repository has any tags, checking both loose and
// packed refs. gitDir is the repository's git directory; for a linked worktree
// its refs live in the shared common directory named by the commondir file.
func hasTags(gitDir string) bool {
	commonDir := gitDir
	if c, err := os.ReadFile(filepath.Join(gitDir, "commondir")); err == nil {
		if cd := strings.TrimSpace(string(c)); cd != "" {
			if filepath.IsAbs(cd) {
				commonDir = cd
			} else {
				commonDir = filepath.Join(gitDir, cd)
			}
		}
	}

	if entries, err := os.ReadDir(filepath.Join(commonDir, "refs", "tags")); err == nil && len(entries) > 0 {
		return true
	}

	f, err := os.Open(filepath.Join(commonDir, "packed-refs"))
	if err != nil {
		return false
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.Contains(sc.Text(), " refs/tags/") {
			return true
		}
	}
	return false
}

// fillGit runs the working-tree status and tag lookup concurrently. It assumes
// fillRepos has already run and reported a git repository.
func (s *gitStatus) fillGit() {
	defer measure("git", time.Now())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		gitStatusPorcelain(s)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.tag = gitTag(s.repos[len(s.repos)-1].gitPath)
	}()

	wg.Wait()
}
