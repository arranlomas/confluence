package trickl

import (
	"log"
	"net"
	"net/http"
	"time"
	"os"

	"github.com/anacrolix/dht"
	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/missinggo/filecache"
	"github.com/anacrolix/missinggo/x"
	"github.com/anacrolix/tagflag"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/storage"

	"github.com/anacrolix/confluence/confluence"
)

var flags = struct {
	Addr          string        `help:"HTTP listen address"`
	DHTPublicIP   net.IP        `help:"DHT secure IP"`
	CacheCapacity tagflag.Bytes `help:"Data cache capacity"`
	TorrentGrace  time.Duration `help:"How long to wait to drop a torrent after its last request"`
	FileDir       string        `help:"File-based storage directory, overrides piece storage"`
	Seed          bool          `help:"Seed data"`
}{
	Addr:          ":8080",
	CacheCapacity: 10 << 30,
	TorrentGrace:  time.Minute,
	FileDir: "/storage/emulated/0/confluence",
}

func newAndroidTorrentClient(mWorkingDir string) (ret *torrent.Client, err error) {
	blocklist, err := iplist.MMapPacked("packed-blocklist")
	if err != nil {
		log.Print(err)
	}
	storage := func() storage.ClientImpl {
		return storage.NewFile(mWorkingDir)
		log.Printf("FILE DIR FLAG %s", flags.FileDir)
		if flags.FileDir != "" {
			return storage.NewFile(flags.FileDir)
		}
		fc, err := filecache.NewCache("filecache")
		x.Pie(err)
		fc.SetCapacity(flags.CacheCapacity.Int64())
		storageProvider := fc.AsResourceProvider()
		return storage.NewResourcePieces(storageProvider)
	}()

	log.Printf("STORAGE %s", storage)
	return torrent.NewClient(&torrent.Config{
		IPBlocklist:    blocklist,
		DefaultStorage: storage,
		DHTConfig: dht.ServerConfig{
			PublicIP: flags.DHTPublicIP,
		},
		Seed: flags.Seed,
	})
}

func AndroidMain(mWorkingDir string) {
	log.Printf("WD INPUT %s", mWorkingDir)
	wd, _ := os.Getwd()
	log.Printf("START WD %s", wd)
	os.MkdirAll(mWorkingDir, 0777)
	os.Chdir(mWorkingDir)
	wd2, _ := os.Getwd()
	log.Printf("START WD after Chdir %s", wd2)
	log.Printf("AFTER flag %s", flags.FileDir)
	flags.FileDir = mWorkingDir
	log.Printf("BEFORE flag %s", flags.FileDir)

	log.SetFlags(log.Flags() | log.Lshortfile)
	tagflag.Parse(&flags)
	cl, err := newAndroidTorrentClient(mWorkingDir)
	if err != nil {
		log.Fatalf("error creating torrent client: %s", err)
	}
	defer cl.Close()
	l, err := net.Listen("tcp", flags.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("serving http at %s", l.Addr())
	h := &confluence.Handler{cl, flags.TorrentGrace}
	err = http.Serve(l, h)
	if err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	log.SetFlags(log.Flags() | log.Lshortfile)
// 	tagflag.Parse(&flags)
// 	cl, err := newTorrentClient()
// 	if err != nil {
// 		log.Fatalf("error creating torrent client: %s", err)
// 	}
// 	defer cl.Close()
// 	l, err := net.Listen("tcp", flags.Addr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer l.Close()
// 	log.Printf("serving http at %s", l.Addr())
// 	h := &confluence.Handler{cl, flags.TorrentGrace}
// 	err = http.Serve(l, h)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func newTorrentClient() (ret *torrent.Client, err error) {
// 	blocklist, err := iplist.MMapPacked("packed-blocklist")
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	storage := func() storage.ClientImpl {
// 		if flags.FileDir != "" {
// 			return storage.NewFile(flags.FileDir)
// 		}
// 		fc, err := filecache.NewCache("filecache")
// 		x.Pie(err)
// 		fc.SetCapacity(flags.CacheCapacity.Int64())
// 		storageProvider := fc.AsResourceProvider()
// 		return storage.NewResourcePieces(storageProvider)
// 	}()
// 	return torrent.NewClient(&torrent.Config{
// 		IPBlocklist:    blocklist,
// 		DefaultStorage: storage,
// 		DHTConfig: dht.ServerConfig{
// 			PublicIP: flags.DHTPublicIP,
// 		},
// 		Seed: flags.Seed,
// 	})
// }
