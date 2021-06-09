/*
Copyright 2020 The goARRG Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"goarrg.com/debug"
)

type CMD struct {
	Run   func([]string) bool
	Name  string
	Usage string
	Short string
	Long  string
	Flags flag.FlagSet
	CMDs  map[string]*CMD
}

func (cmd *CMD) usage() {
	args := ""
	cmds := ""

	cmd.Flags.VisitAll(func(f *flag.Flag) {
		n, u := flag.UnquoteUsage(f)
		args += "\t-" + f.Name + " " + n + "\n\t\t" + strings.ReplaceAll(strings.TrimSpace(u), "\n", "\n\t\t") + "\n"
	})

	sortedCMDs := make([]string, 0, len(cmd.CMDs))

	for _, child := range cmd.CMDs {
		sortedCMDs = append(sortedCMDs, child.Name)
	}

	sort.Strings(sortedCMDs)

	for _, child := range sortedCMDs {
		cmds += "\t" + cmd.CMDs[child].Name + "\n\t\t" + strings.ReplaceAll(strings.TrimSpace(cmd.CMDs[child].Short), "\n", "\n\t\t") + "\n"
	}

	fmt.Fprintf(os.Stderr, "Usage:\n\t%s %s", "go run goarrg.com/cmd/goarrg", cmd.Name)

	switch {
	case cmds != "" && args != "":
		fmt.Fprintf(os.Stderr, " [command] [arguments] %s\n\nCommands:\n%s\nArguments:\n%s", cmd.Usage, cmds, args)
	case cmds != "":
		fmt.Fprintf(os.Stderr, " [command] %s\n\nCommands:\n%s", cmd.Usage, cmds)
	case args != "":
		fmt.Fprintf(os.Stderr, " [arguments] %s\n\nArguments:\n%s", cmd.Usage, args)
	default:
		fmt.Fprintf(os.Stderr, "%s\n", cmd.Usage)
	}

	if cmd.Long != "" {
		fmt.Fprintf(os.Stderr, "\nDescription:\n\t%s\n", strings.TrimSpace(cmd.Long))
	} else if cmd.Short != "" {
		fmt.Fprintf(os.Stderr, "\nDescription:\n\t%s\n", strings.TrimSpace(cmd.Short))
	}
}

func (cmd *CMD) traverse(args []string) (*CMD, []string) {
	if cmd.CMDs == nil || len(cmd.CMDs) == 0 || len(args) == 0 {
		return cmd, args
	}

	child, ok := cmd.CMDs[args[0]]

	if !ok {
		return cmd, args
	}

	child.Name = cmd.Name + " " + child.Name
	return child.traverse(args[1:])
}

func (cmd *CMD) Exec(args []string) {
	cmd, args = cmd.traverse(args)
	cmd.Flags.Usage = cmd.usage
	cmd.Flags.Init("", flag.ExitOnError)
	cmd.Flags.BoolVar(&flagVerbose, "v", false, "Verbose - Print high level tasks")
	cmd.Flags.BoolVar(&flagVeryVerbose, "vv", false, "Very Verbose - Print everything")
	_ = cmd.Flags.Parse(args)

	switch {
	case flagVeryVerbose:
		debug.LogSetLevel(debug.LogLevelVerbose)
	case flagVerbose:
		debug.LogSetLevel(debug.LogLevelInfo)
	default:
		debug.LogSetLevel(debug.LogLevelWarn)
	}

	if !cmd.Run(cmd.Flags.Args()) {
		cmd.usage()
		os.Exit(2)
	}
}
