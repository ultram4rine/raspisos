package data

//User is a struct to save user's settings
type User struct {
	ID            int64    `db:"id"`
	Faculty       string   `db:"fac"`
	Group         string   `db:"groupnum"`
	Notifications bool     `db:"notifications"`
	IgnoreList    []string `db:"ignorelist"`
}

//Fac is a struct describing faculty
type Fac struct {
	Name string
	Link string
	Abbr string
}
