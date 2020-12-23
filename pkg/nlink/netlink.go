package nlink

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// ParseLinks function
func (nl *NLink) ParseLinks() error {
	ll, err := netlink.LinkList()
	if err != nil {
		return err
	}
	for _, l := range ll {
		attr := l.Attrs()
		link := new(Link)
		link.Name = attr.Name
		link.OperState = attr.OperState.String()
		link.MTU = attr.MTU
		link.HardwareAddr = attr.HardwareAddr.String()
		link.Vfs = attr.Vfs

		// only validate link changes for interfaces that are relevant, meaning
		// configured through Network Attachements
		if linkInNetworkAtatchement(&link.Name) {
			nl.checkLinkDelta(&link.Name, link)
		}
	}
	fmt.Println("#######################")
	for ifName, l := range nl.Links {
		fmt.Printf("Link: %s %s %s &s %d\n", ifName, l.OperState, l.HardwareAddr, l.MTU)
		for i, vf := range l.Vfs {
			fmt.Printf("   VFs: %d %d\n", i, vf.Vlan)
		}
	}
	fmt.Println("#######################")
	return nil
}

func (nl *NLink) checkLinkDelta(ifName *string, link *Link) (c bool) {
	if l, ok := nl.Links[*ifName]; ok {
		if l.MTU != link.MTU {
			c = true
			log.Infof("Link %s, MTU changed from %d -> %d", *ifName, l.MTU, link.MTU)
		}
		if l.HardwareAddr != link.HardwareAddr {
			c = true
			log.Infof("Link %s, HardwareAddr changed from %s -> %s", *ifName, l.HardwareAddr, link.HardwareAddr)
		}
		if l.OperState != link.OperState {
			c = true
			log.Infof("Link %s, OperState changed from %s -> %s", *ifName, l.OperState, link.OperState)
		}
		// TODO check if VF(s) are added or deleted
		for i, vf := range link.Vfs {
			if l.Vfs[i].Vlan != vf.Vlan && vf.Vlan == 0 {
				c = true
				log.Infof("Link %s, VLAN deleted %d -> %d", *ifName, l.Vfs[i].Vlan, vf.Vlan)
			}
			if l.Vfs[i].Vlan != vf.Vlan && vf.Vlan != 0 {
				c = true
				log.Infof("Link %s, VLAN changed or added %d -> %d", *ifName, l.Vfs[i].Vlan, vf.Vlan)
			}
		}

	} else {
		c = true
		log.Infof("New Link detected %s", *ifName)
		nl.Links[*ifName] = link
		for _, vf := range link.Vfs {
			if vf.Vlan != 0 {
				log.Infof("Link %s, VLAN to be provisioned %d -> %d", *ifName, 0, vf.Vlan)
			}
		}
	}
	return c
}

func linkInNetworkAtatchement(ifName *string) bool {
	for _, name := range NetworkAttachements {
		if name == *ifName {
			return true
		}
	}
	return false

}


