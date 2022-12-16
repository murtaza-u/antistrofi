package ping

import "time"

// DefaultInterval is the default ping interval if nothing is passed in
// the config.
const DefaultInterval = time.Second * 10

// Cfg is a list of configuration options that can be set to modify the
// behaviour of pings.
type Cfg struct {
	// time interval between each ping
	Interval time.Duration
}
