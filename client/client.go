package client

import (
	"crypto/sha1"
	"errors"
	"go-bittorrent/lib"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var log = lib.Log

type BittorentClient struct {
	Id          []byte
	Port        int
	DownloadDir string
	WillSeed    bool
	Repository  TorrentRepository
	Network     BittorentNetwork

	IsChoked    bool
	Connections []Peer
}

func NewClient(downloadDir string, willSeed bool) BittorentClient {

	network := BittorentNetwork{}
	repo := NewRepository(downloadDir)

	return BittorentClient{
		Id:         generatePeerId(),
		WillSeed:   willSeed,
		Repository: repo,
		Network:    network}
}

// Download begins downloading the given torrent
func (c *BittorentClient) Download(torrentFile string) error {

	torrent, err := NewTorrent(torrentFile)
	if err != nil {
		return errors.New("Failed to read torrent file: " + err.Error())
	}

	// Open server
	port, err := c.Network.BindToPort(c.HandlePeerConnection)
	if err != nil {
		return errors.New("Failed to bind to port:\n" + err.Error())
	}

	c.Port = port
	c.Repository.AddFiles(torrent.Data.Info.Files)

	// Announce to tracker
	trackerData, ok := c.GetTrackerData(torrent)
	if !ok {
		return errors.New("Unable to find usable tracker")
	}

	log.Infof("Tracker data: %s", trackerData.TrackerId)

	for _, peer := range trackerData.PeerList {
		log.Infof("Peer: %s", peer.Address.String())
	}

	return nil
}

// True if all files have been downloaded
func (c *BittorentClient) Status() bool {
	return true
}

// Stops all current downloads
func (c *BittorentClient) Stop() {
	// Notify all clients of shutdown
	// Shutdown server
}

func (c *BittorentClient) GetTrackerData(torrent Torrent) (TrackerResponse, bool) {

	trackerRequest := c.getTrackerRequest(torrent)

	// First check the announce list
	if len(torrent.Data.AnnounceList) > 0 {
		for _, tier := range torrent.Data.AnnounceList {
			// Randomly select trackers from each tier
			for _, index := range rand.Perm(len(tier)) {
				tracker := Tracker{tier[index]}
				log.Debugf("(%s) Announcing...", tracker.Url)
				resp, err := tracker.Announce(trackerRequest)
				if err == nil && len(resp.FailureReason) == 0 {
					return resp, true
				}

				log.Debugf("(%s) Announce failed: .", tracker.Url)
			}
		}
	}

	// If no announce list tracker was successful,
	// try using the top level announce
	tracker := &Tracker{Url: torrent.Data.Announce}
	log.Debugf("(%s) Announcing...", tracker.Url)
	resp, err := tracker.Announce(trackerRequest)
	if err == nil {
		return resp, true
	}

	log.Debugf("(%s) Announce failed.", tracker.Url)
	return TrackerResponse{}, false
}

func (c *BittorentClient) getTrackerRequest(torrent Torrent) TrackerRequest {
	return TrackerRequest{
		InfoHash:   torrent.Hash,
		PeerId:     c.Id,
		Port:       c.Port,
		Uploaded:   c.Repository.GetUploaded(),
		Downloaded: c.Repository.GetDownloaded(),
		Left:       c.Repository.GetLeft(),
		Compact:    1, // Always go with compact
		Event:      Started}
}

func generatePeerId() []byte {
	hash := sha1.New()
	// Current time
	hash.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	hash.Write([]byte(strconv.Itoa(os.Getpid())))

	// Process ID
	return hash.Sum(nil)
}

func (c *BittorentClient) HandlePeerConnection(peer Peer) {
	defer peer.Connection.Close()
}
