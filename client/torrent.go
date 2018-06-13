package client

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type Torrent struct {
	Path string
	Data MetaInfo
	Hash []byte
}

func NewTorrent(path string) (Torrent, error) {
	log.Debugf("Opening %s", path)
	file, err := os.Open(path)
	if err != nil {
		return Torrent{}, errors.New("Failed to open torrent file: " + err.Error())
	}
	defer file.Close()

	log.Debug("Decoding torrent file")
	info := MetaInfo{}
	err = bencode.Unmarshal(file, &info)
	if err != nil {
		return Torrent{}, errors.New("Failed to decode torrent file: " + err.Error())
	}

	log.Debug("Computing torrent info hash")
	infoHash, err := computeInfoHash(path)
	if err != nil {
		return Torrent{}, errors.New("Failed to compute info hash: " + err.Error())
	}

	log.Debugf("Announce URL: %s", info.Announce)
	log.Debugf("Hash: %x", infoHash)

	return Torrent{Path: path, Data: info, Hash: infoHash}, nil
}

func computeInfoHash(torrentPath string) ([]byte, error) {

	file, err := os.Open(torrentPath)
	if err != nil {
		return nil, errors.New("Failed to open torrent: " + err.Error())
	}

	data, err := bencode.Decode(file)
	if err != nil {
		return nil, errors.New("Failed to decode torrent file: " + err.Error())
	}

	torrentDict, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("Torrent file is not a dictionary")
	}

	infoBuffer := bytes.Buffer{}
	err = bencode.Marshal(&infoBuffer, torrentDict["info"])
	if err != nil {
		return nil, errors.New("Failed to encode info dict: " + err.Error())
	}

	hash := sha1.New()
	hash.Write(infoBuffer.Bytes())
	return hash.Sum(nil), nil
}
