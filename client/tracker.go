package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

// The Tracker handles talking to a bittorrent tracker
type Tracker struct {
	Url string
}

// Announce to the client the tracker exists
func (c *Tracker) Announce(request TrackerRequest) (TrackerResponse, error) {

	urlObj, err := url.Parse(c.Url)

	if err != nil {
		return TrackerResponse{}, fmt.Errorf("Unable to parse url: %s"+err.Error(), c.Url)
	}

	values := url.Values{}
	values.Add("info_hash", string(request.InfoHash))
	values.Add("peer_id", string(request.PeerId))
	values.Add("port", strconv.Itoa(request.Port))
	values.Add("uploaded", strconv.Itoa(request.Uploaded))
	values.Add("downloaded", strconv.Itoa(request.Downloaded))
	values.Add("left", strconv.Itoa(request.Left))
	values.Add("compact", strconv.Itoa(request.Compact))
	values.Add("event", string(request.Event))

	urlObj.RawQuery = values.Encode()

	result, err := http.Get(urlObj.String())
	if err != nil || result.StatusCode < 200 || result.StatusCode >= 300 {
		return TrackerResponse{}, errors.New("Tracker announce failed")
	}

	var response TrackerResponse
	bencode.Unmarshal(result.Body, &response)

	if response.PeerString != "" {
		peerList, err := c.ParsePeers(response.PeerString)
		if err != nil && len(response.PeerList) == 0 {
			return TrackerResponse{}, errors.New("Failed to parse peer string")
		}

		response.PeerList = peerList
	}

	return response, nil
}

func (c *Tracker) ParsePeers(peers string) ([]Peer, error) {
	if len(peers)%6 != 0 {
		return nil, errors.New("Peer string length is not a multiple of 6")
	}

	peerCount := len(peers) / 6
	peerList := make([]Peer, peerCount)
	for i := 0; i < peerCount; i++ {
		ipBytes := peers[i*6 : (i*6)+6]
		ip := net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
		port := int(binary.BigEndian.Uint16([]byte(ipBytes[4:6])))
		peer := Peer{Address: net.TCPAddr{IP: ip, Port: port}}

		peerList[i] = peer
	}

	return peerList, nil
}
