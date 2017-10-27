# IPFS-Dashboard
> A un-official dashboard of IPFS and it's network

## !! WARNING: Does not account for the entire network, only nodes who are registered in IPFS-Dashboard !!

## Nodes List
> List of connected nodes

Properties per node:

- Status: online/offline
- Name: self-declared name of node
- Useragent: implementation, version, supported protocols and so on
- Ping: how long time it took to ping it
- Peers: how many peers are this node connected to
- Last Contact: how long time ago it was we last contacted this node
- Uptime: how long time this node have been online for

## Map of connected peers and the peers connected

## Found hashes - table

## Total number of peer, line-graph

## Total Bandwidth

# How stats are collected

You'll need to have a IPFS daemon running, together with the ipfs-dashboard-companion
software. The ipfs-dashboard-companion will connect to your IPFS daemon and
extract the details needed for submission, and sends the data via pubsub.

If you node is not reachable, the ipfs-dashboard-companion won't work either as
you node will need to be connected to the IPFS daemon serving ipfs-dashboard.
