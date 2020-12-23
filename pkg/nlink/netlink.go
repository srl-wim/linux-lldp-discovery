package nlink

import (
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

//GetLinks function
func GetLinks() {
	ll, err := netlink.LinkList()
	if err != nil {
		log.Error(err)
	}
	for _, l := range ll {
		log.Infof("Link Info: %v", l)
	}

}
