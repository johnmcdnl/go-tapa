package tapa

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Tapa struct {
}

func New() *Tapa {
	return &Tapa{}
}

func (t *Tapa) String() string {
	j, err := json.Marshal(t)
	if err != nil {
		panic(errors.Wrap(err, "failed to get a tapa.String()"))
	}
	return string(j)
}
