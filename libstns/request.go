package libstns

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/context"

	stns_settings "github.com/STNS/STNS/settings"
	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/cache"
	"github.com/STNS/libnss_stns/settings"
)

type Request struct {
	ApiPath string
	Config  *Config
}

func NewRequest(config *Config, paths ...string) (*Request, error) {
	r := Request{
		ApiPath: path.Clean(strings.Join(paths, "/")),
		Config:  config,
	}
	return &r, nil
}

// only use wrapper command
func (r *Request) GetRawData() ([]byte, error) {
	if len(r.Config.ApiEndPoint) == 0 {
		return nil, errors.New("endpoint not defined")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rch := make(chan []byte, len(r.Config.ApiEndPoint))
	ech := make(chan error, len(r.Config.ApiEndPoint))
	for _, e := range r.Config.ApiEndPoint {
		go func(endPoint string) {
			if cache.IsLockEndPoint(endPoint) {
				ech <- fmt.Errorf("endpoint %s is locked", endPoint)
				return
			}

			u := strings.TrimRight(endPoint, "/") + "/" + strings.TrimLeft(r.ApiPath, "/")
			req, err := http.NewRequest("GET", u, nil)

			for k, v := range r.Config.RequestHeader {
				req.Header.Add(k, v)
			}

			if err != nil {
				ech <- err
				return
			}

			if r.Config.User != "" && r.Config.Password != "" {
				req.SetBasicAuth(r.Config.User, r.Config.Password)
			}

			r.httpDo(
				ctx,
				req,
				func(res *http.Response, err error) {
					if err != nil {
						if _, ok := err.(*url.Error); ok {
							cache.LockEndPoint(endPoint)
						}
						ech <- err
						return
					}
					defer res.Body.Close()
					body, err := ioutil.ReadAll(res.Body)
					switch res.StatusCode {
					case http.StatusOK, http.StatusNotFound:
						reg := regexp.MustCompile(`/v2[/]?$`)
						switch {
						// version1
						case !reg.MatchString(endPoint):
							buffer, err := r.migrateV2Format(body)
							if err != nil {
								ech <- err
								return
							}
							rch <- buffer
							return
						default:
							rch <- body
							return
						}
					case http.StatusUnauthorized:
						ech <- fmt.Errorf("authenticate error: %s", u)
						return
					default:
						ech <- fmt.Errorf("error: %s", u)
						return
					}
				},
			)
		}(e)
	}

	var cnt int
	for {
		select {
		case r := <-rch:
			return r, nil
		case e := <-ech:
			cnt++
			if cnt == len(r.Config.ApiEndPoint) {
				return nil, e
			}
		}
	}
}

func (r *Request) httpDo(
	ctx context.Context,
	req *http.Request,
	f func(*http.Response, error),
) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !r.Config.SslVerify},
		Dial: (&net.Dialer{
			Timeout:   time.Duration(r.Config.RequestTimeOut) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	tr.Proxy = http.ProxyFromEnvironment
	if r.Config.HttpProxy != "" {
		proxyUrl, err := url.Parse(r.Config.HttpProxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	client := &http.Client{Transport: tr}

	go func() { f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		return
	}
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
			"sha256",
			stns_settings.SUCCESS,
			0,
		},
		&attr,
	}

	j, err := json.Marshal(mig)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func (r *Request) GetByWrapperCmd() (stns.ResponseFormat, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(r.Config.WrapperCommand, r.ApiPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return stns.ResponseFormat{}, err
	}

	if len(stderr.Bytes()) > 0 {
		return stns.ResponseFormat{}, fmt.Errorf("command error:%s", stderr.String())
	}

	var res stns.ResponseFormat
	err = json.Unmarshal(stdout.Bytes(), &res)
	if err != nil {
		return stns.ResponseFormat{}, err
	}
	return res, nil
}
