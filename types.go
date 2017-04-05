package gondole

import "time"

type Gondole struct {
	Name string
    RedirectURI string
}

type Account struct {
	ID         int
	Acct       string
	Avatar     string
	Followers  int
	Followings int
	Header     string
	Note       string
	Statuses   int
	URL        string
	Username   string
}

type Client struct {
	BaseURL     string
	BearerToken string
}

type Status struct {
	ID         int
	Account    *Account
	Content    string
	CreatedAT  time.Time
	Favourited bool
	Favourites int
	InReplyTo  int
	Reblog     *Status
	Reblogged  bool
	Reblogs    int
	URI        string
	URL        string
}
