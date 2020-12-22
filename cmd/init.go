package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/srl-wim/linux-lldp-discovery/lldptopo"
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

		lt.CheckLldpDaemon()

		lt.GetLldpTopology()

		/*
			for i := 0; i < 10; i++ {
				d, err := lt.ParseInputFile(lt.InputFile)
				if err != nil {
					log.Fatal(err)
				}
				if err := lt.ParseLldpDiscovery(d); err != nil {
					log.Fatal(err)
				}

				fmt.Println("#######################")
				for dName, dev := range lt.Devices {
					fmt.Printf("Device: %s %s %s\n", dName, dev.ID, dev.Kind)
					for eName, ep := range dev.Endpoints {
						fmt.Printf("   Port: %s %s\n", eName, ep.ID)
					}
				}
				fmt.Println("#######################")
				time.Sleep(10 * time.Second)

			}
		*/
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
