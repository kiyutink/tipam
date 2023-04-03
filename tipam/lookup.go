package tipam

func (t *Tipam) lookupTags(cidr string) []string {
	return t.TagsByCIDR[cidr]
}
