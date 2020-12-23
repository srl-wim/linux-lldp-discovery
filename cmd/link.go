package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/srl-wim/linux-lldp-discovery/pkg/nlink"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:          "link",
	Short:        "get link information",
	Aliases:      []string{"l"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("get link information ...")

		nlink.GetLinks()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
