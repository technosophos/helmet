package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

var globalOpts = optList{
	{Suggest: desc("--debug", "enable verbose output")},
	{Suggest: desc("--home", "location of your Helm config. Overrides $HELM_HOME.")},
	{Suggest: desc("--host", "address of Tiller, Overrides $HELM_HOST")},
	{Suggest: desc("--kube-context", "name of the kubeconfig context to use")},
	{Suggest: desc("--tiller-namespace", "namespace of Tiller (default: kube-system)")},
}

var tlsOpts = optList{
	{Suggest: desc("--tls", "enable TLS for a request")},
	{Suggest: desc("--tls-ca-cert", "path to TLS CA certificate file")},
	{Suggest: desc("--tls-cert", "path to TLS certificate file")},
	{Suggest: desc("--tls-key", "path to TLS key file")},
	{Suggest: desc("--tls-verify", "enable TLS for request and verify remote")},
}

var getOpts = optList{
	{Suggest: desc("--revision", "revision number")},
}

var listOpts = optList{
	{Suggest: desc("--all", "show all releases, not just DEPLOYED")},
	{Suggest: desc("--date", "sort release by date")},
	{Suggest: desc("--deleted", "show deleted releases")},
	{Suggest: desc("--deleting", "show releases that are currently being deleted")},
	{Suggest: desc("--deployed", "show deployed releases (default filter")},
	{Suggest: desc("--failed", "show failed releases")},
	{Suggest: desc("--max", "maximum number of releases to fetch")},
	{Suggest: desc("--namespace", "show releases in this namespace")},
	{Suggest: desc("--offset", "next release name in list")},
	{Suggest: desc("--pending", "show pending releases")},
	{Suggest: desc("--reverse", "reverse the sort order")},
	{Suggest: desc("--short", "output short (quite) listing format")},
}

var subcommands = cmdList{
	{Suggest: desc("help", "print help"), Resolve: empty},
	{Suggest: desc("home", "print Helm home"), Resolve: empty},
	{
		Suggest: desc("version", "Helm and Tiller version"),
		Options: optBuilder(tlsOpts, globalOpts, []opt{
			{Suggest: desc("--client", "client version only")},
			{Suggest: desc("--server", "server version only")},
			{Suggest: desc("--short", "print the version number")},
		}),
		Resolve: empty,
	},
	{
		Suggest: desc("get", "get details of a release"),
		Subcommands: cmdList{
			{
				Suggest: desc("hooks", "get all hooks for a release"),
				Options: optBuilder(getOpts, tlsOpts, globalOpts),
				Resolve: listReleases,
			},
			{
				Suggest: desc("manifest", "get manifest for a release"),
				Options: optBuilder(getOpts, tlsOpts, globalOpts),
				Resolve: listReleases,
			},
			{
				Suggest: desc("values", "get values for a release"),
				Options: optBuilder(getOpts, tlsOpts, globalOpts),
				Resolve: listReleases,
			},
		},
		Options: optBuilder(getOpts, tlsOpts, globalOpts),
		Resolve: listReleases,
	},
	{
		Suggest: desc("list", "list releases"),
		Options: optBuilder(listOpts, tlsOpts, globalOpts),
		Resolve: empty,
	},
	{
		Suggest: desc("ls", "list releases"),
		Options: optBuilder(listOpts, tlsOpts, globalOpts),
		Resolve: empty,
	},
	{
		Suggest: desc("search", "search for charts"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--regexp", "treat search string as regular expression")},
			{Suggest: desc("--version", "search using semver constraints")},
			{Suggest: desc("--versions", "show long listing, including each version of the chart")},
		}),
		Resolve: empty,
	},
	{
		Suggest: desc("serve", "serve local chart repo (not recommended in Helmet)"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--address", "address to listen on")},
			{Suggest: desc("--repo-path", "path to charts to serve")},
			{Suggest: desc("--url", "external URL of the repository")},
		}),
		Resolve: empty,
	},
	{
		Suggest: desc("test", "test a release"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--cleanup", "delete test pods upon completion")},
			{Suggest: desc("--timeout", "timeout in seconds to wait for an operation")},
		}),
		Resolve: listReleases,
	},
	{
		Suggest: desc("delete", "delete a release"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--dry-run", "simulate a delete")},
			{Suggest: desc("--no-hooks", "prevent hooks from running during deletion")},
			{Suggest: desc("--purge", "remove the release from storage and free its name for later use")},
			{Suggest: desc("--timeout", "timeout in seconds to wait for an operation")},
		}),
		Resolve: listReleases,
	},
	{
		Suggest: desc("rollback", "rollback a release"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--dry-run", "simulate a rollback")},
			{Suggest: desc("--force", "force resource update through delete/recreate if needed")},
			{Suggest: desc("--no-hooks", "prevent hooks from running during rollback")},
			{Suggest: desc("--recreate-pods", "performs pod restart if necessary")},
			{Suggest: desc("--timeout", "timeout in seconds to wait for an operation")},
			{Suggest: desc("--wait", "wait for resources to enter ready state before returning")},
		}),
		Resolve: listReleases,
	},
	{
		Suggest: desc("history", "get the history of a release"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--max", "maximum number of revisions to include in history")},
		}),
		Resolve: listReleases,
	},
	{
		Suggest: desc("status", "show status of a release"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--revision", "revision number")},
		}),
		Resolve: listReleases,
	},
	{
		Suggest: desc("lint", "check a chart for compliance"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--strict", "fail lint on warnings")},
		}),
		Resolve: listReleases,
	},
	{
		Suggest: desc("upgrade", "upgrade a release"),
		Options: optBuilder(tlsOpts, globalOpts, optList{
			{Suggest: desc("--ca-file", "verify certificates of HTTPS-enabled servers using this CA bundle")},
			{Suggest: desc("--cert-file", "identify HTTPS client using this SSL certificate file")},
			{Suggest: desc("--devel", "use development versions, too")},
			{Suggest: desc("--dry-run", "simulate aun upgrade")},
			{Suggest: desc("--force", "force resource update through delete/recreate if needed")},
			{Suggest: desc("--install", "if a release by this name does not exist, install it.")},
			{Suggest: desc("--key-file", "identify HTTPS client using this SSL key file")},
			{Suggest: desc("--keyring", "path to the kerying that contains public singing keys")},
			{Suggest: desc("--namespace", "namespace to install into")},
			{Suggest: desc("--no-hooks", "prevent hooks from running during upgrade")},
			{Suggest: desc("--recreate-pods", "performs pod restart if necessary")},
			{Suggest: desc("--repo", "chart repository URL where requested chart is located")},
			{Suggest: desc("--reset-values", "when upgrading, reset the values to the ones built into the chart")},
			{Suggest: desc("--reuse-values", "when upgrading, reuse the values sent with the last operation and merge in new values")},
			{Suggest: desc("--set", "set one or more values '--set key1=val1,key2=val2'")},
			{Suggest: desc("--timeout", "timeout in seconds to wait for an operation")},
			{Suggest: desc("--values", "path to values YAML file")},
			{Suggest: desc("-f", "path to values YAML file")},
			{Suggest: desc("--verify", "verify provenance of chart before upgrading")},
			{Suggest: desc("--version", "specify the exact chart version to use")},
			{Suggest: desc("--wait", "wait for resources to enter ready state before returning")},
		}),
		Resolve: listReleases,
	},
}

func optBuilder(sets ...optList) optList {
	// TODO: optimize this function
	ret := []opt{}
	for _, s := range sets {
		ret = append(ret, s...)
	}
	return ret
}

func desc(text, desc string) prompt.Suggest {
	return prompt.Suggest{Text: text, Description: desc}
}

type cmd struct {
	prompt.Suggest
	Subcommands cmdList
	Options     optList
	Resolve     resolver
}

func (c cmd) suggestFor(args []string) []prompt.Suggest {
	if len(args) == 0 {
		return noSuggestions
	}
	// If arg0 starts with --, send opts
	a0 := args[0]
	if strings.HasPrefix(a0, "-") {
		return prompt.FilterHasPrefix(c.Options.suggestions(), a0, true)
	}
	sug := c.Subcommands.suggestions()
	sug = append(sug, c.Resolve(args)...)
	return prompt.FilterHasPrefix(sug, a0, true)
}

type cmdList []cmd

func (cmds cmdList) suggestions() []prompt.Suggest {
	ret := make([]prompt.Suggest, len(cmds))
	for i, d := range cmds {
		ret[i] = d.Suggest
	}
	return ret
}

func (cmds cmdList) get(name string) (cmd, bool) {
	for _, c := range cmds {
		if c.Text == name {
			return c, true
		}
	}
	return cmd{}, false
}

type opt struct {
	prompt.Suggest
}

type optList []opt

func (o optList) suggestions() []prompt.Suggest {
	ret := make([]prompt.Suggest, len(o))
	for i, d := range o {
		ret[i] = d.Suggest
	}
	return ret
}
