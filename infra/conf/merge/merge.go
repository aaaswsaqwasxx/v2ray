package merge

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/cmdarg"
	"v2ray.com/core/infra/conf/serial"
)

// ToJSON merges multiple jsons into one.
// It accepts []string for URLs, files, [][]byte for json contents
func ToJSON(args interface{}) ([]byte, error) {
	m, err := ToMap(args)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// ToMap merges multiple jsons into one map.
// It accepts []string for URLs, files, [][]byte for json contents
func ToMap(args interface{}) (m map[string]interface{}, err error) {
	switch v := args.(type) {
	case cmdarg.Arg:
		m, err = filesToMap([]string(v))
	case []string:
		m, err = filesToMap(v)
	case [][]byte:
		m, err = bytesToMap(v)
	default:
		return nil, newError("unsupport args type")
	}
	if err != nil {
		return nil, err
	}
	sortSlicesInMap(m)
	err = mergeSameTag(m)
	if err != nil {
		return nil, err
	}
	removeHelperFields(m)
	return m, nil
}

func filesToMap(args []string) (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	for _, arg := range args {
		r, err := loadArg(arg)
		if err != nil {
			return nil, err
		}
		m, err := readerToMap(r)
		if err != nil {
			return nil, err
		}
		if err = mergeMaps(conf, m); err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func bytesToMap(args [][]byte) (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	for _, arg := range args {
		r := bytes.NewReader(arg)
		m, err := readerToMap(r)
		if err != nil {
			return nil, err
		}
		if err = mergeMaps(conf, m); err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func readerToMap(r io.Reader) (map[string]interface{}, error) {
	c := make(map[string]interface{})
	err := serial.DecodeJSON(r, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// loadArg loads one arg, maybe an remote url, or local file path
func loadArg(arg string) (out io.Reader, err error) {
	var data []byte
	switch {
	case strings.HasPrefix(arg, "http://"), strings.HasPrefix(arg, "https://"):
		data, err = fetchHTTPContent(arg)
	case (arg == "stdin:"):
		data, err = ioutil.ReadAll(os.Stdin)
	default:
		data, err = ioutil.ReadFile(arg)
	}
	if err != nil {
		return
	}
	out = bytes.NewBuffer(data)
	return
}

// fetchHTTPContent dials https for remote content
func fetchHTTPContent(target string) ([]byte, error) {
	parsedTarget, err := url.Parse(target)
	if err != nil {
		return nil, newError("invalid URL: ", target).Base(err)
	}

	if s := strings.ToLower(parsedTarget.Scheme); s != "http" && s != "https" {
		return nil, newError("invalid scheme: ", parsedTarget.Scheme)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL:    parsedTarget,
		Close:  true,
	})
	if err != nil {
		return nil, newError("failed to dial to ", target).Base(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError("unexpected HTTP status code: ", resp.StatusCode)
	}

	content, err := buf.ReadAllToBytes(resp.Body)
	if err != nil {
		return nil, newError("failed to read HTTP response").Base(err)
	}

	return content, nil
}
