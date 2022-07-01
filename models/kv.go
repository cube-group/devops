package models

import "errors"

type Kv struct {
	K   string `json:"key"`
	V string `json:"value"`
}

func (t *Kv) Validator() error {
	if t.K == "" || t.V == "" {
		return errors.New("key or value is nil")
	}
	return nil
}
