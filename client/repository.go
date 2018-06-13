package client

// Handles managing torrent file state
type TorrentRepository struct {
	Directory string
	Files     []File
}

func NewRepository(path string) TorrentRepository {
	return TorrentRepository{Directory: path, Files: make([]File, 1)}
}

func (c *TorrentRepository) AddFiles(files []File) {
	for _, file := range files {
		c.Files = append(c.Files, file)
	}
}

// TODO: Actually calculate values

func (c *TorrentRepository) GetUploaded() int {
	return 0
}

func (c *TorrentRepository) GetDownloaded() int {
	return 0
}

func (c *TorrentRepository) GetLeft() int {
	sum := 0
	for _, file := range c.Files {
		sum += file.Length
	}

	return sum
}
