package consul

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
)

var kv *api.KV

// New initialize consul API
// ex : https://user:pass@consul.okadoc.co/dc1
func New(val string) error {
	// return if kv already set
	if kv != nil {
		return nil
	}

	u, err := url.Parse(val)
	if err != nil {
		return err
	}

	pwd, _ := u.User.Password()

	client, err := api.NewClient(&api.Config{
		Address:    u.Host,
		Scheme:     u.Scheme,
		Datacenter: strings.Replace(u.Path, "/", "", -1),

		// if basic auth exists
		HttpAuth: &api.HttpBasicAuth{
			Username: u.User.Username(),
			Password: pwd,
		},
	})

	if err != nil {
		return err
	}

	// set kv into global variable
	kv = client.KV()
	return nil
}

// LoadENV by prefix
func LoadENV(url, prefix string) error {
	err := New(url)
	if err != nil {
		return err
	}

	prefix = strings.ToUpper(prefix)
	pairs, err := getAllKV(prefix)
	if err != nil {
		return err
	}

	setENV(prefix, pairs)
	return nil
}

// KVPut add key/value into consul KV
// accept map[string]interface{} returns errors
func KVPut(p map[string]interface{}) error {
	var kvpairs api.KVPairs
	for k, v := range p {
		kvpairs = append(kvpairs, &api.KVPair{
			Key:   k,
			Value: []byte(fmt.Sprintf("%v", v)),
		})
	}

	for _, x := range kvpairs {
		_, err := kv.Put(x, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// KVGet accept string and return empty string if not found
func KVGet(key string) string {
	pair, _, err := kv.Get(key, &api.QueryOptions{
		UseCache: false,
	})

	if err != nil {
		return ""
	}

	return string(pair.Value)
}

func getAllKV(prefix string) (api.KVPairs, error) {
	kvpairs, _, err := kv.List(prefix, nil)

	if err != nil {
		return nil, err
	}

	return kvpairs, nil
}

func setENV(prefix string, lists api.KVPairs) {
	for _, v := range lists {
		key := strings.Replace(v.Key, prefix+"/", "", -1)
		val := strings.Trim(string(v.Value), `"`)
		os.Setenv(key, val)
	}
}
