package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
)

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
	flag.StringVar(&configPath, "c", "./config.json", "Set config path")

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
	fmt.Fprintln(os.Stderr, "      -c  :  set configureation file (default: ./config.json)")
}

func listenForConnect() {
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

// run dispatcher
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

	// parse config
	ConfigInit(configPath)

	// listen for other module to connect
	go listenForConnect()

}
