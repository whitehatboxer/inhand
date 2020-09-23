package main

import (
	"net"
)

// define struct about a conn of Crawler and other info
type Crawler struct {
	Conn   net.Conn
	Status bool
}

type Seeder struct {
	Conn   net.Conn
	Status bool
}

// A seed from seeder
type Seed struct {
	MaxDepth int32
	Depth    int32
	//ParentID string
	ID       string
	Urls     []string
}