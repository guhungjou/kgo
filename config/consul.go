package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
)

var (
	prefix string
	kv     *api.KV
)

func ConsulInit(_prefix string) error {
	prefix = _prefix
	cfg := api.DefaultConfig()
	cfg.Token = os.Getenv("CONSUL_TOKEN")
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	kv = client.KV()
	return nil
}

func ConsulUnmarshal(k string, v interface{}) error {
	key := fmt.Sprintf("%s/%s", prefix, strings.Trim(k, "/"))
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	} else if pair == nil {
		return fmt.Errorf("key %s not found", k)
	}
	if err := json.Unmarshal(pair.Value, v); err != nil {
		return err
	}
	return nil
}
