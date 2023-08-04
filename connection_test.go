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
	"github.com/go-wares/log"
	"testing"
)

type (
	testModel struct {
		Id int
	}

	testService struct {
		Service
	}
)

func (o *testModel) TableName() string {
	return "user"
}

func newTestService(opts ...Option) *testService {
	v := &testService{}
	v.With(opts...)
	return v
}

func (o *testService) Get(ctx context.Context) {
	bean := &testModel{
		Id: 1,
	}
	o.Master(ctx).Get(bean)
}

func TestConnection(t *testing.T) {
	defer log.Stop()
	span := log.NewSpan("testing")
	Connection.Master(span.Context()).Query("SELECT 1")
}

func TestService(t *testing.T) {
	span := log.NewSpan("test-service")
	sess := Connection.Master(span.Context())

	s := newTestService(WithSession(sess))
	s.Get(span.Context())
}
