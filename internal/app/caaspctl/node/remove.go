/*
 * Copyright (c) 2019 SUSE LLC. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package node

import (
	"github.com/spf13/cobra"

	node "suse.com/caaspctl/pkg/caaspctl/actions/node/remove"
)

func NewRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <node-name>",
		Short: "Removes a node from the cluster",
		Run: func(cmd *cobra.Command, nodenames []string) {
			node.Remove(nodenames[0])
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}
