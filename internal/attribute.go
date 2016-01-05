package libnss_stns

type UserAttr struct {
	Group_Id  int
	Directory string
	Shell     string
	Password  string
	Gecos     string
	Keys      []string
}
type GroupAttr struct {
	Users []string
}
type Attr struct {
	Id   int
	Name string
	*UserAttr
	*GroupAttr
}
