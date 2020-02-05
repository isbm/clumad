package uccd

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
}

func NewCDStats() *CDStats {
	stats := new(CDStats)
	return stats
}
