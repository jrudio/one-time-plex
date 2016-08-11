package ledis

import (
	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

// Open initializes LedisDB and connects to it's instance
func Open() (*ledis.DB, error) {
	cfg := config.NewConfigDefault()

	cfg.DataDir = "./otp-data"

	lds, err := ledis.Open(cfg)

	if err != nil {
		return &ledis.DB{}, err
	}

	return lds.Select(0)
}
