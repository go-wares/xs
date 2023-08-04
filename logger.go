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

package service

import (
	"github.com/go-wares/log"
	l "xorm.io/xorm/log"
)

type (
	// logger
	// adapter for xorm package.
	logger struct{ c *Config }
)

// AfterSQL
// send log message after query.
func (o *logger) AfterSQL(c l.LogContext) {
	var (
		spa, exists = log.SpanExists(c.Ctx)
		format      = "[XORM] %s, username=%s, schema=%s, duration=%dms"
		arguments   = []interface{}{c.SQL, o.c.username, o.c.schema, c.ExecuteTime.Milliseconds()}
	)

	// Append arguments
	// on logger.
	if len(c.Args) > 0 {
		format += ", arguments=%v"
		arguments = append(arguments, c.Args)
	}

	// Send info log.
	if exists {
		spa.Info(format, arguments...)
	} else {
		log.Infof(format, arguments...)
	}

	// Send error log
	// if mistake occurred.
	if c.Err != nil {
		if exists {
			spa.Error("[XORM] %v", c.Err)
		} else {
			log.Errorf("[XORM] %v", c.Err)
		}
	}
}

// +---------------------------------------------------------------------------+
// | Not used methods                                                          |
// +---------------------------------------------------------------------------+

func (o *logger) BeforeSQL(_ l.LogContext)          {}
func (o *logger) Debugf(_ string, _ ...interface{}) {}
func (o *logger) Infof(_ string, _ ...interface{})  {}
func (o *logger) Warnf(_ string, _ ...interface{})  {}
func (o *logger) Errorf(_ string, _ ...interface{}) {}
func (o *logger) Level() l.LogLevel                 { return l.LOG_INFO }
func (o *logger) SetLevel(_ l.LogLevel)             {}
func (o *logger) ShowSQL(_ ...bool)                 {}
func (o *logger) IsShowSQL() bool                   { return true }
