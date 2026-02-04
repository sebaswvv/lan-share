/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package cmd

import (
	"context"

	"github.com/sebaswvv/lan-share/internal/server"

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
		uploadHandler := server.NewUploadHandler()
		srv := setupReceiveServer(uploadHandler)

		// start processing uploads in background with shutdown signal
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go uploadHandler.ProcessUploads(ctx)

		displayServerInfo(localIP, receivePort, "upload")
		runServerWithGracefulShutdown(srv, func() {
			cancel() // signal upload processor to stop
		})
	},
}

func setupReceiveServer(uploadHandler *server.UploadHandler) *server.Server {
	mux := uploadHandler.SetupRoutes()
	return server.New(receivePort, mux)
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	// add port flag
	receiveCmd.Flags().StringVarP(&receivePort, "port", "p", server.DefaultPort, "Port to run the server on")
}
