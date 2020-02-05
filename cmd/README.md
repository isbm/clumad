## Required Setup

The `clumad` should have a capability to connect to a raw sockets in
order to perform ICMP pings:

	setcap cap_net_raw=+ep /path/to/binary/clumad

To reverse the above:

	setcap -r /path/to/binary/clumad

List the capabilities:

	getcap /path/to/binary/clumad

