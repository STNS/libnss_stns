package request

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	stns_settings "github.com/STNS/STNS/settings"
	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/logger"
	"github.com/STNS/libnss_stns/settings"
)

type Request struct {
	ApiPath string
	Config  *config.Config
}

func NewRequest(config *config.Config, paths ...string) (*Request, error) {
	logger.Setlog()
	r := Request{}

	r.Config = config
	r.SetPath(paths...)
	return &r, nil
}

func (r *Request) SetPath(paths ...string) {
	r.ApiPath = path.Clean(strings.Join(paths, "/"))
}

// only use wrapper command
func (r *Request) GetRawData() ([]byte, error) {
	var lastError error
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(r.Config.ApiEndPoint))

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: !r.Config.SslVerify}
	http.DefaultTransport.(*http.Transport).Dial = (&net.Dialer{
		Timeout:   settings.HTTP_TIMEOUT * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial

	for _, v := range perm {
		endPoint := r.Config.ApiEndPoint[v]
		url := strings.TrimRight(endPoint, "/") + "/" + strings.TrimLeft(r.ApiPath, "/")
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			lastError = err
			continue
		}

		if r.Config.User != "" && r.Config.Password != "" {
			req.SetBasicAuth(r.Config.User, r.Config.Password)
		}

		if r.checkLockFile(endPoint) {
			res, err := http.DefaultClient.Do(req)

			if err != nil {
				r.writeLockFile(endPoint)
				lastError = err
				continue
			}

			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				lastError = err
				continue
			}

			switch res.StatusCode {
			case http.StatusOK:
				reg := regexp.MustCompile(`/v2[/]?$`)
				switch {
				// version1
				case !reg.MatchString(endPoint):
					buffer, err := r.migrateV2Format(body)
					if err != nil {
						lastError = err
						continue
					}
					return buffer, nil
				default:
					return body, nil
				}
			// only direct return notfonud
			case http.StatusNotFound:
				return nil, fmt.Errorf("resource notfound: %s", url)
			case http.StatusUnauthorized:
				return nil, fmt.Errorf("authenticate error: %s", url)
			default:
				continue
			}
		}
	}
	return nil, lastError
}

func (r *Request) migrateV2Format(body []byte) ([]byte, error) {
	var attr stns.Attributes
	err := json.Unmarshal(body, &attr)

	if err != nil {
		return nil, err
	}

	if attr == nil {
		return nil, errors.New(settings.V2_FORMAT_ERROR)
	}

	mig := stns.ResponseFormat{
		&stns.MetaData{
			1.0,
			false,
			0,
			stns_settings.SUCCESS,
		},
		&attr,
	}

	j, err := json.Marshal(mig)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func (r *Request) checkLockFile(endPoint string) bool {
	fileName := "/tmp/libnss_stns." + r.GetMD5Hash(endPoint)
	_, err := os.Stat(fileName)

	// lockfile not exists
	if err != nil {
		return true
	}

	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
		os.Remove(fileName)
		return false
	}

	buf := bytes.NewBuffer(buff)
	lastTime, err := binary.ReadVarint(buf)
	if err != nil {
		log.Println(err)
		os.Remove(fileName)
		return false
	}

	if time.Now().Unix() > lastTime+settings.LOCK_TIME || lastTime > time.Now().Unix()+settings.LOCK_TIME {
		os.Remove(fileName)
		return true
	}

	return false
}

func (r *Request) writeLockFile(endPoint string) {
	fileName := "/tmp/libnss_stns." + r.GetMD5Hash(endPoint)
	now := time.Now()
	log.Println("create lockfile:" + endPoint + " time:" + now.String() + " unix_time:" + strconv.FormatInt(now.Unix(), 10))

	result := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(result, now.Unix())
	ioutil.WriteFile(fileName, result, os.ModePerm)
}

// only use wrapper command
func (r *Request) GetAttributes() (stns.Attributes, error) {
	var res stns.ResponseFormat

	body, err := r.GetRawData()

	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &res)

		if err != nil {
			return nil, err
		}
	}

	return *res.Items, nil
}

func (r *Request) GetByWrapperCmd() (stns.ResponseFormat, error) {
	out, err := exec.Command(r.Config.WrapperCommand, r.ApiPath).Output()
	if err != nil {
		return stns.ResponseFormat{}, err
	}
	var res stns.ResponseFormat
	err = json.Unmarshal(out, &res)
	if err != nil {
		return stns.ResponseFormat{}, err
	}
	return res, nil
}

func (r *Request) GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
