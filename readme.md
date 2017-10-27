# ipdash
> InterPlanetary Dashboard

A community of self-reporting ipfs nodes around the world!

## Architecture

ipdash has three parts to it.

- ipdash-agent - runs alongside with a ipfs daemon to report stats
- ipdash-server - an instance to collects stats and broadcast them. Also keeps history
- ipdash-web - webui for viewing the stats from ipdash-server

## Channels

- `ipfs-dashboard-stats` - raw stats meant for server
- `ipfs-dashboard-feed` - processed stats meant for webui
