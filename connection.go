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
	"context"
	"github.com/go-wares/log"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// Connection
	// instance for mysql connection manager.
	Connection ConnectionManager

	defaultEngineGroup, _ = xorm.NewEngineGroup(DefaultDriver, []string{DefaultDsn})
)

const (
	DefaultConnectionName = "db"
)

type (
	// ConnectionManager
	// interface for mysql connection manager.
	ConnectionManager interface {
		// Master
		// return master connection session.
		Master(ctx context.Context) *xorm.Session

		// MasterWith
		// return master connection session with configuration name.
		MasterWith(ctx context.Context, name string) *xorm.Session

		// Slave
		// return slave connection session.
		Slave(ctx context.Context) *xorm.Session

		// SlaveWith
		// return slave connection session with configuration name.
		SlaveWith(ctx context.Context, name string) *xorm.Session
	}

	connection struct {
		configs map[string]*Config
		engines map[string]*xorm.EngineGroup
		mu      *sync.RWMutex
	}
)

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *connection) Master(ctx context.Context) *xorm.Session {
	return o.MasterWith(ctx, DefaultConnectionName)
}

func (o *connection) MasterWith(ctx context.Context, name string) *xorm.Session {
	return o.get(name).Master().NewSession().Context(ctx)
}

func (o *connection) Slave(ctx context.Context) *xorm.Session {
	return o.SlaveWith(ctx, DefaultConnectionName)
}

func (o *connection) SlaveWith(ctx context.Context, name string) *xorm.Session {
	return o.get(name).Slave().NewSession().Context(ctx)
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *connection) get(name string) (eg *xorm.EngineGroup) {
	var (
		c   *Config
		err error
		ok  bool
	)

	// Open lock resource then close when ended.
	o.mu.Lock()
	defer func() {
		o.mu.Unlock()

		// Build default engine group.
		if eg == nil {
			eg = defaultEngineGroup
		}
	}()

	// Return
	// xorm engine group.
	if eg, ok = o.engines[name]; ok {
		log.Debugf("[XORM] reuse xorm engine group: %s", name)
		return
	}

	// Return
	// if config name not configured.
	if c, ok = o.configs[name]; !ok {
		log.Warnf("[XORM] engine group config not specified: %s", name)
		return
	}

	// Return
	// if create new engine group error.
	if eg, err = xorm.NewEngineGroup(c.Driver, c.Dsn); err != nil {
		eg = nil
		log.Errorf("[XORM] engine group create error: name=%s, %v", name, err)
		return
	}

	// Generate engine options.
	eg.SetMapper(c.GetMapper())
	eg.SetMaxOpenConns(c.MaxOpen)
	eg.SetMaxIdleConns(c.MaxIdle)
	eg.SetConnMaxLifetime(c.maxLifetime)
	eg.SetLogger(&logger{c: c})
	eg.ShowSQL(c.ShowSql)
	eg.EnableSessionID(c.EnableSessionId)

	// Update engines mapper.
	o.engines[name] = eg
	return
}

func (o *connection) init() *connection {
	o.engines = make(map[string]*xorm.EngineGroup)
	o.mu = new(sync.RWMutex)
	return o.scan()
}

func (o *connection) scan() *connection {
	cs := make(map[string]*Config)

	for _, s := range []string{"config/db.yaml", "../config/db.yaml"} {
		if si, se := os.Stat(s); se == nil && !si.IsDir() {
			if rb, re := os.ReadFile(s); re == nil {
				if err := yaml.Unmarshal(rb, cs); err == nil {
					break
				}
			}
		}
	}

	for _, c := range cs {
		c.Default()
	}

	o.configs = cs
	return o
}
