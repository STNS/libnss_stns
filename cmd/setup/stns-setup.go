package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/STNS/libnss_stns/libstns"
)

func main() {
	config, err := libstns.LoadConfig("/etc/stns/libnss_stns.conf")
	if err == nil {
		if err := Run(config, os.Getenv("USER")); err != nil {
			fmt.Println(err)
		}
	}
}

func Run(config *libstns.Config, name string) error {
	var res *libstns.ResponseFormat

	r, err := libstns.NewRequest(config, "user", "name", name)
	if err != nil {
		return err
	}

	body, err := r.GetRawData()
	if err != nil {
		return err
	}
	if len(body) > 0 {
		err = json.Unmarshal(body, &res)

		if err != nil {
			return err
		}

		if res.Items != nil {
			for _, u := range res.Items {
				if len(u.SetupCommands) > 0 {
					for _, c := range u.SetupCommands {
						cmd := strings.Split(c, " ")
						fmt.Printf("run_command: %s\n", c)
						out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()

						if out != nil {
							fmt.Println(string(out))
						}

						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
