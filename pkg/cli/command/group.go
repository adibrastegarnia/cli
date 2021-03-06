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
	"github.com/atomix/api/proto/atomix/protocols/log"
	"github.com/atomix/api/proto/atomix/protocols/raft"
	"github.com/atomix/go-client/pkg/client"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

func newGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group {set,get,create,delete}",
		Short: "Manage partition groups and partitions",
		Run:   runGroupGetCommand,
	}
	cmd.PersistentFlags().Duration("timeout", 15*time.Second, "the operation timeout")
	cmd.AddCommand(newGroupSetCommand())
	cmd.AddCommand(newGroupGetCommand())
	cmd.AddCommand(newGroupCreateCommand())
	cmd.AddCommand(newGroupDeleteCommand())
	return cmd
}

func newGroupsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Get a list of partition groups",
		Run:   runGroupsCommand,
	}
	cmd.PersistentFlags().Duration("timeout", 15*time.Second, "the operation timeout")
	cmd.Flags().Bool("no-headers", false, "exclude headers from the output")
	return cmd
}

func printGroups(groups []*client.PartitionGroup, includeHeaders bool) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if includeHeaders {
		fmt.Fprintln(writer, "NAME\tPARTITIONS\tSIZE")
	}
	for _, group := range groups {
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%d\t%d", group.Name, group.Partitions, group.PartitionSize))
	}
	writer.Flush()
}

func printGroup(group *client.PartitionGroup) {
	fmt.Println(fmt.Sprintf("Name:            %s", group.Name))
	fmt.Println(fmt.Sprintf("Namespace:       %s", group.Namespace))
	fmt.Println(fmt.Sprintf("Partitions:      %d", group.Partitions))
	fmt.Println(fmt.Sprintf("Partitions Size: %d", group.PartitionSize))
}

func runGroupsCommand(cmd *cobra.Command, _ []string) {
	client := newClientFromEnv()
	groups, err := client.GetGroups(newTimeoutContext(cmd))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		noHeaders, _ := cmd.Flags().GetBool("no-headers")
		printGroups(groups, !noHeaders)
		ExitWithSuccess()
	}
}

func newGroupSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "set <group>",
		Args: cobra.ExactArgs(1),
		Run:  runGroupSetCommand,
	}
}

func runGroupSetCommand(cmd *cobra.Command, args []string) {
	if err := setClientGroup(args[0]); err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(getClientGroup())
	}
}

func newGroupGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get [group]>",
		Args: cobra.MaximumNArgs(1),
		Run:  runGroupGetCommand,
	}
}

func runGroupGetCommand(cmd *cobra.Command, args []string) {
	var name string
	if len(args) == 0 {
		name = getClientGroup()
	} else {
		name = args[0]
	}

	client := newClientFromGroup(name)
	group, err := client.GetGroup(newTimeoutContext(cmd), getGroupName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		printGroup(group)
		ExitWithSuccess()
	}
}

func newGroupCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create <group>",
		Args: cobra.ExactArgs(1),
		Run:  runGroupCreateCommand,
	}
	cmd.Flags().String("protocol", "raft", "the protocol to run in the partition group")
	cmd.Flags().IntP("partitions", "p", 1, "the number of partitions to create")
	cmd.Flags().IntP("partitionSize", "s", 1, "the size of partitions in the group")
	return cmd
}

func runGroupCreateCommand(cmd *cobra.Command, args []string) {
	name := args[0]
	client := newClientFromGroup(name)

	partitions, _ := cmd.Flags().GetInt("partitions")
	partitionSize, _ := cmd.Flags().GetInt("partitionSize")
	protocolName, _ := cmd.Flags().GetString("protocol")

	var protocol proto.Message
	switch protocolName {
	case "raft":
		protocol = &raft.RaftProtocol{}
	case "log":
		protocol = &log.LogProtocol{}
	}

	group, err := client.CreateGroup(newTimeoutContext(cmd), getGroupName(name), partitions, partitionSize, protocol)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		printGroup(group)
		ExitWithSuccess()
	}
}

func newGroupDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <group>",
		Args: cobra.ExactArgs(1),
		Run:  runGroupDeleteCommand,
	}
}

func runGroupDeleteCommand(cmd *cobra.Command, args []string) {
	name := args[0]
	client := newClientFromGroup(name)
	err := client.DeleteGroup(newTimeoutContext(cmd), getGroupName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithSuccess()
	}
}
