// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"
	"github.com/atomix/go-client/pkg/client/set"
	"github.com/spf13/cobra"
)

func newSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {create,add,contains,remove,size,clear,delete}",
		Short: "Manage the state of a distributed set",
	}
	addClientFlags(cmd)
	cmd.AddCommand(newSetCreateCommand())
	cmd.AddCommand(newSetAddCommand())
	cmd.AddCommand(newSetContainsCommand())
	cmd.AddCommand(newSetRemoveCommand())
	cmd.AddCommand(newSetSizeCommand())
	cmd.AddCommand(newSetClearCommand())
	cmd.AddCommand(newSetDeleteCommand())
	return cmd
}

func addSetFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "the set name")
	cmd.Flags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_sets"},
	}
	cmd.MarkPersistentFlagRequired("name")
}

func getSetName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func newSetFromName(cmd *cobra.Command, name string) set.Set {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetSet(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newSetCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create <name>",
		Args: cobra.ExactArgs(1),
		Run:  runSetCreateCommand,
	}
}

func runSetCreateCommand(cmd *cobra.Command, args []string) {
	set := newSetFromName(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	set.Close(ctx)
	ExitWithOutput(fmt.Sprintf("Created %s", set.Name().String()))
}

func newSetDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <name>",
		Args: cobra.ExactArgs(1),
		Run:  runSetDeleteCommand,
	}
}

func runSetDeleteCommand(cmd *cobra.Command, args []string) {
	set := newSetFromName(cmd, args[0])
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := set.Delete(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", set.Name().String()))
	}
}

func newSetAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add <value>",
		Args: cobra.ExactArgs(1),
		Run:  runSetAddCommand,
	}
	addSetFlags(cmd)
	return cmd
}

func runSetAddCommand(cmd *cobra.Command, args []string) {
	set := newSetFromName(cmd, getSetName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	added, err := set.Add(ctx, value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(added)
	}
}

func newSetContainsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "contains <value>",
		Args: cobra.ExactArgs(1),
		Run:  runSetContainsCommand,
	}
	addSetFlags(cmd)
	return cmd
}

func runSetContainsCommand(cmd *cobra.Command, args []string) {
	set := newSetFromName(cmd, getSetName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	contains, err := set.Contains(ctx, value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(contains)
	}
}

func newSetRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "remove <value>",
		Args: cobra.ExactArgs(1),
		Run:  runSetRemoveCommand,
	}
	addSetFlags(cmd)
	return cmd
}

func runSetRemoveCommand(cmd *cobra.Command, args []string) {
	set := newSetFromName(cmd, getSetName(cmd))
	value := args[0]
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	removed, err := set.Remove(ctx, value)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(removed)
	}
}

func newSetSizeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "size",
		Args: cobra.NoArgs,
		Run:  runSetSizeCommand,
	}
	addSetFlags(cmd)
	return cmd
}

func runSetSizeCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd, getSetName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	size, err := set.Len(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(size)
	}
}

func newSetClearCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "clear",
		Args: cobra.NoArgs,
		Run:  runSetClearCommand,
	}
	addSetFlags(cmd)
	return cmd
}

func runSetClearCommand(cmd *cobra.Command, _ []string) {
	set := newSetFromName(cmd, getSetName(cmd))
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	err := set.Clear(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
