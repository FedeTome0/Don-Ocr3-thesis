package main  

import (
	"fmt"
    "log"
	"net"
	"os"

)

// =============================================================================
// CONFIGURATION HELPERS 
// =============================================================================

// Helper that reads the environment variables with a fallback string
func getEnvironment(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}


// resolveHostnamePortToIP resolves a hostname ("bootstrap_node") to an IPv4 address.
// This is critical for Docker internal networking, as libp2p often requires raw IPs
// for peer discovery rather than DNS hostnames.
func resolveHostnamePortToIP(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil { return "", err }

	if net.ParseIP(host) != nil { return addr, nil } //already an IP

	ips, err := net.LookupIP(host)
	if err != nil { return "", err }
	if len(ips) == 0 { return "", fmt.Errorf("no IP for host %s", host) }
	return net.JoinHostPort(ips[0].String(), port), nil
}

// Helper to set the errors in a safe way
func setJobError(jobIdUint uint64, err error) {
	log.Printf("[Async Task] Job #%d FAILED: %v\n", jobIdUint, err)
	
	JobCache.Lock()
	defer JobCache.Unlock()
	
	if job, exists := JobCache.jobs[jobIdUint]; exists {
		job.State = StateFailed
		job.Err = err
	}
}