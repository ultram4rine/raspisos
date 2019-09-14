package data

//User is a struct to save user's settings
type User struct {
	ID            int64
	Faculty       string
	Group         string
	Notifications bool
	IgnoreList    []string
}

//Fac is a struct describing faculty
type Fac struct {
	Name string
	Link string
	Abbr string
}
