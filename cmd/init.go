package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/srl-wim/linux-lldp-discovery/pkg/lldptopo"
)

// deployCmd represents the deploy command
var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "initialize the lldp topology discovery tool",
	Aliases:      []string{"dep"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("deploying nuage aws tgw network manager configuration ...")
		opts := []lldptopo.Option{
			lldptopo.WithDebug(debug),
			lldptopo.WithTimeout(timeout),
			lldptopo.WithInputFile(config),
		}

		lt, err := lldptopo.NewLldpTopo(opts...)
		if err != nil {
			log.Fatal(err)
		}

		if err := lt.CheckLldpDaemon(); err != nil {
			log.Fatal(err)
		}

		go lt.ListAndWatch()

		go lt.TimeoutLoop()

		lt.Run()


		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
