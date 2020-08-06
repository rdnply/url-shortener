package link

type Link struct {
	ID         uint `db:"link_id"`
	URL        string
	ShortID    string `db:"short_link_id"`
	ShortIDInt uint   `db:"short_link_id_int"`
	Clicks     uint   `db:"count_clicks"`
}

type Storage interface {
	AddLink(link *Link) (uint, error)
	GetLinkByShortID(shortID string) (*Link, error)
	IncrementLinkCounter(link *Link) (uint, error)
}
