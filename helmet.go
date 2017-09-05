package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

func main() {
	p := prompt.New(
		execute,
		complete,
		prompt.OptionTitle("helmet: interactive Helm client (type \"q\" to exit)"),
		prompt.OptionPrefix("⎈ > "),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	fmt.Fprintln(os.Stdout, "helmet (type \"q\" to quit)")
	p.Run()
}

func execute(s string) {
	s = strings.TrimSpace(s)
	switch s {
	case "":
		return
	case "?", "h":
		fmt.Fprintln(os.Stdout, "This is a wrapper around the 'helm' command.")
	case "quit", "exit", "q":
		fmt.Fprintln(os.Stdout, "Bye!")
		os.Exit(0)
	default:
		execHelm(s)
	}
}

func execHelm(s string) {
	c := exec.Command("/bin/sh", "-c", "helm "+s)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}

func complete(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{}
	text := d.TextBeforeCursor()
	if text == "" || strings.Contains(text, "|") {
		return suggestions
	}

	args := strings.Split(text, " ")
	//word := d.GetWordBeforeCursor()

	return completeSubcommands(args)
}

func completeOptions(args []string, longOpt bool) []prompt.Suggest {
	return globalOpts.suggestions()
}

var noSuggestions = []prompt.Suggest{}

func completeSubcommands(args []string) []prompt.Suggest {

	if len(args) < 2 {
		return prompt.FilterHasPrefix(subcommands.suggestions(), args[0], true)
	}

	sc, ok := subcommands.get(args[0])
	if !ok {
		return noSuggestions
	}

	return sc.suggestFor(args[1:])
}
