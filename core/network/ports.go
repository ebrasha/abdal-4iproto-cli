/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : ports.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Port availability checks and automatic suggestions
 *                for server and panel listeners before any download.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package network

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"abdal-4iproto-cli/core/config"
)

// PortCheckError is returned when a port is already reserved/in use.
type PortCheckError struct {
	Port int
}

func (e *PortCheckError) Error() string {
	return fmt.Sprintf("port %d is reserved or already in use on this host", e.Port)
}

// IsPortAvailable returns true when the TCP port can be bound locally.
func IsPortAvailable(port int) bool {
	if port < config.MinPort || port > config.MaxPort {
		return false
	}
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// ValidatePorts ensures every port in the slice is free. Returns the first
// offending port and a PortCheckError when one is busy.
func ValidatePorts(ports []int) error {
	for _, p := range ports {
		if !IsPortAvailable(p) {
			return &PortCheckError{Port: p}
		}
	}
	return nil
}

// SuggestFreePorts returns up to `count` random-looking ports in the
// configured suggestion range that are currently available.
func SuggestFreePorts(count int) ([]int, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}
	var found []int
	// Deterministic scan with a pseudo-random stride to avoid well-known ports.
	stride := 137
	for candidate := config.PortSuggestionMin; candidate <= config.PortSuggestionMax && len(found) < count; candidate += stride {
		if IsPortAvailable(candidate) {
			found = append(found, candidate)
		}
	}
	if len(found) < count {
		// Fallback: linear scan.
		for candidate := config.PortSuggestionMin; candidate <= config.PortSuggestionMax && len(found) < count; candidate++ {
			if IsPortAvailable(candidate) {
				dup := false
				for _, f := range found {
					if f == candidate {
						dup = true
						break
					}
				}
				if !dup {
					found = append(found, candidate)
				}
			}
		}
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("could not find any free ports in range %d-%d", config.PortSuggestionMin, config.PortSuggestionMax)
	}
	return found, nil
}

// ParsePortList converts a comma-separated list of ports into integers.
func ParsePortList(raw string) ([]int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty port list")
	}
	parts := strings.Split(raw, ",")
	var ports []int
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid port '%s': %w", part, err)
		}
		if n < config.MinPort || n > config.MaxPort {
			return nil, fmt.Errorf("port %d out of range", n)
		}
		ports = append(ports, n)
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("no valid ports parsed")
	}
	return ports, nil
}

// FormatPortList renders ports as a comma-separated string.
func FormatPortList(ports []int) string {
	var parts []string
	for _, p := range ports {
		parts = append(parts, strconv.Itoa(p))
	}
	return strings.Join(parts, ",")
}

// WaitForPortRelease polls until a port becomes free or the timeout elapses.
func WaitForPortRelease(port int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if IsPortAvailable(port) {
			return true
		}
		time.Sleep(250 * time.Millisecond)
	}
	return IsPortAvailable(port)
}

// ReservedSystemPorts returns a small set of ports that should never be
// suggested automatically (well-known services).
func ReservedSystemPorts() map[int]bool {
	common := []int{22, 80, 443, 3389, 52202}
	if runtime.GOOS == "windows" {
		common = append(common, 135, 139, 445)
	}
	m := make(map[int]bool, len(common))
	for _, p := range common {
		m[p] = true
	}
	return m
}

// IsWellKnownReserved returns true for ports we should not auto-suggest.
func IsWellKnownReserved(port int) bool {
	return ReservedSystemPorts()[port]
}
