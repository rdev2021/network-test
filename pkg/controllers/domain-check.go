package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/rdev2021/network-test/pkg/utils"
)

type DomainCheckRequst struct {
	DomainName string `json:"domainName"`
	DnsServer  string `json:"dnsServer"`
}

type DomainCheckResponse struct {
	DomainName   string   `json:"domainName"`
	DnsServer    string   `json:"dnsServer,omitempty"`
	A            []string `json:"A"`
	CNAME        []string `json:"CNAME"`
	MX           []string `json:"MX"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
	ResponseTime int64    `json:"responseTime"`
	Status       string   `json:"status"`
}

func CheckDomain(domain string, dnsServer string, resolver *net.Resolver) (res DomainCheckResponse) {

	result := DomainCheckResponse{
		DomainName: domain,
		DnsServer:  dnsServer,
	}

	// Resolve A records
	ips, err := resolver.LookupIP(context.Background(), "ip", domain)
	if err == nil {
		for _, ip := range ips {
			result.A = append(result.A, ip.String())
		}
	} else {
		result.A = []string{}
	}

	// Resolve CNAME record
	result.CNAME = []string{}
	cname, err := resolver.LookupCNAME(context.Background(), domain)
	if err == nil {
		if cname != domain+"." {
			result.CNAME = []string{cname}
		}
	}

	// Resolve MX records
	result.MX = []string{}
	mxRecords, err := resolver.LookupMX(context.Background(), domain)
	if err == nil {
		for _, mx := range mxRecords {
			result.MX = append(result.MX, mx.Host)
		}
	}

	result.Status = "Success"
	return result
}

func CheckDomainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)
	start := time.Now()

	var data DomainCheckRequst
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := DomainCheckResponse{
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: fmt.Sprintf("Error decoding JSON: %v", err),
		}

		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)
		return
	}

	portCheckResult := utils.PortCheck(data.DnsServer, 53, "TCP")
	if portCheckResult.Status != "Connected" && data.DnsServer != "" {
		result := DomainCheckResponse{
			ErrorMessage: portCheckResult.ErrorMessage,
			Status:       "Error",
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)
		return
	}

	var resolver *net.Resolver
	if data.DnsServer != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 5 * time.Second,
				}
				return d.DialContext(ctx, "udp", data.DnsServer+":53")
			},
		}
	} else {
		data.DnsServer, _ = utils.GetDefaultDNS()
		resolver = net.DefaultResolver
	}

	response := CheckDomain(data.DomainName, data.DnsServer, resolver)
	response.ResponseTime = time.Since(start).Milliseconds()

	w.Header().Set(ContentType, JsonContentType)
	json.NewEncoder(w).Encode(response)
}
