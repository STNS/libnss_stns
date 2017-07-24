package libstns

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/STNS/libnss_stns/cache"
)

type Request struct {
	ApiPath      string
	Config       *Config
	ResourceType string
}

func NewRequest(config *Config, paths ...string) (*Request, error) {
	if len(paths) == 3 && strings.Contains(paths[2], ".") {
		paths[2] = urlencode(paths[2])
	}

	r := Request{
		ApiPath:      path.Clean(strings.Join(paths, "/")),
		Config:       config,
		ResourceType: paths[0],
	}
	return &r, nil
}

// only use wrapper command
func (r *Request) GetRawData() ([]byte, error) {
	var b []byte
	var e error

	if len(r.Config.ApiEndPoint) == 0 {
		return nil, errors.New("endpoint not defined")
	}

	retry := 3
	if r.Config.RequestRetry != 0 {
		retry = r.Config.RequestRetry
	}

	for i := 0; i < retry; i++ {
		b, e = r.request()
		if e == nil {
			break
		}
	}
	return b, e
}

func (r *Request) request() ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rch := make(chan *ResponseFormat, len(r.Config.ApiEndPoint))
	ech := make(chan error, len(r.Config.ApiEndPoint))
	for _, e := range r.Config.ApiEndPoint {
		go func(endPoint string) {
			if cache.IsLockEndPoint(endPoint) {
				ech <- fmt.Errorf("endpoint %s is locked", endPoint)
				return
			}

			u := strings.TrimRight(endPoint, "/") + "/" + strings.TrimLeft(r.ApiPath, "/")
			req, err := http.NewRequest("GET", u, nil)
			if err != nil {
				ech <- err
				return
			}

			for k, v := range r.Config.RequestHeader {
				req.Header.Add(k, v)
			}

			if r.Config.User != "" && r.Config.Password != "" {
				req.SetBasicAuth(r.Config.User, r.Config.Password)
			}

			r.httpDo(
				ctx,
				req,
				func(res *http.Response, err error) {
					if err != nil {
						if _, ok := err.(*url.Error); ok && len(r.Config.ApiEndPoint) != 1 {
							cache.LockEndPoint(endPoint)
						}
						ech <- err
						return
					}

					defer res.Body.Close()
					body, err := ioutil.ReadAll(res.Body)

					switch res.StatusCode {
					case http.StatusOK:
						v2 := regexp.MustCompile(`/v2[/]?$`)
						v3 := regexp.MustCompile(`/v3[/]?$`)
						switch {
						// version1
						case !v2.MatchString(endPoint) && !v3.MatchString(endPoint):
							buffer, err := convertV1toV3Format(body)
							if err != nil {
								ech <- err
								return
							}
							rch <- buffer
							return
						// version2
						case v2.MatchString(endPoint):
							buffer, err := convertV2toV3Format(body)
							if err != nil {
								ech <- err
								return
							}
							rch <- buffer
							return
						default:
							buffer, err := convertV3Format(body, r.ApiPath, r.Config)
							if err != nil {
								ech <- err
								return
							}
							rch <- buffer
							return
						}
					case http.StatusNotFound:
						ids := map[string]int{}
						for _, t := range []string{"Min", "Max"} {
							id := res.Header.Get(fmt.Sprintf("Stns-%s-Id", t))
							if id != "" {
								i, err := strconv.Atoi(id)
								if err != nil {
									ech <- err
									return
								}
								ids[t] = i
							}
						}

						if len(ids) > 0 {
							ech <- fmt.Errorf("resource not found min_id: %v max_id %v url: %s", ids["Min"], ids["Max"], u)
						} else {
							ech <- fmt.Errorf("resource not found url: %s", u)
						}
						return
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
			j, err := json.Marshal(r)
			if err != nil {
				return nil, err
			}
			return j, nil
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
	tc := r.TlsConfig()
	tr := &http.Transport{
		TLSClientConfig: tc,
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

func (r *Request) TlsConfig() *tls.Config {
	tc := &tls.Config{InsecureSkipVerify: !r.Config.SslVerify}

	if r.TlsKeysExists() {
		cert, err := tls.LoadX509KeyPair(r.Config.TlsCert, r.Config.TlsKey)
		if err != nil {
			log.Println(err)
			goto ret
		}

		if _, err := os.Stat(r.Config.TlsCa); err == nil {
			// Load CA cert
			caCert, err := ioutil.ReadFile(r.Config.TlsCa)
			if err != nil {
				log.Println(err)
				goto ret
			}
			caPool := x509.NewCertPool()
			caPool.AppendCertsFromPEM(caCert)

			tc.Certificates = []tls.Certificate{cert}
			tc.RootCAs = caPool

			tc.BuildNameToCertificate()
		}

	}
ret:
	return tc
}

func (r *Request) TlsKeysExists() bool {
	if r.Config.TlsCert != "" && r.Config.TlsKey != "" {
		for _, v := range []string{r.Config.TlsCert, r.Config.TlsKey} {
			if _, err := os.Stat(v); err != nil {
				log.Println(err)
				return false
			}
		}
		return true
	}
	return false
}

func (r *Request) GetByWrapperCmd() (*ResponseFormat, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(r.Config.WrapperCommand, r.ApiPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	if len(stderr.Bytes()) > 0 {
		reg := regexp.MustCompile(`resource not found min_id: ([\d]+) max_id ([\d]+) url: .*`)
		if result := reg.FindStringSubmatch(stderr.String()); result != nil {
			for index, t := range []string{"min", "max"} {
				i, err := strconv.Atoi(string(result[index+1]))
				if err != nil {
					return nil, fmt.Errorf("command error:%s", err)
				}
				cache.WriteID(r.ResourceType, t, i)
			}
		} else {

		}
		return nil, fmt.Errorf("command error:%s", stderr.String())
	}

	var res *ResponseFormat
	if len(stdout.Bytes()) > 0 {
		err = json.Unmarshal(stdout.Bytes(), &res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func urlencode(s string) (result string) {
	for _, c := range s {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%%%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%%%X%%%X%%%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%%%X%%%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}

	return result
}
