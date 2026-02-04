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
	"strings"
)

// server represents the HTTP server for file sharing
type Server struct {
	httpServer *http.Server
	port       string
}

// new creates a new server instance
func New(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: MaxHeaderBytes,
			// no ReadTimeout/WriteTimeout for large file transfers
		},
		port: port,
	}
}

// start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting server on port %s", s.port)
	return s.httpServer.ListenAndServe()
}

// shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// getLocalIP returns the local IP address of the machine
func GetLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// prioritize physical network interfaces (especially for Mac OS)
	priorityInterfaces := []string{"en0", "en1", "eth0", "wlan0"}

	// first pass: try priority interfaces
	for _, priorityName := range priorityInterfaces {
		for _, iface := range ifaces {
			if iface.Name != priorityName {
				continue
			}

			// skip loopback and down interfaces
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}

			if ip := getValidIPFromInterface(iface); ip != "" {
				return ip, nil
			}
		}
	}

	// second pass: any valid interface
	for _, iface := range ifaces {
		// skip loopback and down interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if ip := getValidIPFromInterface(iface); ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("no local IP address found")
}

// getValidIPFromInterface extracts a valid IPv4 address from a network interface
func getValidIPFromInterface(iface net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipv4 := ipNet.IP.To4(); ipv4 != nil {
				ip := ipv4.String()
				// skip link-local addresses (169.254.x.x)
				if strings.HasPrefix(ip, "169.254") {
					continue
				}
				// skip common virtual adapter ranges used by VirtualBox and VMware
				// 192.168.176.x - VMware NAT
				// 192.168.224.x - VirtualBox host-only adapter
				if strings.HasPrefix(ip, "192.168.176") || strings.HasPrefix(ip, "192.168.224") {
					continue
				}
				// found a good IP
				return ip
			}
		}
	}

	return ""
}
