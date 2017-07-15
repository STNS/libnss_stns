package libstns

import (
	"encoding/json"
	"errors"
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
		res.Items,
	}, nil

}

func uidShift(attr stns.Attributes, u *v3User, config *Config) {
	if u.Name != "" && u.ID+config.UIDShift > settings.MIN_LIMIT_ID {
		prevID := 0
		nextID := 0

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

		if u.PrevID+config.UIDShift > settings.MIN_LIMIT_ID {
			prevID = u.PrevID + config.UIDShift
		}

		if u.NextID+config.UIDShift > settings.MIN_LIMIT_ID {
			nextID = u.NextID + config.UIDShift
		}
		attr[u.Name] = &stns.Attribute{
			ID:     u.ID + config.UIDShift,
			PrevID: prevID,
			NextID: nextID,
			User:   tmpUser,
		}
	}
}

func gidShift(attr stns.Attributes, g *v3Group, config *Config) {
	if g.ID+config.GIDShift > settings.MIN_LIMIT_ID {
		prevID := 0
		nextID := 0
		tmpGroup := &stns.Group{
			Users: g.Users,
		}

		if g.PrevID+config.GIDShift > settings.MIN_LIMIT_ID {
			prevID = g.PrevID + config.GIDShift
		}

		if g.NextID+config.GIDShift > settings.MIN_LIMIT_ID {
			nextID = g.NextID + config.GIDShift
		}

		attr[g.Name] = &stns.Attribute{
			ID:     g.ID + config.GIDShift,
			PrevID: prevID,
			NextID: nextID,
			Group:  tmpGroup,
		}
	}
}

func convertV3Format(b []byte, path string, config *Config) (*ResponseFormat, error) {
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
				uidShift(attr, &u, config)
			}
		} else {
			u := v3User{}
			err = json.Unmarshal(b, &u)
			if err != nil {
				return nil, err
			}
			uidShift(attr, &u, config)
		}
	case "group":
		if sp[1] == "list" {
			groups := v3Groups{}
			err = json.Unmarshal(b, &groups)
			if err != nil {
				return nil, err
			}
			for _, g := range groups {
				gidShift(attr, &g, config)
			}
		} else {
			g := v3Group{}
			err = json.Unmarshal(b, &g)
			if err != nil {
				return nil, err
			}

			gidShift(attr, &g, config)
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

	return &ResponseFormat{
		attr,
	}, nil
}

type ResponseFormat struct {
	Items stns.Attributes `json:"items"`
}

func (r ResponseFormat) First() *stns.Attribute {
	for _, v := range r.Items {
		return v
	}
	return &stns.Attribute{}
}

type v2MetaData struct {
	APIVersion float64 `json:"api_version"`
	Result     string  `json:"result"`
}

type v2ResponseFormat struct {
	MetaData *v2MetaData     `json:"metadata"`
	Items    stns.Attributes `json:"items"`
}

type v3User struct {
	ID            int      `json:"id"`
	PrevID        int      `json:"prev_id"`
	NextID        int      `json:"next_id"`
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
	ID     int      `json:"id"`
	PrevID int      `json:"prev_id"`
	NextID int      `json:"next_id"`
	Name   string   `json:"name"`
	Users  []string `json:"users"`
}

type v3Sudo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type v3Users []v3User

type v3Groups []v3Group

type v3Sudoers []v3Sudo
