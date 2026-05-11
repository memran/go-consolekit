package console

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const daemonEnvVar = "_CONSOLEKIT_DAEMON_CHILD"

func IsDaemonChild() bool {
	return os.Getenv(daemonEnvVar) == "1"
}

func hasDaemonFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--daemon" {
			return true
		}
	}
	return false
}

func stripDaemonFlag(args []string) []string {
	result := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == "--daemon" {
			continue
		}
		result = append(result, args[i])
	}
	return result
}

func extractFlagValue(args []string, name string) string {
	for i, arg := range args {
		if arg == name && i+1 < len(args) {
			return args[i+1]
		}
		prefix := name + "="
		if strings.HasPrefix(arg, prefix) {
			return strings.TrimPrefix(arg, prefix)
		}
	}
	return ""
}

func writePID(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644)
}

func readPID(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func removePID(path string) error {
	return os.Remove(path)
}

func redirectOutput(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	os.Stdout = f
	os.Stderr = f
	return nil
}


