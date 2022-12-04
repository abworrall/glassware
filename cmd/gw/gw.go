package main

import(
	"flag"
	"log"
	"os"
	"strings"

	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/controller"
	"github.com/abworrall/glassware/pkg/eventloop"
	"github.com/abworrall/glassware/pkg/spot"
)

var(
	fVerbosity int
	fCacheDir string

	fActiveSensors string // unwired sensors may generate ghost events; we need to be told which sensors aren't ghosts
	
	// Get ID & Secret from Spotify dev portal, https://developer.spotify.com/dashboard/
	fSpotifyId string
	fSpotifySecret string
	fSpotifyPlayerDevice string
	fSpotifyInit bool

	c config.Config
)

func init() {
	flag.IntVar(&fVerbosity, "v", 0, "logging verbosity")
	flag.StringVar(&fCacheDir, "dir", os.Getenv("HOME")+"/.gw", "where to store auth tokens etc.")

	flag.StringVar(&fActiveSensors, "sensors", "C0/A0", "comma-sep list of the sensors that should be monitored")

	// These two values are needed just as a one-off for --spotify-init; they are cached
	flag.StringVar(&fSpotifyId, "spotify-id", "", "The spotify ID for Spotify Connect (see spotify dev portal)")
	flag.StringVar(&fSpotifySecret, "spotify-secret", "", "The spotify Secret for Spotify Connect (see spotify dev portal)")

	flag.StringVar(&fSpotifyPlayerDevice, "spotify-device", "raspotify", "full/partial name of the device to play music on")
	flag.BoolVar(&fSpotifyInit, "spotify-init", false, "perform OAuth web auth flow (needs a browser) and cache token")

	flag.Parse()

	c = config.NewConfig()

	c.Verbosity = fVerbosity
	c.CacheDir = fCacheDir
	if err := os.MkdirAll(fCacheDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for _, s := range strings.Split(fActiveSensors, ",") {
		c.ActiveSensors[s] = 1
	}

	c.SpotifyId = fSpotifyId
	c.SpotifySecret = fSpotifySecret
	c.SpotifyPlayerDevice = fSpotifyPlayerDevice

	if c.SpotifyId == "" {
		// If they're not configured, try and load them from the cache
		if id, secret, err := spot.LoadCreds(c.CacheDir); err == nil {
			c.SpotifyId = id
			c.SpotifySecret = secret
		}
	}
}

func main() {
	if fSpotifyInit {
		log.Printf("Initializing spotify, with OAuth flow; you will need a browser that can see localhost:8081 on this host\n")
		if client := spot.GetClient(c, true); client != nil {
			// save the creds
			if err := spot.StoreCreds(c.CacheDir, fSpotifyId, fSpotifySecret); err != nil {
				log.Fatal(err)
			}
		}
		return
	}

	eventloop.New(c).Run(controller.InitControllers())
}
