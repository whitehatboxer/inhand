package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type CrawlerCtx struct {
	Role           string
	DispatcherConn net.Conn
	Established    bool
	TaskQueue      chan CrawlerTask
	Config         Config
	ResponseQueue  chan Response
}

type CrawlerTask struct {
	ID  string
	URL string
}

func connectDispatcher(conf *Config, ctx *CrawlerCtx) error {
	server := fmt.Sprintf("%s:%d", conf.DispatcherConfig.Host, conf.DispatcherConfig.Port)

	conn, err := net.Dial("tcp4", server)
	if err != nil {
		return err
	}

	ctx.DispatcherConn = conn
	return nil
}


func collectTask(ctx *CrawlerCtx) {
	reader := bufio.NewReader(ctx.DispatcherConn)
	req, err := parseRequest(reader)
	if err != nil {
		return
	}

	items := strings.Split(req.Body, "\n")
	for _, item := range items {
		ctx.TaskQueue <- CrawlerTask{
			ID:  req.ID,
			URL: item,
		}
	}
}

func handleConnection(ctx *CrawlerCtx) {
	for {
		if err := connect(ctx.DispatcherConn, ctx.Role); err != nil {
			log.Println("failed to Connect to dispatcher: ", err)
			time.Sleep(ctx.Config.DispatcherConfig.Interval * time.Second)
		} else {
			log.Println("Connect success to dispatcher, waiting for task")
			ctx.Established = true
			break
		}
	}

	for {
		collectTask(ctx)
	}
}

func main() {
	// init log
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)

	var ctx CrawlerCtx

	// init config
	err := ConfigInit("./crawler.conf", &ctx.Config)
	if err != nil {
		log.Fatalln("parse config err: ", err)
	}
	log.Println("init config.")

	ctx.Role = "Crawler"

	// try connect to dispatcher
	if err := connectDispatcher(&ctx.Config, &ctx); err != nil {
		log.Fatalln("can not connect dispatcher: ", err)
	}
	log.Println("build tcp connection to dispatcher.")

	// handle dispatcher conn
	ctx.TaskQueue = make(chan CrawlerTask , 1024)
	go handleConnection(&ctx)

	for {
		time.Sleep(ctx.Config.DispatcherConfig.Interval * time.Second)

		// waiting for established
		if ctx.Established {
			break
		}
		log.Println("waiting connect to dispatcher...")
	}

	// start crawler and wait for task from dispatcher
	ctx.ResponseQueue = make(chan Response, 1024)
	var crawlerWaitGroup sync.WaitGroup
	for i := 0; i < ctx.Config.CrawlerConfig.Goroutines; i++ {
		crawlerWaitGroup.Add(1)

		go func() {
			for {
				task := <- ctx.TaskQueue
				log.Println("get task, url: ", task.URL)
				ctx.ResponseQueue <- Crawl(task.URL)
			}
		}()
	}

	// upload crawl res to dispatcher
	go func() {
		for {
			res := <- ctx.ResponseQueue
			if res.Status >= 200 && res.Status < 300 {
				log.Println("upload response to dispatcher")
				if err := upload(res, ctx.DispatcherConn); err != nil {
					log.Println("upload task to dispatcher failed: ", err)
				}
			}
		}
	}()


	crawlerWaitGroup.Wait()
}