package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

type DispatcherCtx struct {
	ID           string
	Crawlers     []Crawler
	Config       Config
	CrawlerRound uint32
	Seeders      []Seeder
	Queue        chan *Seed
}

const (
	exitSuccess = iota
	exitFailed
	Version = "0.0.1"
)

var (
	showVersion bool
	showHelp    bool
	configPath  string
)

func getOptions() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Version: dispatcher/%s\n", Version)
		fmt.Fprintf(os.Stderr, "Useage: -c config\n")
	}

	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.StringVar(&configPath, "c", "./dispatcher.conf", "Set config path")

	flag.Parse()
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "Version: dispatcher/%s\n", Version)
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "Version: dispatcher %s (build by %s %s) \n\n",
		Version, runtime.Compiler, runtime.Version())
	fmt.Fprintf(os.Stderr, "Useage: dispatcher -c config\n")
	fmt.Fprintln(os.Stderr, "Options: ")
	fmt.Fprintln(os.Stderr, "      -h  :  show help")
	fmt.Fprintln(os.Stderr, "      -v  :  show version")
	fmt.Fprintln(os.Stderr, "      -c  :  set configureation file (default: ./dispatcher.conf)")
}


func ListenForConnect(ctx *DispatcherCtx) {
	address := fmt.Sprintf("%s:%d", ctx.Config.Host, ctx.Config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("can not bind: ", err)
	}

	log.Println("start accept connect from module...")
	for {
		// Dispatcher accept connection from other module.
		// It will wait their connection and save it if connect success
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error: ", err, ", from: ", conn.RemoteAddr())
		}

		log.Println("accept success from: ", conn.RemoteAddr())

		if request, err := AcceptConnect(conn); err != nil {
			log.Println("failed to establish connect with module: ", err, " , from: ", conn.RemoteAddr())
		} else {
			log.Printf("success establish connect with module: %s, from: %v", request.Role, conn.RemoteAddr())
			switch request.Role {
			case "Crawler":
				ctx.Crawlers = append(ctx.Crawlers, Crawler{
					Conn:   conn,
					Status: true,
				})
			case "Seeder":
				ctx.Seeders = append(ctx.Seeders, Seeder{
					Conn:   conn,
					Status: true,
				})
			}
		}
	}
}

func cycle(ctx *DispatcherCtx) {
	go func() {
		for {
			if len(ctx.Seeders) == 0 {
				continue
			}

			for _, seeder := range ctx.Seeders {
				if !seeder.Status {
					continue
				}
				seed, err := acceptSeed(seeder.Conn);
				if err != nil {
					log.Printf("accept seed fail: %v, from: %v\n", err, seeder.Conn.RemoteAddr())
					time.Sleep(1 * time.Second)
				}

				ctx.Queue <- seed
			}
		}
	}()


	for {
		if len(ctx.Crawlers) == 0 {
			continue
		}

		seed := <-ctx.Queue
		for _, url := range seed.Urls {
			crawler := ctx.Crawlers[ctx.CrawlerRound % uint32(len(ctx.Crawlers))]
			ctx.CrawlerRound += 1
			if crawler.Status {
				if err := Crawl(url, seed.ID, crawler.Conn); err != nil {
					log.Println("send crawl fail, mark it down: ", crawler.Conn.RemoteAddr())
					ctx.Crawlers[ctx.CrawlerRound].Status = false
				}
			}
		}
	}
}

func main() {
	getOptions()

	if showVersion {
		printVersion()
		os.Exit(exitSuccess)
	}

	if showHelp {
		printHelp()
		os.Exit(exitSuccess)
	}

	// init logger
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)

	var ctx DispatcherCtx
	// init queue
	ctx.Queue = make(chan *Seed, 1024)

	// parse config
	err := ConfigInit(configPath, &ctx.Config)
	if err != nil {
		log.Fatalln("parse config err: ", err)
	}

	// listen for other module to connect
	go ListenForConnect(&ctx)

	// seeder cycle
	cycle(&ctx)
}
