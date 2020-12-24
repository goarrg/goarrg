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

package base

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"goarrg.com/debug"
)

type CMD struct {
	Run   func([]string) bool
	Name  string
	Short string
	Long  string
	Flag  flag.FlagSet
	CMDs  map[string]*CMD
}

func (cmd *CMD) usage() {
	args := ""
	cmds := ""

	cmd.Flag.VisitAll(func(f *flag.Flag) {
		n, u := flag.UnquoteUsage(f)
		args += fmt.Sprintf("\t-" + f.Name + " " + n + "\n\t\t" + u + "\n")
	})

	for _, child := range cmd.CMDs {
		cmds += fmt.Sprintf("\t" + child.Name + "\n\t\t" + child.Short + "\n")
	}

	fmt.Fprintf(os.Stderr, "Usage:\n\t"+filepath.Base(os.Args[0])+" "+cmd.Name)

	switch {
	case cmds != "" && args != "":
		fmt.Fprintf(os.Stderr, " [command] [arguments]\n\nCommands:\n"+cmds+"\nArguments:\n"+args)
	case cmds != "":
		fmt.Fprintf(os.Stderr, " [command]\n\nCommands:\n"+cmds)
	case args != "":
		fmt.Fprintf(os.Stderr, " [arguments]\n\nArguments:\n"+args)
	default:
		fmt.Fprintf(os.Stderr, "\n")
	}

	if cmd.Long != "" {
		fmt.Fprintf(os.Stderr, "\nDescription:\n\t%s\n", cmd.Long)
	} else if cmd.Short != "" {
		fmt.Fprintf(os.Stderr, "\nDescription:\n\t%s\n", cmd.Short)
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
	cmd.Flag.Usage = cmd.usage
	cmd.Flag.Init("", flag.ExitOnError)
	cmd.Flag.BoolVar(&flagVerbose, "v", false, "Verbose - Print high level tasks")
	cmd.Flag.BoolVar(&flagVeryVerbose, "vv", false, "Very Verbose - Print everything")
	_ = cmd.Flag.Parse(args)

	switch {
	case flagVeryVerbose:
		debug.LogSetLevel(debug.LogLevelVerbose)
	case flagVerbose:
		debug.LogSetLevel(debug.LogLevelInfo)
	default:
		debug.LogSetLevel(debug.LogLevelError)
	}

	if !cmd.Run(cmd.Flag.Args()) {
		cmd.usage()
		os.Exit(2)
	}
}
