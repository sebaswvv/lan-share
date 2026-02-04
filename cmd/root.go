/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lanshare",
	Short: "Share files quickly over a local network using a temporary link or QR code",
	Long: `lanshare is a lightweight CLI tool for sharing files directly over a local network (LAN).

It starts a temporary local server on your machine and generates a one-time download link
and QR code that can be opened from another laptop or phone on the same network.

No cloud storage, no accounts, no internet required.
The file is streamed directly from the sender to the receiver and the session
automatically expires or stops after download.`,
}

// execute adds all child commands to the root command and sets flags appropriately.
// this is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
