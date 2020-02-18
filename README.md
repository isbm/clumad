# ccd

Client controller daemon.

Currently it is working only with Uyuni Server Cluster Extension. But
essentially it is a daemon that does specific actions on events,
coming from the Cluster controller and that should be abstracted to a
configuration and plugins layer, rather then embedded directly into
the code.

As of now, the `ccd` does the following:

- Configures Salt Minion against specific Uyuni Server node in the cluster
- Starts Salt Minion
- Listens to the cluster controller endpoint for events

If the cluster controller responds with another status and/or points
to another cluster node, `ccd` will stop Salt Minion, reconfigure it
and restart it again.

All that should be a configurable scenario, rather then a code inside
the daemon.
