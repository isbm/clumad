## Required Setup

The `mgr-uccd` should have a capability to connect to a raw sockets in
order to perform ICMP pings:

	setcap cap_net_raw=+ep /path/to/binary/mgr-uccd

To reverse the above:

	setcap -r /path/to/binary/mgr-uccd

List the capabilities:

	getcap /path/to/binary/mgr-uccd

