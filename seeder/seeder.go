package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type SeederCtx struct {
	Role           string
	DispatcherConn net.Conn
	Established    bool
	Config         Config
}


func connectDispatcher(conf *Config, ctx *SeederCtx) error {
	server := fmt.Sprintf("%s:%d", conf.DispatcherConfig.Host, conf.DispatcherConfig.Port)

	conn, err := net.Dial("tcp4", server)
	if err != nil {
		return err
	}

	ctx.DispatcherConn = conn
	return nil
}

func doSeed(ctx *SeederCtx) {
	for {
		if err := Connect(ctx.DispatcherConn, ctx.Role); err != nil {
			log.Println("failed to Connect to dispatcher: ", err)
			time.Sleep(ctx.Config.DispatcherConfig.Interval * time.Second)
		} else {
			log.Println("Connect success to dispatcher, waiting for task")
			ctx.Established = true
			break
		}
	}

	for {
		if ctx.Established {
			url := "http://segmentfault.com/a/1190000024564062"
			if err := seed(url, ctx.DispatcherConn); err != nil {
				log.Println("seed error: ", err)
				ctx.Established = false
			} else {
				log.Println("seed url to dispatcher, url: ", url)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// init log
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)

	var ctx SeederCtx

	// init config
	err := ConfigInit("./seeder.conf", &ctx.Config)
	if err != nil {
		log.Fatalln("parse config err: ", err)
	}
	log.Println("init config.")

	ctx.Role = "Seeder"

	// try connect to dispatcher
	if err := connectDispatcher(&ctx.Config, &ctx); err != nil {
		log.Fatalln("can not connect dispatcher: ", err)
	}
	log.Println("build tcp connection to dispatcher.")

	doSeed(&ctx)
}