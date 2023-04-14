package core

import "net"

type Claim struct {
	IPNet *net.IPNet
	Name  string
}
