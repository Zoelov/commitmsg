package main

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"commitmsg/git"

	"github.com/apex/log"
	"github.com/fatih/color"
)

func getGitLog() ([]string, error) {
	var args = []string{"log", "--pretty=oneline", "--abbrev-commit", "--no-color", "--max-count=4"}
	logs, err := git.Run(args...)
	if err != nil {
		return nil, err
	}

	var entries = strings.Split(logs, "\n")
	return entries[0 : len(entries)-1], nil
}

type ErrList struct {
	errs []string
}

func (e *ErrList) Error() string {
	msg := strings.Join(e.errs, "\n")
	return msg
}

func main() {
	entries, err := getGitLog()
	if err != nil {
		log.WithError(err).Error("get git log error")
		os.Exit(1)
		return
	}

	errCommits := make([]string, 0)
	for _, commitMsg := range entries {
		splitMsgs := strings.SplitN(commitMsg, " ", 2)
		message := splitMsgs[1]
		err = valid(message)
		if err != nil {
			errCommits = append(errCommits, commitMsg)
			continue
		}

		log.Debugf("%v passed.", commitMsg)
	}

	if len(errCommits) == 0 {
		log.Info("success")
		return
	}

	err = &ErrList{errCommits}
	log.WithError(err).Error("valid commit msg failed")
	os.Exit(1)
}

const CommitMessagePattern = `^(?:fixup!\s*)?(\w*)(\(([\w\$\.\*/-].*)\))?\: (.*)|^Merge\ branch(.*)`

type CommitType string

const (
	FEAT     CommitType = "feat"
	FIX      CommitType = "fix"
	DOCS     CommitType = "docs"
	DOC      CommitType = "doc"
	STYLE    CommitType = "style"
	REFACTOR CommitType = "refactor"
	TEST     CommitType = "test"
	CHORE    CommitType = "chore"
	PERF     CommitType = "perf"
	HOTFIX   CommitType = "hotfix"
)

func valid(message string) error {
	if strings.HasPrefix(message, "Merge branch") {
		return nil
	}
	if strings.HasPrefix(message, "Merge remote") {
		return nil
	}

	var commitMsgReg = regexp.MustCompile(CommitMessagePattern)

	commitTypes := commitMsgReg.FindAllStringSubmatch(message, -1)
	if len(commitTypes) != 1 {
		msg := color.New(color.Bold).Sprintf("[%v] not match", message)
		return errors.New(msg)
	} else {
		switch commitTypes[0][1] {
		case string(FEAT):
		case string(FIX):
		case string(DOCS):
		case string(DOC):
		case string(STYLE):
		case string(REFACTOR):
		case string(TEST):
		case string(CHORE):
		case string(PERF):
		case string(HOTFIX):
		default:
			msg := color.New(color.Bold).Sprintf("[%v] not match", message)
			return errors.New(msg)
		}
	}

	return nil
}
