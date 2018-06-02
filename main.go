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

	client := client.NewClient()
	err := client.Download(args.Torrent)
	if err != nil {
		log.Error("Failed to download torrent: " + err.Error())
	}
}

type Args struct {
	Torrent  string
	WillSeed bool
	Verbose  bool
}

func ParseArgs() Args {
	torrentOption := flag.String("f", "", "Torrent file to download (required)")
	seedOption := flag.Bool("s", false, "Will seed if provided")
	verbosityOption := flag.Bool("v", false, "Turns on debug logging")

	flag.Parse()

	if *torrentOption == "" {
		flag.Usage()
		os.Exit(1)
	}

	return Args{*torrentOption, *seedOption, *verbosityOption}
}
