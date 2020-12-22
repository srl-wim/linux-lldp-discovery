package lldptopo

import (
	"fmt"
	"os"
	"os/exec"
)

// CheckLldpDaemon function
func (lt *LldpTopo) CheckLldpDaemon() error {
	cmd := exec.Command("systemctl", "check", "lldpd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("systemctl finished with non-zero: %v\n", exitErr)
		} else {
			fmt.Printf("failed to run systemctl: %v", err)
			os.Exit(1)
		}
	}
	fmt.Printf("Status is: %s\n", string(out))
	return nil
}

// GetLldpTopology function
func (lt *LldpTopo) GetLldpTopology() error {
	cmd := exec.Command("lldpcli", "show", "neighbors", "-f", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("lldpcli finished with non-zero: %v\n", exitErr)
		} else {
			fmt.Printf("failed to run lldpcli: %v", err)
			os.Exit(1)
		}
	}
	fmt.Printf("Status is: %s\n", string(out))
	return nil
}