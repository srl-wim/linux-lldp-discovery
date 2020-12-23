package lldptopo

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// LldpTopo struct holds the structure of the
type LldpTopo struct {
	InputFile *string
	Devices   map[string]*Device

	timeoutChan chan *cMsg
	workChan    chan *cMsg

	ctx     context.Context
	debug   bool
	timeout time.Duration
}

// Device is a struct that contains the information of a device element
type Device struct {
	Kind      string
	ID        string
	IDType    string
	Endpoints map[string]*Endpoint
}

// Endpoint is a struct that contains information of a link endpoint
type Endpoint struct {
	Device *Device
	ID     string
	IDType string
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
		lt.InputFile = &file
		log.Info(file)
		//if err := lt.ParseInputFile(&file); err != nil {
		//	log.Fatalf("failed to read topology file: %v", err)
		//}
	}
}

// NewLldpTopo function defines a new lldptopo
func NewLldpTopo(opts ...Option) (*LldpTopo, error) {
	lt := &LldpTopo{
		InputFile:   new(string),
		Devices:     make(map[string]*Device),
		timeoutChan: make(chan *cMsg),
		workChan:    make(chan *cMsg),
		ctx:         context.Background(),
	}
	for _, o := range opts {
		o(lt)
	}

	return lt, nil
}
