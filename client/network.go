package client

import (
	"errors"
	"net"
	"strconv"
)

type BittorentNetwork struct {
	listener net.Listener
}

// Start begins listening on a port and
// returns the port it selected
func (c *BittorentNetwork) Start(callback func(net.Conn)) (int, error) {

	var port = -1
	for i := 6881; i < 6890; i++ {
		log.Debug("Attempting to listen on %d", i)
		listen, err := net.Listen("tcp", ":"+strconv.Itoa(i))
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

func (c *BittorentNetwork) listenOnPort(callback func(net.Conn)) {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			log.Error("Connection to client failed: ", err)
			continue
		}

		go callback(conn)
	}
}

// Stop stops listening for peer requests
func (c *BittorentNetwork) Stop() error {
	if c.listener != nil {
		return c.listener.Close()
	}

	return nil
}
