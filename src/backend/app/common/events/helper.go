package events

import (
	"encoding/json"

	"github.com/growerlab/growerlab/src/backend/app/common/errors"
	"github.com/growerlab/growerlab/src/backend/app/common/mq"
)

func async(name, field string, t interface{}) error {
	body, err := json.Marshal(t)
	if err != nil {
		return errors.Trace(err)
	}
	_, err = MQ.Add(name, field, string(body))
	return err
}

func getPayload(pd *mq.Payload, fd string, out interface{}) error {
	if v := pd.Get(fd); v != nil {
		raw := []byte(v.(string))
		if err := json.Unmarshal(raw, out); err != nil {
			return err
		}
	}
	return nil
}
