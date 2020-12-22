package lldptopo

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// LldpTopo struct holds the structure of the
type LldpTopo struct {
	InputFile   *string
	Devices     map[string]*Device
	Connections map[int]*Connection

	ctx     context.Context
	debug   bool
	timeout time.Duration
}

// Device is a struct that contains the information of a device element
type Device struct {
	Name      string
	DeviceID  *string
	DeviceARN *string
	Index     int
	Kind      string
	Model     string
	Serial    string
	Vendor    string
	Region    string
	Endpoints map[string]*Endpoint
}

// Connection is a struct that contains the information of a link between 2 containers
type Connection struct {
	A      *Endpoint
	B      *Endpoint
	Labels map[string]string
}

// Endpoint is a struct that contains information of a link endpoint
type Endpoint struct {
	Device   *Device
	Name     string
	LinkID   *string
	LinkARN  *string
	Provider string
	Kind     string
}

// Option struct
type Option func(lt *LldpTopo)

// WithDebug function
func WithDebug(d bool) Option {
	return func(lt *LldpTopo) {
		lt.debug = d
	}
}

// WithTimeout function
func WithTimeout(dur time.Duration) Option {
	return func(lt *LldpTopo) {
		lt.timeout = dur
	}
}

// WithInputFile function
func WithInputFile(file string) Option {
	return func(lt *LldpTopo) {
		if file == "" {
			return
		}
		log.Info(file)
		if err := lt.ParseInputFile(&file); err != nil {
			log.Fatalf("failed to read topology file: %v", err)
		}
	}
}

// NewLldpTopo function defines a new lldptopo
func NewLldpTopo(opts ...Option) (*LldpTopo, error) {
	lt := &LldpTopo{
		InputFile: new(string),
		ctx:       context.Background(),
	}
	for _, o := range opts {
		o(lt)
	}

	return lt, nil
}
