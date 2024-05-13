package utils

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rdev2021/network-test/pkg/models"
)

func PortCheck(host string, port int, checkType string) *models.PortCheckResponse {
	portStr := strconv.Itoa(port)
	start := time.Now()

	if checkType != "TCP" && checkType != "UDP" {
		responseTime := time.Since(start).Milliseconds()
		result := &models.PortCheckResponse{
			Host:         host,
			Port:         port,
			Type:         checkType,
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: "Invalid check type: " + checkType,
		}
		return result
	}

	conn, err := net.DialTimeout(strings.ToLower(checkType), net.JoinHostPort(host, portStr), 2*time.Second)
	responseTime := time.Since(start).Milliseconds()

	result := &models.PortCheckResponse{
		Host:         host,
		Port:         port,
		Type:         checkType,
		ResponseTime: responseTime,
	}

	if err != nil {
		result.Status = "Error"
		result.ErrorMessage = err.Error()
	} else {
		defer conn.Close()
		result.Status = "Connected"
	}

	return result
}

func GetDefaultDNS() (string, error) {
	// On Unix-like systems, the default DNS server can be read from /etc/resolv.conf
	// On Windows, the GetSystemDefaultDNS method can be used
	if os.PathSeparator == '/' {
		// Unix-like system
		file, err := os.Open("/etc/resolv.conf")
		if err != nil {
			return "", err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "nameserver") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					return parts[1], nil
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	} else {
		// Windows system
		addrs, err := net.LookupHost("")
		if err != nil {
			return "", err
		}
		return addrs[0], nil
	}
	return "", nil
}
