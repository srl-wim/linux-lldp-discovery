package lldptopo

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ListAndWatch function
func (lt *LldpTopo) ListAndWatch() error {
	for {
		log.Info("ListandWatch message received")
		select {
		case msg := <-lt.workChan:
			switch msg.id {
			case doWork:
				// trigger to see if link status has changed
				d, err := lt.GetLldpTopology()
				if err != nil {
					log.Errorf("Get LLDP topology failure: %s", err)
				}
				if err := lt.ParseLldpDiscovery(d); err != nil {
					log.Errorf("Parse LLDP discovery failed: %s", err)
				}
				fmt.Println("#######################")
				for dName, dev := range lt.Devices {
					fmt.Printf("Device: %s %s %s\n", dName, dev.ID, dev.Kind)
					for eName, ep := range dev.Endpoints {
						fmt.Printf("   Port: %s %s\n", eName, ep.ID)
					}
				}
				fmt.Println("#######################")
			default:
				log.Errorf("Unexpected message: %d", msg)
			}

		}
	}
}

// TimeoutLoop runs timeout loop until the program stops
func (lt *LldpTopo) TimeoutLoop() {
	// Period, in seconds, to dump stats if only counting.
	const TIMEOUT = 5
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(TIMEOUT * time.Second)
		timeout <- true
	}()

	for {
		select {
		case <-timeout:
			log.Infof("Timeout handling")
			go func() {
				time.Sleep(5 * time.Second)
				log.Infof("Timeout done")
				timeout <- true
			}()
			log.Infof("Request compareLinkStatus")
			request := &cMsg{
				id:       doWork,
				respChan: nil,
			}
			lt.workChan <- request
		case msg := <-lt.timeoutChan:
			switch msg.id {
			case shutdown:
				log.Infof("Shutdown message received")
				resp := &cMsg{
					id:       ack,
					respChan: nil,
				}
				msg.respChan <- resp
			default:
				log.Errorf("Unexpected message: %d", msg)
			}
		}
	}
}

// Run the loop to report what is going on, and waiting for interrupt to exit clean.
func (lt *LldpTopo) Run() {
	// Wait for a SIGINT, (typically triggered from CTRL-C), TERM,
	// QUIT. Run cleanup when signal is received. Ideally use os
	// independent os.Interrupt, Kill (but need an exhaustive
	// list.
	sysSignal := make(chan os.Signal, 1)
	signal.Notify(sysSignal,
		syscall.SIGTERM, // ^C
		syscall.SIGINT,  // kill
		syscall.SIGQUIT, // QUIT
		syscall.SIGABRT)
	doneSignal := make(chan bool)
	log.Info("watching for shutdown...")
	go func() {
		for range sysSignal {
			log.Info("Interrupt, stopping gracefully...")
			// add a stop function to stop the go-routine

			shut(lt.timeoutChan)

			// Now that they are all done. Unblock
			doneSignal <- true

		}
	}()
	//
	// Block here waiting for cleanup. This is likely to be in a
	// main select along other possible conditions (like a timeout
	// to update stats?)
	<-doneSignal

	log.Info(" Bye...")
}

func shut(cChan chan *cMsg) {
	respChan := make(chan *cMsg)
	request := &cMsg{
		id:       shutdown,
		respChan: respChan,
	}
	// Send shutdown message
	cChan <- request
	// Wait for ACK
	<-respChan
	close(cChan)
}
