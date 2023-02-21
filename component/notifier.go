package component

import (
	"time"

	"toolkit/config"
	log "toolkit/logger"
	"toolkit/netdev"
)

type Notifier struct {
	done chan struct{}
	stop bool
}

func (n *Notifier) init() error {
	if n.done == nil {
		n.done = make(chan struct{}, 1)
	}

	go func() {
		for !n.stop {
			log.Infof("trigger notify auto")
			netdev.NotifyChanges(config.Device(), config.Receiver())
			if !n.stop {
				time.Sleep(config.Period())
			}
		}
		n.done <- struct{}{}
	}()
	return nil
}

func (n *Notifier) free() error {
	n.stop = true
	<-n.done
	return nil
}
