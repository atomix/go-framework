// Copyright 2019-present Open Networking Foundation.
//
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

package counter

import (
	"github.com/atomix/atomix-api/go/atomix/primitive/counter"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
)

func init() {
	registerServiceFunc(newService)
}

func newService(scheduler rsm.Scheduler, context rsm.ServiceContext) Service {
	return &counterService{
		Service: rsm.NewService(scheduler, context),
	}
}

// counterService is a state machine for a counter primitive
type counterService struct {
	rsm.Service
	value int64
}

func (c *counterService) Set(input *counter.SetRequest) (*counter.SetResponse, error) {
	if err := checkPreconditions(c.value, input.Preconditions); err != nil {
		return nil, err
	}
	c.value = input.Value
	return &counter.SetResponse{
		Value: c.value,
	}, nil
}

func (c *counterService) Get(input *counter.GetRequest) (*counter.GetResponse, error) {
	return &counter.GetResponse{
		Value: c.value,
	}, nil
}

func (c *counterService) Increment(input *counter.IncrementRequest) (*counter.IncrementResponse, error) {
	c.value += input.Delta
	return &counter.IncrementResponse{
		Value: c.value,
	}, nil
}

func (c *counterService) Decrement(input *counter.DecrementRequest) (*counter.DecrementResponse, error) {
	c.value -= input.Delta
	return &counter.DecrementResponse{
		Value: c.value,
	}, nil
}

func checkPreconditions(value int64, preconditions []counter.Precondition) error {
	for _, precondition := range preconditions {
		switch p := precondition.Precondition.(type) {
		case *counter.Precondition_Value:
			if value != p.Value {
				return errors.NewConflict("value precondition failed")
			}
		}
	}
	return nil
}