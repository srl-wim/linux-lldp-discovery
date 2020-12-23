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

		nl, err := nlink.NewNLink()
		if err != nil {
			log.Fatal(err)
		}

		go nl.ListAndWatch()

		go nl.TimeoutLoop()

		nl.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
