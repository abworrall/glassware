package config

type Config struct {
	Verbosity int
	CacheDir string  // where we store auth tokens etc

	SpotifyId string
	SpotifySecret string
	SpotifyPlayerDevice string

	ActiveSensors map[string]int // the sensors we've been told to pay attention to
}

func NewConfig() Config {
	return Config{
		ActiveSensors: map[string]int{},
	}
}
