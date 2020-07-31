package link

type Link struct {
	ID             uint
	URL            string
	ShortLinkID    string
	ShortLinkIDInt uint
	Count          uint
}

type Storage interface {
	AddLink(link *Link) (uint, error)
}
