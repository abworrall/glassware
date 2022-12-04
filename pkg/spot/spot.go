package spot

import (
	"context"
	"fmt"
	"log"
	"strings"

	spotify "github.com/zmb3/spotify/v2"

	"github.com/abworrall/glassware/pkg/config"
)

func StartPlayback(c config.Config, playlistNameFrag string) error {
	ctx := context.Background()
	client := GetClient(c, false)
	if client == nil {
		return fmt.Errorf("GetClient failed")
	}

	d := findPlayer(client, c.SpotifyPlayerDevice)
	if d == nil {
		return fmt.Errorf("StartPlayback: could not find device '%s'", c.SpotifyPlayerDevice)
	}

	pl := findPlaylist(client, playlistNameFrag)
	if pl == nil {
		return fmt.Errorf("StartPlayback: could not find playlist '%s'", playlistNameFrag)
	}

	opt := spotify.PlayOptions{
		DeviceID: &d.ID,
		PlaybackContext: &pl.URI,
	}

	if err := client.PlayOpt(ctx, &opt); err != nil {
		return err
	}

	return nil
}

func StopPlayback(c config.Config) error {
	ctx := context.Background()
	client := GetClient(c, false)
	if client == nil {
		return fmt.Errorf("GetClient failed")
	}

	d := findPlayer(client, c.SpotifyPlayerDevice)
	if d == nil {
		return fmt.Errorf("StopPlayback: could not find device '%s'", c.SpotifyPlayerDevice)
	}

	opt := spotify.PlayOptions{
		DeviceID: &d.ID,
	}

	return client.PauseOpt(ctx, &opt)
}


// findDevice finds the first device whose name *contains* namefrag
func findPlayer(c *spotify.Client, namefrag string) *spotify.PlayerDevice {
	ctx := context.Background()

	devices, err := c.PlayerDevices(ctx)
	if err != nil { log.Fatal(err) }

	if len(devices) == 0 {
		log.Printf("Problem: no PlayerDevices found via Spotify web API\n")
		log.Printf("This is probably because librespot lost its session; to repair it, play something to your raspotify device via the Spotify app on your phone.\n")
	}

	for _, d := range devices {
		if strings.Contains(d.Name, namefrag) {
			return &d
		}
	}

	return nil
}

func findPlaylist(c *spotify.Client, namefrag string) *spotify.SimplePlaylist {
	ctx := context.Background()

	results, err := c.Search(ctx, namefrag, spotify.SearchTypePlaylist)
	if err != nil {
		log.Fatal(err)
	}

	for _, pl := range results.Playlists.Playlists {
		if strings.Contains(pl.Name, namefrag) {
			return &pl
		}
	}

	return nil
}
