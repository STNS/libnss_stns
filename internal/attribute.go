package libnss_stns

type UserAttr struct {
	GroupId   int      `json:"group_id"`
	Directory string   `json:"directory"`
	Shell     string   `json:"shell"`
	Gecos     string   `json:"gecos"`
	Keys      []string `json:"keys"`
}
type GroupAttr struct {
	Users []string `json:"users"`
}
type Attr struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	*UserAttr
	*GroupAttr
}
