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
	"fmt"
	"github.com/go-wares/log"
	"sync"
	"xorm.io/xorm"
)

var (
	// Transaction instance pool.
	transactionPool sync.Pool
)

type (
	// TransactionHandler
	// is a process handler of transaction.
	TransactionHandler func(ctx context.Context, session *xorm.Session) (err error)

	// TransactionManager
	// interface for transaction manager.
	TransactionManager interface {
		Add(handlers ...TransactionHandler) TransactionManager
		Release()
		Run(ctx context.Context) error
		WithSession(session *xorm.Session) TransactionManager
	}

	transaction struct {
		handlers                      []TransactionHandler
		session                       *xorm.Session
		sessionCreated, sessionOpened bool
	}
)

// Transaction
// return a transaction manager instance.
func Transaction() TransactionManager {
	// Acquire an instance from pool.
	if v := transactionPool.Get(); v != nil {
		return v.(*transaction).acquire()
	}

	// Create an instance if null got.
	return (&transaction{}).init().acquire()
}

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

// Add
// handlers into transaction process.
func (o *transaction) Add(handlers ...TransactionHandler) TransactionManager {
	o.handlers = append(o.handlers, handlers...)
	return o
}

// Release
// instance into pool.
func (o *transaction) Release() {
	transactionPool.Put(o.clean())
}

// Run
// transaction process.
func (o *transaction) Run(ctx context.Context) (err error) {
	// Call
	// when process ended.
	defer func() {
		// Catch
		// runtime panic.
		if r := recover(); r != nil {
			log.Fatalfc(ctx, "[XORM-TX] process fatal: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", r)
			}
		}

		// Send transaction
		// if session opened.
		if o.sessionOpened {
			if err != nil {
				if te := o.session.Rollback(); te != nil {
					log.Errorfc(ctx, "[XORM][TX] process rollback: %v", te)
				} else {
					log.Debugfc(ctx, "[XORM][TX] process rollback completed")
				}
			} else {
				if te := o.session.Commit(); te != nil {
					log.Errorfc(ctx, "[XORM][TX] process commit: %v", te)
				} else {
					log.Debugfc(ctx, "[XORM][TX] process commit")
				}
			}
		}

		// Close connection session
		// if created in runtime.
		if o.sessionCreated {
			if te := o.session.Close(); te != nil {
				log.Errorfc(ctx, "[XORM][TX] session close: %v", te)
			} else {
				log.Infofc(ctx, "[XORM][TX] session closed")
			}
		}
	}()

	// Create connection session
	// if not specified.
	if o.session == nil {
		log.Debugfc(ctx, "[XORM][TX] session created")
		o.session = Connection.Master(ctx)
		o.sessionCreated = true
	}

	// Transaction opener.
	if err = o.session.Begin(); err != nil {
		return
	}

	// Transaction opened.
	o.sessionOpened = true
	log.Debugfc(ctx, "[XORM][TX] session opened")

	// Range transaction handlers.
	for _, handler := range o.handlers {
		if err = handler(ctx, o.session); err != nil {
			break
		}
	}
	return
}

// WithSession
// bind xorm session on transaction process.
func (o *transaction) WithSession(session *xorm.Session) TransactionManager {
	o.session = session
	return o
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *transaction) acquire() *transaction {
	o.handlers = make([]TransactionHandler, 0)
	o.sessionCreated = false
	o.sessionOpened = false
	return o
}

func (o *transaction) clean() *transaction {
	o.handlers = nil
	o.session = nil
	return o
}

func (o *transaction) init() *transaction {
	return o
}
