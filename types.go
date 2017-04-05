package gondole

import (
	"time"
)

type Gondole struct {
	Name        string
	RedirectURI string
}

type Client struct {
	BaseURL     string
	BearerToken string
}

/*

Entities:

Everything manipulated/returned by the API
*/

type Account struct {
	ID          int64
	Username    string
	Acct        string
	DisplayName string
	Note        string
	URL         string
	Avatar      string
	Header      string
	Locked      bool
	Followers   int
	Followings  int
	Statuses    int
}

type Application struct {
	Name    string
	Website string
}

type Attachement struct {
	ID         int64
	Type       string
	URL        string
	RemoteURL  string
	PreviewURL string
	TextURL    string
}

type Card struct {
	URL         string
	Title       string
	Description string
	Image       string
}

type Context struct {
	Ancestors   []Status
	Descendents []Status
}

type Error struct {
	Text string
}

type Instance struct {
	URI         string
	Title       string
	Description string
	Email       string
}

type Mention struct {
	ID       int64
	URL      string
	Username string
	Acct     string
}

type Notification struct {
	ID        int64
	Type      string
	CreatedAt time.Time
	Account   *Account
	Status    *Status
}

type Relationship struct {
	Following  bool
	FollowedBy bool
	Blocking   bool
	Muting     bool
	Requested  bool
}
type Report struct {
	ID          int64
	ActionTaken string
}

type Result struct {
	Accounts []Account
	Statutes []Status
	Hashtags []Tag
}

type Status struct {
	ID                 int64
	URI                string
	URL                string
	Account            *Account
	InReplyToId        int64
	InReplyToAccountID int64
	Reblog             *Status
	Content            string
	CreatedAT          time.Time
	Reblogs            int
	Favourites         int
	Reblogged          bool
	Favourited         bool
	Sensitive          bool
	SpoilerText        string
	Visibility         string
	MediaAttachments   []Attachement
	Mentions           []Mention
	Tags               []Tag
	App                Application
}

type Tag struct {
	Name string
	URL  string
}
