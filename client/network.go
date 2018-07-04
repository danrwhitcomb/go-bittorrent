package client

import (
	"errors"
	"net"
)

type BittorentNetwork struct {
	listener *net.TCPListener
}

// BindToPort begins listening on a port and
// returns the port it selected
func (c *BittorentNetwork) BindToPort(callback func(Peer)) (int, error) {

	port := -1
	for i := 6881; i < 6890; i++ {
		log.Debug("Attempting to listen on %d", i)
		addr := net.TCPAddr{Port: i}
		listen, err := net.ListenTCP("tcp", &addr)
		if err == nil {
			c.listener = listen
			port = i
			break
		}
	}

	if c.listener == nil {
		return 0, errors.New("Unable to bind to port between 6881 and 6889")
	}

	log.Infof("Listening on port %d", port)
	go c.listenOnPort(callback)

	return port, nil
}

func (c *BittorentNetwork) listenOnPort(callback func(Peer)) {
	for {
		conn, err := c.listener.AcceptTCP()
		if err != nil {
			log.Error("Connection to client failed: ", err)
			continue
		}

		addr, err := net.ResolveTCPAddr(conn.RemoteAddr().Network(), conn.RemoteAddr().String())
		peer := Peer{Connection: conn, Address: *addr}
		go callback(peer)
	}
}

// Stop stops listening for peer requests
func (c *BittorentNetwork) Stop() error {
	if c.listener != nil {
		return c.listener.Close()
	}

	return nil
}
