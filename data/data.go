package data

//User is a struct to save user's settings
type User struct {
	ID            int64    `json:"id"`
	Faculty       string   `json:"faculty"`
	Group         string   `json:"group"`
	Notifications bool     `json:"notifications"`
	IgnoreList    []string `json:"ignore"`
}

//Fac is a struct describing faculty
type Fac struct {
	Name string
	Link string
	Abbr string
}
