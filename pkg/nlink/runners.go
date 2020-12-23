package nlink

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ListAndWatch function
func (nl *NLink) ListAndWatch() error {
	for {
		log.Info("ListandWatch message received")
		select {
		case msg := <-nl.linkChan:
			switch msg.id {
			case compareLinkStatus:
				// trigger to see if link status has changed
				log.Info("ListandWatch compareLinkStatus message received")
				if err := nl.ParseLinks(); err != nil {
					log.Error(err)
				}
			default:
				log.Errorf("Unexpected message: %d", msg)
			}

		}
	}
}

// TimeoutLoop runs timeout loop until the program stops
func (nl *NLink) TimeoutLoop() {
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
				id:       compareLinkStatus,
				respChan: nil,
			}
			nl.linkChan <- request
		case msg := <-nl.timeoutChan:
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
func (nl *NLink) Run() {
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

			shut(nl.timeoutChan)

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
