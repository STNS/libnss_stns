package libstns

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/settings"
)

func convertV1toV3Format(body []byte) (*ResponseFormat, error) {
	var attr stns.Attributes
	err := json.Unmarshal(body, &attr)

	if err != nil {
		return nil, err
	}

	return &ResponseFormat{
		0,
		attr,
	}, nil
}

func convertV2toV3Format(body []byte) (*ResponseFormat, error) {
	var res v2ResponseFormat
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &ResponseFormat{
		res.MetaData.MinID,
		res.Items,
	}, nil

}

func convertV3Format(b []byte, path string, minID string, config *Config) (*ResponseFormat, error) {
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
				if u.Name != "" && u.ID+config.UIDShift > settings.MIN_LIMIT_ID {
					tmpUser := &stns.User{
						Password:      u.Password,
						Directory:     u.Directory,
						Shell:         u.Shell,
						Gecos:         u.Gecos,
						Keys:          u.Keys,
						SetupCommands: u.SetupCommands,
					}

					if u.GroupID+config.GIDShift > settings.MIN_LIMIT_ID {
						tmpUser.GroupID = u.GroupID + config.GIDShift
					}

					attr[u.Name] = &stns.Attribute{
						ID:   u.ID + config.UIDShift,
						User: tmpUser,
					}
				}
			}
		} else {
			u := v3User{}
			err = json.Unmarshal(b, &u)
			if err != nil {
				return nil, err
			}

			if u.Name != "" && u.ID+config.UIDShift > settings.MIN_LIMIT_ID {
				tmpUser := &stns.User{
					Password:      u.Password,
					Directory:     u.Directory,
					Shell:         u.Shell,
					Gecos:         u.Gecos,
					Keys:          u.Keys,
					SetupCommands: u.SetupCommands,
				}

				if u.GroupID+config.GIDShift > settings.MIN_LIMIT_ID {
					tmpUser.GroupID = u.GroupID + config.GIDShift
				}

				attr[u.Name] = &stns.Attribute{
					ID:   u.ID + config.UIDShift,
					User: tmpUser,
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
				if g.ID+config.GIDShift > settings.MIN_LIMIT_ID {
					attr[g.Name] = &stns.Attribute{
						ID: g.ID + config.GIDShift,
						Group: &stns.Group{
							Users: g.Users,
						},
					}
				}
			}
		} else {
			g := v3Group{}
			err = json.Unmarshal(b, &g)
			if err != nil {
				return nil, err
			}

			if g.ID+config.GIDShift > settings.MIN_LIMIT_ID {
				attr[g.Name] = &stns.Attribute{
					ID: g.ID + config.GIDShift,
					Group: &stns.Group{
						Users: g.Users,
					},
				}
			}
		}
	case "sudo":
		u := v3Sudo{}
		err = json.Unmarshal(b, &u)
		if err != nil {
			return nil, err
		}

		if u.Name != "" && u.Password != "" {
			attr[u.Name] = &stns.Attribute{
				ID: 0,
				User: &stns.User{
					Password: u.Password,
				},
			}
		}
	}

	m, _ := strconv.Atoi(minID)
	return &ResponseFormat{
		m,
		attr,
	}, nil
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
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Password      string   `json:"password"`
	GroupID       int      `json:"group_id"`
	Directory     string   `json:"directory"`
	Shell         string   `json:"shell"`
	Gecos         string   `json:"gecos"`
	Keys          []string `json:"keys"`
	SetupCommands []string `json:"setup_commands"`
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
