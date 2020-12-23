package nlink

import (
	"time"

	"github.com/vishvananda/netlink"
)

// NetworkAttachements variable
var NetworkAttachements = []string{"bond0.100", "bond0.102", "bond0.103", "ens2f0", "ens2f1"}

// Option struct
type Option func(nl *NLink)

// WithDebug function
func WithDebug(d bool) Option {
	return func(nl *NLink) {
		nl.debug = d
	}
}

// WithTimeout function
func WithTimeout(dur time.Duration) Option {
	return func(nl *NLink) {
		nl.timeout = dur
	}
}

// NLink struct holds the structure
type NLink struct {
	Links map[string]*Link

	timeoutChan chan *cMsg
	workChan    chan *cMsg

	debug   bool
	timeout time.Duration
}

// Link struct
type Link struct {
	Name         string
	MTU          int
	HardwareAddr string
	OperState    string
	Parent       *Link
	Child        *Link
	Vfs          []netlink.VfInfo
}

// NewNLink function defines a new lldptopo
func NewNLink(opts ...Option) (*NLink, error) {
	nl := &NLink{
		Links:       make(map[string]*Link),
		timeoutChan: make(chan *cMsg),
		workChan:    make(chan *cMsg),
	}
	for _, o := range opts {
		o(nl)
	}

	return nl, nil
}
