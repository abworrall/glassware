package action

import(
	"log"

	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/event"
	"github.com/abworrall/glassware/pkg/spot"
)

// SpotifyAction will start/stop a playlist when it sees
// EvDrop/EvRestore events. The playlist name is automatically derived
// from the sensor name, e.g. "Glassware C0/A0"
type SpotifyAction struct {
	Config config.Config
}

func NewSpotifyAction(c config.Config) *SpotifyAction {
	return &SpotifyAction{Config:c}
}

func (sa *SpotifyAction)String() string {
	return "SpotifyAction"
}

func (sa *SpotifyAction)ActOnEvent(e event.Event) error {
	var err error
	
	switch e.Kind {
	case event.EvDrop:
		playlistName := "Glassware " + e.SourceName
		log.Printf("    Starting playback of %q\n", playlistName)
		err = spot.StartPlayback(sa.Config, playlistName)

	case event.EvRestore:
		log.Printf("    Stopping playback\n")
		err = spot.StopPlayback(sa.Config)
	}

	return err
}
