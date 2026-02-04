/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

// Server represents the HTTP server for file sharing
type Server struct {
	httpServer *http.Server
	port       string
}

// New creates a new server instance
func New(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: MaxHeaderBytes,
			// No ReadTimeout/WriteTimeout for large file transfers
		},
		port: port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting server on port %s", s.port)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// GetLocalIP returns the local IP address of the machine
func GetLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		// skip loopback and down interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipv4 := ipNet.IP.To4(); ipv4 != nil {
					ip := ipv4.String()
					// skip link-local addresses (169.254.x.x)
					// prefer 192.168.x.x or 10.x.x.x ranges
					if ip[:7] == "169.254" {
						continue
					}
					// skip common virtual adapter ranges
					if ip[:11] == "192.168.176" || ip[:11] == "192.168.224" {
						continue
					}
					// found a good IP
					return ip, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no local IP address found")
}
