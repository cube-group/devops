package models

import "errors"

type Kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (t *Kv) Validator() error {
	if t.Key == "" || t.Value == "" {
		return errors.New("key or value is nil")
	}
	return nil
}
