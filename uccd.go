package uccd

import (
	"fmt"
	"log"
	"time"
)

type Uccd struct {
	sysop   *SysOp
	pinger  *Pinger
	cdstats *CDStats

	masterFqdn    string
	minionCfgPath string
	pingInterval  int64
}

func NewUccd() *Uccd {
	d := new(Uccd)
	d.sysop = NewSysOp()
	d.cdstats = NewCDStats()
	d.pinger = NewPinger()
	d.pingInterval = 10

	return d
}

func (d *Uccd) setup() {
	// Get current Salt Minion master's hostname
	// and if it is set in the main Minion config,
	// rewrite it as an overlay configuration.
	var _mfqdn interface{}
	_mfqdn, d.minionCfgPath = d.sysop.GetSaltOps().GetConfDOption("minion", "master")
	if _mfqdn != nil {
		d.masterFqdn = _mfqdn.(string)
	}

	// Overlay wasn't found, try main config.
	if d.masterFqdn == "" {
		d.masterFqdn = d.sysop.GetSaltOps().GetConfOpion("minion", "master").(string)
	}

	// No hostname found at all, so just set localhost for now and transfer it to the overlay.
	if d.masterFqdn == "" || d.masterFqdn == "salt" {
		d.masterFqdn = "localhost"
	}
	d.sysop.GetSaltOps().SetConfDOption("minion", "master", d.masterFqdn)

	// Set pubkey PEM fingerprint
	d.cdstats.SetPubKeyFP(d.cdstats.GetPubKeyFp("pki/minion/minion.pub"))
}

// SetSaltConfigPath sets configuration path to the Salt Minion
func (d *Uccd) SetSaltConfigPath(confpath string) *Uccd {
	d.sysop.GetSaltOps().SetSaltConfigDir(confpath)
	d.cdstats.SetSaltConfigDir(confpath)
	return d
}

// SetSaltConfigPath sets configuration path to the Salt Minion
func (d *Uccd) SetSaltExec(execpath string) *Uccd {
	d.sysop.GetSaltOps().SetSaltExecPath(execpath)
	return d
}

// SetPingTimeDuration sets ping pause duration to ping Cluster Node
func (d *Uccd) SetPingTimeDuration(duration int64) *Uccd {
	d.pingInterval = duration
	return d
}

// SetClusterURL sets cluster main entry URL
func (d *Uccd) SetClusterURL(url string) *Uccd {
	d.cdstats.SetClusterDirectorURL(url)
	return d
}

// Pinger loop is used to gather physical networn statistics to the Cluster Node.
func (d *Uccd) pingerLoop() {
	for {
		netstat := d.pinger.Ping(d.masterFqdn)
		if netstat != nil {
			log.Printf("Host '%s' max RTT (ms): %d; min RTT (ms): %d; packets sent: %d",
				d.masterFqdn, netstat.MaxRtt.Milliseconds(), netstat.MinRtt.Milliseconds(), netstat.PacketsSent)
		}
		time.Sleep(time.Duration(d.pingInterval) * time.Second)
	}
}

// Heartbeat loop is used to get current Cluster Node FQDN.
// If new FQDN comes, it means that currently registered
// system should switch elsewhere. For this, the current
// Salt Minion is reconfigured and then restarted against
// new Salt Master.
func (d *Uccd) heartbeatLoop() {
	for {
		fmt.Println("Poke")
		time.Sleep(time.Duration(d.pingInterval) * time.Second)
	}
}

// Start the Salt Minion service
func (d *Uccd) Start() {
	d.setup()
	go d.pingerLoop()
	d.heartbeatLoop()
}
