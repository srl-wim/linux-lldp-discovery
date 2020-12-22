package lldptopo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Discovery struct
type Discovery struct {
	Lldp Lldp `json:"lldp,omitempty"`
}

// Lldp struct
type Lldp struct {
	Itfce map[string]Itfce `json:"interface"`
}

// Itfce struct
type Itfce struct {
	Via     string             `json:"via"`
	Rid     string             `json:"rid"`
	Age     string             `json:"age"`
	Chassis map[string]Chassis `json:"chassis"`
	Port    Port               `json:"port"`
}

// Chassis struct
type Chassis struct {
	ID         ID           `json:"id"`
	Descr      string       `json:"descr"`
	MgmtIP     string       `json:"mgmt-ip"`
	Capability []Capability `json:"capability"`
}

// Capability struct
type Capability struct {
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

// Port struct
type Port struct {
	ID    ID     `json:"id"`
	Descr string `json:"descr"`
	TTL   string `json:"ttl"`
}

// ID struct
type ID struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ParseInputFile parses the input topology
func (lt *LldpTopo) ParseInputFile(f *string) error {
	jsonFile, err := ioutil.ReadFile(*f)
	if err != nil {
		return err
	}
	log.Infof(fmt.Sprintf("Topology file contents:\n%s\n", jsonFile))

	var d Discovery
	err = json.Unmarshal(jsonFile, &d)
	if err != nil {
		return err
	}

	for n, i := range d.Lldp.Itfce {
		fmt.Printf("Interface name: %s, Info: %v", n, i)
	}
	fmt.Println(d.Lldp.Itfce)
	return nil
}
