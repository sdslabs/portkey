package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sdslabs/portkey/pkg/connection"
	"github.com/sdslabs/portkey/pkg/utils"
)

var (
	key            string
	sendPath       string
	receive        bool
	receivePath    string
	certPath       string
	privateKeyPath string
	doBenchmarking bool
)

// rootCmd represents the run command
var rootCmd = &cobra.Command{
	Use:   "portkey",
	Short: "Portkey is a p2p file transfer tool.",
	Long:  `Portkey is a p2p file transfer tool that uses ORTC p2p API over QUIC protocol to achieve very fast file transfer speeds`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.Validate(key, sendPath, receive, receivePath, certPath, privateKeyPath, doBenchmarking)
	},
	Run: func(cmd *cobra.Command, args []string) {
		connection.Connect(key, sendPath, receive, receivePath, certPath, privateKeyPath, doBenchmarking)
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	rootCmd.Flags().StringVarP(&key, "key", "k", "", "Key to connect to peer")
	rootCmd.Flags().StringVarP(&sendPath, "send", "s", "", "Absolute path of directory/file to send")
	rootCmd.Flags().BoolVarP(&receive, "receive", "r", false, "Set to receive files")
	rootCmd.Flags().StringVarP(&receivePath, "dir", "d", "", "Absolute path of where to receive files, pwd by default")
	rootCmd.Flags().StringVarP(&certPath, "cert", "c", "", "Absolute path of cert.pem (x509 certificate in PEM format)")
	rootCmd.Flags().StringVarP(&privateKeyPath, "pkey", "p", "", "Absolute path of key.pem (Private key corresponding to certificate in PEM format)")
	rootCmd.Flags().BoolVarP(&doBenchmarking, "benchmark", "b", false, "Set to benchmark locally(for local testing)")
}
