/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package cmd

import (
	"fmt"
	"lanshare/internal/server"
	"os"

	"github.com/fatih/color"
	qrterminal "github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

var receivePort string

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive files from other devices on your network",
	Long:  `Start a server that allows other devices to upload files to your computer.`,
	Run: func(cmd *cobra.Command, args []string) {
		localIP := getLocalIP()
		srv := setupReceiveServer()

		displayReceiveInfo(localIP)
		runServer(srv)
	},
}

func setupReceiveServer() *server.Server {
	uploadHandler := server.NewUploadHandler()

	// start processing uploads in background
	go uploadHandler.ProcessUploads()

	mux := uploadHandler.SetupRoutes()
	return server.New(receivePort, mux)
}

func displayReceiveInfo(localIP string) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	yellow := color.New(color.FgYellow)

	url := fmt.Sprintf("http://%s:%s", localIP, receivePort)

	fmt.Println()
	green.Println("✓ Server started successfully!")
	fmt.Println()
	magenta.Println("📱 Scan QR code to upload files:")
	fmt.Println()
	qrterminal.GenerateHalfBlock(url, qrterminal.L, os.Stdout)
	fmt.Println()
	magenta.Print("🌐  URL: ")
	cyan.Println(url)
	fmt.Println()
	yellow.Println("📥 Waiting for uploads... Press Ctrl+C to stop")
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	// add port flag
	receiveCmd.Flags().StringVarP(&receivePort, "port", "p", "8080", "Port to run the server on")
}
