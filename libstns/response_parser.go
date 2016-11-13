package libstns

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/STNS/STNS/stns"
)

func convertV1toV3Format(body []byte) ([]byte, error) {
	var attr stns.Attributes
	err := json.Unmarshal(body, &attr)

	if err != nil {
		return nil, err
	}

	mig := ResponseFormat{
		0,
		attr,
	}

	j, err := json.Marshal(mig)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func convertV2toV3Format(body []byte) ([]byte, error) {
	var res v2ResponseFormat
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	mig := ResponseFormat{
		res.MetaData.MinID,
		res.Items,
	}

	j, err := json.Marshal(mig)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func convertV3Format(b []byte, path string, minID string) ([]byte, error) {
	var err error
	attr := stns.Attributes{}
	sp := strings.Split(strings.TrimLeft(path, "/"), "/")

	if len(sp) < 2 {
		return nil, errors.New("parse error: path specification is insufficient")
	}

	switch sp[0] {
	case "user":
		if sp[1] == "list" {
			users := v3Users{}
			err = json.Unmarshal(b, &users)

			if err != nil {
				return nil, err
			}
			for _, u := range users {
				if u.Name != "" && u.ID != 0 {
					attr[u.Name] = &stns.Attribute{
						ID: u.ID,
						User: &stns.User{
							Password:  u.Password,
							GroupID:   u.GroupID,
							Directory: u.Directory,
							Shell:     u.Shell,
							Gecos:     u.Gecos,
							Keys:      u.Keys,
						},
					}
				}
			}
		} else {
			user := v3User{}
			err = json.Unmarshal(b, &user)
			if err != nil {
				return nil, err
			}

			if user.Name != "" && user.ID != 0 {
				attr[user.Name] = &stns.Attribute{
					ID: user.ID,
					User: &stns.User{
						Password:  user.Password,
						GroupID:   user.GroupID,
						Directory: user.Directory,
						Shell:     user.Shell,
						Gecos:     user.Gecos,
						Keys:      user.Keys,
					},
				}
			}
		}
	case "group":
		if sp[1] == "list" {
			groups := v3Groups{}
			err = json.Unmarshal(b, &groups)
			if err != nil {
				return nil, err
			}

			for _, g := range groups {
				if g.ID != 0 {
					attr[g.Name] = &stns.Attribute{
						ID: g.ID,
						Group: &stns.Group{
							Users: g.Users,
						},
					}
				}
			}
		} else {
			group := v3Group{}
			err = json.Unmarshal(b, &group)
			if err != nil {
				return nil, err
			}

			if group.ID != 0 {
				attr[group.Name] = &stns.Attribute{
					ID: group.ID,
					Group: &stns.Group{
						Users: group.Users,
					},
				}
			}
		}
	case "sudo":
		user := v3Sudo{}
		err = json.Unmarshal(b, &user)
		if err != nil {
			return nil, err
		}

		if user.Name != "" && user.Password != "" {
			attr[user.Name] = &stns.Attribute{
				ID: 0,
				User: &stns.User{
					Password: user.Password,
				},
			}
		}
	}

	m, _ := strconv.Atoi(minID)
	mig := ResponseFormat{
		m,
		attr,
	}
	j, err := json.Marshal(mig)
	if err != nil {
		return nil, err
	}

	return j, nil
}

type ResponseFormat struct {
	MinID int             `json:"min_id"`
	Items stns.Attributes `json:"items"`
}

type v2MetaData struct {
	APIVersion float64 `json:"api_version"`
	Result     string  `json:"result"`
	MinID      int     `json:"min_id"`
}

type v2ResponseFormat struct {
	MetaData *v2MetaData     `json:"metadata"`
	Items    stns.Attributes `json:"items"`
}

type v3User struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Password  string   `json:"password"`
	GroupID   int      `json:"group_id"`
	Directory string   `json:"directory"`
	Shell     string   `json:"shell"`
	Gecos     string   `json:"gecos"`
	Keys      []string `json:"keys"`
}

type v3Group struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Users []string `json:"users"`
}

type v3Sudo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type v3Users []v3User

type v3Groups []v3Group

type v3Sudoers []v3Sudo
