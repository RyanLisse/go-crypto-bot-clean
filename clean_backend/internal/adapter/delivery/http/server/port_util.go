package server

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

// isPortInUse checks if a port is already in use
func isPortInUse(port int) bool {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

// killProcessOnPort kills any process using the specified port
func killProcessOnPort(port int, logger *zerolog.Logger) error {
	logger.Debug().Int("port", port).Msg("Checking if port is in use")

	if !isPortInUse(port) {
		logger.Debug().Int("port", port).Msg("Port is not in use, no need to kill any process")
		return nil
	}

	logger.Info().Int("port", port).Msg("Port is in use, attempting to kill process")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// For Windows
		findCmd := exec.Command("cmd", "/c", fmt.Sprintf("netstat -ano | findstr :%d", port))
		output, err := findCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to find process using port %d: %w", port, err)
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) == 0 {
			return fmt.Errorf("no process found using port %d", port)
		}

		// Extract PID from the output
		for _, line := range lines {
			if strings.Contains(line, fmt.Sprintf(":%d", port)) {
				parts := strings.Fields(line)
				if len(parts) > 4 {
					pid := parts[len(parts)-1]
					cmd = exec.Command("taskkill", "/F", "/PID", pid)
					break
				}
			}
		}

	case "darwin", "linux":
		// For macOS and Linux
		logger.Debug().Msg("Using lsof to find process on macOS/Linux")
		findCmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
		output, err := findCmd.Output()
		if err != nil {
			// lsof returns error if no process is using the port
			logger.Error().Err(err).Msg("lsof command failed")
			return fmt.Errorf("failed to find process using port %d: %w", port, err)
		}

		outputStr := string(output)
		logger.Debug().Str("lsof_output", outputStr).Msg("lsof command output")

		lines := strings.Split(outputStr, "\n")
		logger.Debug().Int("line_count", len(lines)).Msg("Number of lines in lsof output")

		if len(lines) <= 1 {
			return fmt.Errorf("no process found using port %d", port)
		}

		// Extract PID from the output (skip header line)
		parts := strings.Fields(lines[1])
		logger.Debug().Strs("line_parts", parts).Msg("Parts of the first data line")

		if len(parts) > 1 {
			pid := parts[1]
			logger.Info().Str("pid", pid).Msg("Found process ID to kill")
			cmd = exec.Command("kill", "-9", pid)
		}

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if cmd == nil {
		return fmt.Errorf("failed to create kill command for port %d", port)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill process on port %d: %w", port, err)
	}

	logger.Info().Int("port", port).Msg("Successfully killed process on port")
	return nil
}
