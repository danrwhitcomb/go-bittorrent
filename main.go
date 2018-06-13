package main

import (
	"bittorrent/client"
	"bittorrent/lib"
	"flag"
	"os"
)

var log = lib.Log

func main() {
	args := ParseArgs()
	lib.ConfigureLogging(args.Verbose)
	log.Infof("Starting download of torrent %s. Seeding: %t\n", args.Torrent, args.WillSeed)

	client := client.NewClient(args.DownloadDir, args.WillSeed)

	err := client.Download(args.Torrent)
	if err != nil {
		log.Fatal("Failed to bind client", err)
	}
}

// Args hold command line arguments
type Args struct {
	Torrent     string
	DownloadDir string
	WillSeed    bool
	Verbose     bool
}

// ParseArgs manages command line arguments
func ParseArgs() Args {
	torrentOption := flag.String("f", "", "Torrent file to download (required)")
	downloadPath := flag.String("d", "", "Directory will files will be downloaded (required). Directory will be created if it does not exist")
	seedOption := flag.Bool("s", false, "Will seed if provided")
	verbosityOption := flag.Bool("v", false, "Turns on debug logging")

	flag.Parse()

	if *torrentOption == "" || *downloadPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	return Args{*torrentOption, *downloadPath, *seedOption, *verbosityOption}
}
