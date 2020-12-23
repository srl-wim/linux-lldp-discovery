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
		attr := l.Attrs()
		log.Infof("Link Name: %s", attr.Name)
		log.Infof("  Link Type: %v", l.Type())
		log.Infof("  Link Attr: %v", attr)
		for _, v := range attr.Vfs {
			log.Infof("    VFs: %v", v)

		}
	}

}
