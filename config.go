// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-08-04

package xs

import (
	"regexp"
	"strings"
	"time"
	"xorm.io/xorm/names"
)

const (
	DefaultDriver      = "mysql"
	DefaultDsn         = "username:password@tcp(127.0.0.1:3306)/undefined?charset=utf8"
	DefaultMapper      = "snake"
	DefaultMaxIdle     = 2
	DefaultMaxLifetime = 60
	DefaultMaxOpen     = 50
)

var (
	regexpConfigDsn = regexp.MustCompile(`^([_a-zA-Z0-9-.]+):(\S+)@tcp\(([^)]+)\)/([_a-zA-Z0-9-.]+)`)
)

type (
	// Config
	// for mysql connection configurations.
	Config struct {
		Driver          string   `json:"driver" yaml:"driver"`
		Dsn             []string `json:"dsn" yaml:"dsn"`
		EnableSessionId bool     `json:"enable_session_id" yaml:"enable_session_id"`
		ShowSql         bool     `json:"show_sql" yaml:"show_sql"`
		Mapper          string   `json:"mapper" yaml:"mapper"`
		MaxIdle         int      `json:"max_idle" yaml:"max_idle"`
		MaxLifetime     int      `json:"max_lifetime" yaml:"max_lifetime"`
		MaxOpen         int      `json:"max_open" yaml:"max_open"`

		maxLifetime                time.Duration
		mapper                     names.Mapper
		username, hostname, schema string
	}
)

// Default
// used to generate fields.
func (o *Config) Default() *Config {
	// Assign driver field
	// with default.
	if o.Driver == "" {
		o.Driver = DefaultDriver
	}

	// Assign mapper name
	// with default.
	if o.Mapper == "" {
		o.Mapper = DefaultMapper
	}

	// Generate mapper name as lower string.
	o.Mapper = strings.ToLower(o.Mapper)

	// Assign max idle connections
	// with default.
	if o.MaxIdle == 0 {
		o.MaxIdle = DefaultMaxIdle
	}

	// Assign max lifetime of connection
	// with default.
	if o.MaxLifetime == 0 {
		o.MaxLifetime = DefaultMaxLifetime
	}

	// Assign max open files
	// with default.
	if o.MaxOpen == 0 {
		o.MaxOpen = DefaultMaxOpen
	}

	// Executions.
	o.maxLifetime = time.Duration(o.MaxLifetime) * time.Second
	o.parseConnectionMapper()
	o.parseDatabaseSourceName()
	return o
}

// GetMapper
// return xorm mapper name.
func (o *Config) GetMapper() names.Mapper { return o.mapper }

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Config) parseConnectionMapper() {
	switch o.Mapper {
	case "snake":
		o.mapper = names.SnakeMapper{}
	case "gonic":
		o.mapper = names.GonicMapper{}
	default:
		o.mapper = names.SameMapper{}
	}
}

func (o *Config) parseDatabaseSourceName() {
	for _, s := range o.Dsn {
		if s = strings.TrimSpace(s); s != "" {
			if m := regexpConfigDsn.FindStringSubmatch(s); len(m) == 5 {
				o.username = m[1]
				o.hostname = m[3]
				o.schema = m[4]
			}
		}
		break
	}
}
