package client

//
type MetaInfoData struct {
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Private     int    `bencode:"private"`
	Length      int    `bencode:"length"`
	Md5sum      string `bencode:"md5sum"`
	Files       []File `bencode:"files"`
}

// Description of an available file in the torrent
type File struct {
	Length int    `bencode:"length"`
	Md5sum string `bencode:"md5sum"`
	Path   string `bencode:"path"`
}

// .torrent file description. Mostly meta data about the torrent
type MetaInfo struct {
	Announce     string       `bencode:"announce"`
	AnnounceList [][]string   `bencode:"announce-list"`
	Info         MetaInfoData `bencode:"info"`
	Encoding     string       `bencode:"encoding"`
	CreationDate int          `bencode:"creation date"`
	CreatedBy    string       `bencode:"created by"`
}

// A ClientEvent indicates what the client is doing
type ClientEvent string

const (
	Started   ClientEvent = "started"
	Stopped               = "stopped"
	Completed             = "completed"
)

type TrackerRequest struct {
	InfoHash   []byte
	PeerId     []byte
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Compact    int
	Event      ClientEvent
}

type TrackerResponse struct {
	FailureReason  string `bencode:"failure reason"`
	WarningMessage string `bencode:"warning message"`
	Interval       int    `bencode:"interval"`
	MinInterval    int    `bencode:"min interval"`
	TrackerId      string `bencode:"tracker id"`
	Complete       int    `bencode:"complete"`
	Incomplete     int    `bencode:"incomplete"`
	PeerString     string `bencode:"peers"`
	PeerList       []Peer `bencode:"peers"`
}
