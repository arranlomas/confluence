This is a branch of the confluence wrapper for anacrolix torrent, this is so that it works nicely on mobile, the main changes are the changes to the main method and the ip address, 

# confluence

Confluence is a torrent client as a HTTP service. This allows for easy use from other processes, languages, and machines, due to the ubiquity of HTTP. It makes use of [anacrolix/torrent](https://github.com/anacrolix/torrent)'s [download-on-demand](https://godoc.org/github.com/anacrolix/torrent#Torrent.NewReader) torrenting, and [custom data backend](https://godoc.org/github.com/anacrolix/torrent/storage#ClientImpl) features to store data in a cache. You can then utilize the BitTorrent network with sensible defaults as though it were just regular HTTP.

before setup ensure that go and gomobile are setup correctly
# Installation
```
$ go get github.com/anacrolix/confluence
$ go get github.com/arranlomas/torrent
$ go get github.com/arranlomas/confluence
copy all files from arranlomas/confluence to ancrolix/confluence
copy all files from arranlomas/torrent to ancrolix/torrent

this is a super quick and dirty setup for now - will be improved at a later date

```
# Usage
```
$ gomobile bind github.com/arranlomas/confluence -target=android

then use the aar as a module in you android project
to start the daemon call trickl.Trickl.androidMain(workingDir.absolutePath) in your android project

Confluence will announce itself to DHT, and wait for HTTP activity. Torrents are added to the client as needed. Without an active request on a torrent, it is kicked from the client after one minute. Its data however may remain in the cache for future uses of that torrent.

# Routes

There are several routes to interact with torrents:

 * `GET /data?ih=<infohash in hex>&path=<display path of file declared in torrent info>`. Note that this handler supports HTTP range requests for bytes. Response will block until the data is available.
 * `GET /status`. This fetches the textual status info page per anacrolix/torrent.Client.WriteStatus. Very useful for debugging.
 * `GET /info?ih=<infohash in hex>`. This returns the info bytes for the matching torrent. It's useful if the caller needs to know about the torrent, such as what files it contains. It will block until the info is available. The response is the full bencoded info dictionary per [BEP 3](http://www.bittorrent.org/beps/bep_0003.html).
 * `/events?ih=<infohash in hex>`. This is a websocket that emits frames with [confluence.Event] encoded as JSON for the torrent. The PieceChanged field for instance is set if the given piece changed [state](https://godoc.org/github.com/anacrolix/torrent#PieceState) within the torrent.
 * `GET /fileState?ih=<infohash in hex>&path=<display path of file declared in torrent info>`. Returns [file state](https://godoc.org/github.com/anacrolix/torrent#File.State) encoded as JSON.
 * `POST /metainfo?ih=<infohash in hex>`. Updates a torrents metainfo. New trackers are added. Info bytes are set if not already known. The request body must be the bencoded metainfo, the same as what would appear in a `.torrent` file.
