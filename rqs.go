package uccd

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"path"
	"strings"
)

/*
Cluster Director Stats are set with the public key fingerprint for the current client machine.
Status and Cluster Node FQDN are received from the Cluster Director.

If status is "orphan", uccd should shutdown Salt Minion.
*/

const (
	CS_NEW        = iota // new to the cluster
	CS_REGISTERED        // registered in the cluster
	CS_ORPHANED          // orphaned from the current cluster node, needs remapping
	CS_DELETED           // deleted from the cluster, needs complete shutdown of all cluster services
)

// Cluster Director Stats
type CDStats struct {
	ClusterNodeFQDN string
	PubKeyFP        string
	Status          int
	saltConfPath    string
}

func NewCDStats() *CDStats {
	stats := new(CDStats)
	return stats
}

// Set Salt configuration directory
func (cds *CDStats) SetSaltConfigDir(confpath string) *CDStats {
	cds.saltConfPath = confpath
	return cds
}

// GetPubKeyFp is just reimplementation of "salt.utils.crypt.pem_finger".
// Salt gets fingerprint by not stripping any newlines "\n" symbol from the
// key cipher text.
//
// Default location of keys is $SALT_CONFIG/pki. If path starts with "/",
// it is used entirely (absolute), otherwise appended to the default root
// of known Salt Config root path.
func (cds *CDStats) GetPubKeyFp(keypath string) string {
	if !strings.HasPrefix(keypath, "/") {
		keypath = path.Join(cds.saltConfPath, keypath)
	}
	var fp string
	fh, err := os.Open(keypath)
	if err != nil {
		log.Printf("Unable to open PEM key file %s: %s\n", keypath, err.Error())
	} else {
		digest := sha256.New()
		defer fh.Close()
		scr := bufio.NewScanner(fh)
		for scr.Scan() {
			cipherline := scr.Text() + "\n"
			if strings.Contains(cipherline, "PUBLIC KEY-----") {
				continue
			}
			digest.Write([]byte(cipherline))
		}
		for idx, ch := range hex.EncodeToString(digest.Sum(nil)) {
			if idx%2 != 0 {
				fp += string(ch) + ":"
			} else {
				fp += string(ch)
			}
		}
		fp = strings.Trim(fp, ":")
	}

	return fp
}
