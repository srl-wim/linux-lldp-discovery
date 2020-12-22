package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootCaCsrTemplate = "/etc/containerlab/templates/ca/csr-root-ca.json"
	certCsrTemplate   = "/etc/containerlab/templates/ca/csr.json"
)

var debug bool
var timeout time.Duration


// path to the topology file
var config string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lldptopo",
	Short: "Discovers the lldp topology",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&config, "input", "i", "", "path to the file with input information for testing")

}