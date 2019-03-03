package main

import (
	"log"

	"github.com/spf13/cobra"

	"suse.com/caaspctl/pkg/caaspctl"
	"suse.com/caaspctl/pkg/caaspctl/actions/join"
	"suse.com/caaspctl/internal/pkg/caaspctl/deployments/ssh"
)

type JoinOptions struct {
	Role string
}

func newJoinCmd() *cobra.Command {
	joinOptions := JoinOptions{}

	cmd := &cobra.Command{
		Use:   "join <target>",
		Short: "Joins a new node to the cluster",
		Run: func(cmd *cobra.Command, targets []string) {
			user, err := cmd.Flags().GetString("user")
			if err != nil {
				log.Fatalf("Unable to parse user flag: %v", err)
			}
			sudo, err := cmd.Flags().GetBool("sudo")
			if err != nil {
				log.Fatalf("Unable to parse sudo flag: %v", err)
			}

			joinConfiguration := join.JoinConfiguration{}

			switch joinOptions.Role {
			case "master":
				joinConfiguration.Role = caaspctl.MasterRole
			case "worker":
				joinConfiguration.Role = caaspctl.WorkerRole
			default:
				log.Fatalf("Invalid role provided: %q, 'master' or 'worker' are the only accepted roles", joinOptions.Role)
			}

			join.Join(
				joinConfiguration,
				ssh.NewTarget(targets[0], user, sudo),
			)
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().StringVarP(&joinOptions.Role, "role", "", "", "Role that this node will have in the cluster (master|worker)")
	cmd.MarkFlagRequired("role")

	cmd.Flags().StringP("user", "u", "root", "User identity used to connect to target")
	cmd.Flags().Bool("sudo", false, "Run remote command via sudo")

	return cmd
}
