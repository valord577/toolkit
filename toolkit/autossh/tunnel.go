package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net"
	"sync"
	"time"

	"toolkit/logs"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

const defaultFlowBufferSize = 16 * 1024

var (
	shutdown  bool
	waitgroup sync.WaitGroup
)

func Shutdown() {
	shutdown = true
	waitgroup.Wait()
}

func Startup() error {
	tuns := tunnels()
	if len(tuns) < 1 {
		return errors.New("empty ssh config or tunnels")
	}

	for _, tun := range tuns {
		waitgroup.Add(1)
		go startup(tun)
	}
	return nil
}

func startup(tun *tunnel) {
	defer waitgroup.Done()

	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		logs.Errorf("read crypto/rand, err: %s", err.Error())
		return
	}
	uuid := hex.EncodeToString(bs)

	l := logs.With(
		zap.String("uuid", uuid),
		zap.String("service", tun.Service),
		zap.String("listenOn", tun.ListenOn.string()),
	)
	l = l.WithOpts(zap.AddCallerSkip(-1))
	l.Infof("%s", "starting ssh tunnel service")

	for !shutdown {
		// restart when not shutdown
		forwarding(tun, l)
		time.Sleep(1 * time.Second)
	}
}

func forwarding(tun *tunnel, log *logs.Logger) {
	listenOn := tun.ListenOn
	listenAt := tun.ListenAt
	forwardTo := tun.ForwardTo
	sshConfig := tun.SshConfig

	transport := func(src net.Conn, sshClient *ssh.Client) {
		var (
			dst net.Conn
			err error
		)
		switch listenOn {
		case listenOnLocal:
			dst, err = sshClient.Dial("tcp", forwardTo)
		case listenOnRemote:
			dst, err = net.DialTimeout("tcp", forwardTo, 10*time.Second)
		}
		if err != nil {
			log.Errorf("target dial, err: %s", err.Error())
			return
		}
		go exflow(dst, src, log)
		go exflow(src, dst, log)
	}

	log.Infof("ssh dial, alias: %s, address: %s", sshConfig.Alias, sshConfig.Address)
	sshClient, err := ssh.Dial("tcp", sshConfig.Address, sshConfig.Config)
	if err != nil {
		log.Errorf("ssh dial, err: %s", err.Error())
		return
	}
	defer sshClient.Close()

	log.Infof("listen, address: %s", listenAt)
	var listener net.Listener
	switch listenOn {
	case listenOnLocal:
		listener, err = net.Listen("tcp", listenAt)
	case listenOnRemote:
		listener, err = sshClient.Listen("tcp", listenAt)
	}
	if err != nil {
		log.Errorf("listen, err: %s", err.Error())
		return
	}
	defer listener.Close()

	done := make(chan struct{}, 1)
	hook := true
	go func() {
		for hook {
			time.Sleep(1 * time.Second)
			if shutdown {
				done <- struct{}{}
				break
			}
		}
	}()
	go func() {
		for {
			conn, e := listener.Accept()
			if e != nil {
				log.Errorf("accept, err: %s", e.Error())
				done <- struct{}{}
				break
			}
			log.Infof("accept: '%s'", conn.RemoteAddr().String())
			go transport(conn, sshClient)
		}
	}()
	<-done
	hook = false
}

func exflow(dst, src net.Conn, log *logs.Logger) {
	defer func() {
		_ = dst.Close()
		_ = src.Close()
	}()

	var (
		buff = make([]byte, defaultFlowBufferSize)
		addr = src.RemoteAddr().String()
	)
	for !shutdown {
		nr, err := src.Read(buff)
		if err != nil {
			log.Errorf("%s", err.Error())
			break
		}
		if nr > 0 {
			nw, ew := dst.Write(buff[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			log.Debugf("exchange flow, addr: '%s', written: %d Bytes", addr, nw)
			if ew != nil {
				err = ew
				log.Errorf("%s", err.Error())
				break
			}
			if nr != nw {
				err = errors.New("short write")
				log.Errorf("%s", err.Error())
				break
			}
		}
	}
}
