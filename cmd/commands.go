package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sdslabs/portkey/pkg/connection"
)

var (
	key         string
	sendPath    string
	receive     bool
	receivePath string
)

// rootCmd represents the run command
var rootCmd = &cobra.Command{
	Use:   "portkey",
	Short: "Portkey is a p2p file transfer tool.",
	Long:  `Portkey is a p2p file transfer tool that uses ORTC p2p API over QUIC protocol to achieve very fast file transfer speeds`,
	Run: func(cmd *cobra.Command, args []string) {
		connection.Connect(key, sendPath, receive, receivePath)
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	rootCmd.Flags().StringVarP(&key, "key", "k", "", "Key to connect to peer")
	rootCmd.Flags().StringVarP(&sendPath, "send", "s", "", "Absolute path of directory/file to send")
	rootCmd.Flags().BoolVarP(&receive, "receive", "r", false, "Set to receive files")
	rootCmd.Flags().StringVar(&receivePath, "rpath", "", "Absolute path of where to receive files, pwd by default")
}
