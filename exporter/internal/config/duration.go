package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unexpected type: %T", value)
	}
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	var err error

	d.Duration, err = time.ParseDuration(node.Value)
	if err != nil {
		val, terr := strconv.ParseFloat(node.Value, 64)
		if terr != nil {
			return fmt.Errorf("error parsing duration: [%w, %w]", terr, err)
		}

		d.Duration = time.Duration(val)
	}

	return nil
}
