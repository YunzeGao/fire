package formatter

import (
	"bytes"
	"encoding/json"
	"github.com/YunzeGao/fire/framework/util"
	"time"

	"github.com/YunzeGao/fire/framework/contract"

	"github.com/pkg/errors"
)

func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	fields["msg"] = msg
	fields["level"] = level
	fields["timestamp"] = util.DefaultTimeFormatMicroSeconds(t)
	content, err := json.Marshal(fields)
	if err != nil {
		return bf.Bytes(), errors.Wrap(err, "json format error")
	}
	bf.Write(content)
	return bf.Bytes(), nil
}
