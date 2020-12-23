package lldptopo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Discovery struct
type Discovery struct {
	Lldp Lldp `json:"lldp,omitempty"`
}

// Lldp struct
type Lldp struct {
	Itfce []map[string]Itfce `json:"interface,omitempty"`
}

// Itfce struct
type Itfce struct {
	Via     string             `json:"via,omitempty"`
	Rid     string             `json:"rid,omitempty"`
	Age     string             `json:"age,omitempty"`
	Chassis map[string]Chassis `json:"chassis,omitempty"`
	Port    Port               `json:"port,omitempty"`
}

// Chassis struct
type Chassis struct {
	ID    ID     `json:"id,omitempty"`
	Descr string `json:"descr,omitempty"`
	//MgmtIP     []string     `json:"mgmt-ip,omitempty"`
	Capability []Capability `json:"capability,omitempty"`
}

// Capability struct
type Capability struct {
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// Port struct
type Port struct {
	ID    ID     `json:"id,omitempty"`
	Descr string `json:"descr,omitempty"`
	TTL   string `json:"ttl,omitempty"`
}

// ID struct
type ID struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

// ParseInputFile parses the input topology
func (lt *LldpTopo) ParseInputFile(f *string) (*Discovery, error) {
	jsonFile, err := ioutil.ReadFile(*f)
	if err != nil {
		return nil, err
	}
	log.Debugf(fmt.Sprintf("Topology file contents:\n%s\n", jsonFile))

	var d Discovery
	err = json.Unmarshal(jsonFile, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ParseLldpDiscovery parses the lldp input and reports on the delta
func (lt *LldpTopo) ParseLldpDiscovery(d *Discovery) error {
	var err error
	//fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@")
	for _, itfce := range d.Lldp.Itfce {
		for ifName, i := range itfce {
			for dName, c := range i.Chassis {
				d := new(Device)
				e := new(Endpoint)
				d.Kind, err = findDeviceKind(&c.Descr)
				d.ID = c.ID.Value
				d.IDType = c.ID.Type
				e.ID = i.Port.ID.Value
				e.IDType = i.Port.ID.Type
				e.Device = d
				//fmt.Printf("Topo: %s, %s\n", dName, ifName)
				lt.checkTopoDelta(&dName, &ifName, d, e)
			}
		}
	}
	//fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@")
	return err
}

func findDeviceKind(d *string) (string, error) {
	if strings.Contains(*d, "TiMOS") && strings.Contains(*d, "NUAGE") {
		return "wbx", nil
	}
	if strings.Contains(*d, "SRLinux") {
		return "srl", nil
	}
	if strings.Contains(*d, "Ubuntu") {
		return "linux", nil
	}
	if strings.Contains(*d, "Centos") {
		return "linux", nil
	}
	if strings.Contains(*d, "Redhat") {
		return "linux", nil
	}
	return "linux", nil
}

func (lt *LldpTopo) checkTopoDelta(dName, ifName *string, d *Device, e *Endpoint) (c bool) {
	if dev, ok := lt.Devices[*dName]; ok {
		// device already exists, check differences
		if dev.Kind != d.Kind {
			c = true
			log.Infof("Kind changed from %s to %s on device %s", dev.Kind, d.Kind, *dName)
		}
		if dev.ID != d.ID {
			c = true
			log.Infof("ID changed from %s to %s on device %s", dev.ID, d.ID, *dName)
		}
		if ep, ok := dev.Endpoints[*ifName]; ok {
			if ep.ID != e.ID {
				c = true
				log.Infof("Port Mac changed from %s to %s on device %s", ep.ID, e.ID, *dName)
				ep.ID = e.ID
			}
		} else {
			// new interface on existing device
			c = true
			dev.Endpoints[*ifName] = e
			log.Infof("New Interface %s detected on exisitng device %s", *ifName, *dName)

		}

	} else {
		// new device discovered
		c = true
		log.Infof("New Device %s detected with interface %s", *dName, *ifName)
		lt.Devices[*dName] = d
		lt.Devices[*dName].Endpoints = make(map[string]*Endpoint)
		lt.Devices[*dName].Endpoints[*ifName] = e
	}
	return c
}
