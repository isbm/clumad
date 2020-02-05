package uccd

import (
	"github.com/sparrc/go-ping"
)

// Pinger object
type Pinger struct {
	pinger *ping.Pinger
}

// Constructor
func NewPinger() *Pinger {
	p := new(Pinger)
	return p
}

// Ping method pings a network address and returns statistics
func (p *Pinger) Ping(addr string) *ping.Statistics {
	var err error
	var stats *ping.Statistics
	if p.pinger == nil {
		p.pinger, err = ping.NewPinger(addr)
		if err != nil {
			panic("Unabpe to initialise pinger to the address " + addr + ": " + err.Error())
		}
		p.pinger.Count = 3
		p.pinger.SetPrivileged(true) // Still needs a raw socket capability to be set
		p.pinger.Run()
		stats = p.pinger.Statistics()
		p.pinger = nil
	}
	return stats
}
