package snapctl

import (
	"fmt"
	"os/exec"
	"strings"
)

// Get reads a config option
func Get(key string) (string, error) {
	output, err := exec.Command("snapctl", "get", key).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("snapctl get error for %s: %s: %s", key, err, output)
	}
	return strings.TrimSpace(string(output)), nil
}

// TODO
// get -d doc
// get -t strict
// get multiple

// Set writes a config option
func Set(key string, val string) error {
	output, err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, val)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl set error for %s: %s: %s", key, err, output)
	}
	return nil
}

// UnsetConfig uses snapctl to unset a config value from a key
func Unset(key string) error {
	output, err := exec.Command("snapctl", "unset", key).CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl unset error for %s: %s: %s", key, err, output)
	}
	return nil
}

// Start start one or more services and optionally enable all
func Start(enable bool, services ...string) error {
	args := []string{"start"}
	if enable {
		args = append(args, "--enable")
	}
	for _, s := range services {
		args = append(args, snapName+"."+s)
	}

	output, err := exec.Command("snapctl", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl start error: %s: %s", err, output)
	}

	return nil
}

// Stop uses snapctl to stop one or more services and optionally disable all
func Stop(disable bool, services ...string) error {
	args := []string{"stop"}
	if disable {
		args = append(args, "--disable")
	}
	for _, s := range services {
		args = append(args, snapName+"."+s)
	}

	output, err := exec.Command("snapctl", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl stop error: %s: %s", err, output)
	}

	return nil
}
