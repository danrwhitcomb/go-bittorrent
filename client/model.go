package client

type MetaInfoData struct {
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Private     int    `bencode:"private"`
	Length      int    `bencode:"length"`
	Md5sum      string `bencode:"md5sum"`
}

type MetaInfo struct {
	Announce     string       `bencode:"announce"`
	AnnouceList  []string     `bencode:"annouce-list"`
	Info         MetaInfoData `bencode:"info"`
	Encoding     string       `bencode:"encoding"`
	CreationDate int          `bencode:"creation date"`
	CreatedBy    string       `bencode:"created by"`
}

type TrackerResponse struct {
	FailureReason  string `bencode:"failure reason"`
	WarningMessage string `bencode:"warning message"`
	Interval       int    `bencode:"interval"`
	MinInterval    int    `bencode:"min interval"`
	TrackerId      string `bencode:"tracker id"`
	Complete       int    `bencode:"complete"`
	Incomplete     int    `bencode:"incomplete"`
}

type Peer struct {
	Id      string `bencode:"peer id"`
	Address string `bencode:"ip"`
	Port    int    `bencode:"prt"`
}
