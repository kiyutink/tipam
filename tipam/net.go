package tipam

import "net"

// isSubnet checks whether 'sub' is a subnet of 'super'
func isSubnet(sub *net.IPNet, super *net.IPNet) bool {
	onesSub, _ := sub.Mask.Size()
	onesSuper, _ := super.Mask.Size()

	// The subnet's mask should be longer, constituting a smaller CIDR range
	if onesSub <= onesSuper {
		return false
	}

	// The subnet's network address must lie within the supernet's CIDR range
	if !super.Contains(sub.IP) {
		return false
	}

	return true
}
