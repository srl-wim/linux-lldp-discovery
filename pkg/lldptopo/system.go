package lldptopo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// CheckLldpDaemon function
func (lt *LldpTopo) CheckLldpDaemon() error {
	cmd := exec.Command("systemctl", "check", "lldpd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("systemctl finished with non-zero: %v, status: %s", exitErr, string(out))
		}
		log.Fatalf("failed to run systemctl: %v", err)
		os.Exit(1)

	}
	log.Infof("Status is: %s", string(out))
	return nil
}

// GetLldpTopology function
func (lt *LldpTopo) GetLldpTopology() (*Discovery, error) {
	cmd := exec.Command("lldpcli", "show", "neighbors", "-f", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("lldpcli finished with non-zero: %v", exitErr)
		}
		log.Fatalf("failed to run lldpcli: %v", err)
		os.Exit(1)

	}
	//fmt.Printf("Status is: %s\n", string(out))

	var d Discovery
	err = json.Unmarshal(out, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
