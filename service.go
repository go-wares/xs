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
	"context"
	"xorm.io/xorm"
)

type (
	// Interface
	// for top level service.
	Interface interface {
		Clean()
		Master(ctx context.Context) *xorm.Session
		MasterWith(ctx context.Context, name string) *xorm.Session
		Slave(ctx context.Context) *xorm.Session
		SlaveWith(ctx context.Context, name string) *xorm.Session
		With(opts ...Option)
	}

	// Service
	// for top level configurations.
	Service struct{ session *xorm.Session }
)

func (o *Service) Clean() {
	if o.session != nil {
		o.session = nil
	}
}

func (o *Service) Master(ctx context.Context) *xorm.Session {
	return Connection.Master(ctx)
}

func (o *Service) MasterWith(ctx context.Context, name string) *xorm.Session {
	return Connection.MasterWith(ctx, name)
}

func (o *Service) Slave(ctx context.Context) *xorm.Session {
	return Connection.Slave(ctx)
}

func (o *Service) SlaveWith(ctx context.Context, name string) *xorm.Session {
	return Connection.SlaveWith(ctx, name)
}

func (o *Service) With(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}
