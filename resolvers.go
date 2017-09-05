package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

type resolver func(input []string) []prompt.Suggest

var lsCache []prompt.Suggest

func empty(v []string) []prompt.Suggest {
	return noSuggestions
}

func listReleases(v []string) []prompt.Suggest {
	if len(lsCache) == 0 {
		releases, err := helmCmd("list", "-m", "500", "-d", "-r", "-q")
		if err != nil {
			panic(err)
			fmt.Fprintln(os.Stderr, err)
			return noSuggestions
		}
		ret := make([]prompt.Suggest, len(releases))
		for i, r := range releases {
			ret[i] = prompt.Suggest{Text: r, Description: "(release)"}
		}
		lsCache = ret
	}
	return lsCache
}

// helmCmd runs a Helm command and returns the output as an array of lines.
func helmCmd(in ...string) ([]string, error) {
	d, err := exec.Command("helm", in...).Output()
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(d), "\n")
	return lines, nil
}
