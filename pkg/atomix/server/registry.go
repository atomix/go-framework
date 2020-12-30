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

package server

import (
	"google.golang.org/grpc"
)

// RegisterServiceFunc is a function for registering a service
type RegisterServiceFunc func(server *grpc.Server)

// Registry is a primitive registry
type Registry interface {
	// Register registers a service
	Register(service RegisterServiceFunc)

	// GetPrimitives gets a list of primitives
	GetServices() []RegisterServiceFunc
}

// primitiveRegistry is the default primitive registry
type primitiveRegistry struct {
	services []RegisterServiceFunc
}

func (r *primitiveRegistry) Register(service RegisterServiceFunc) {
	r.services = append(r.services, service)
}

func (r *primitiveRegistry) GetServices() []RegisterServiceFunc {
	primitives := make([]RegisterServiceFunc, 0, len(r.services))
	for _, primitive := range r.services {
		primitives = append(primitives, primitive)
	}
	return primitives
}

// NewRegistry creates a new primitive registry
func NewRegistry() Registry {
	return &primitiveRegistry{
		services: make([]RegisterServiceFunc, 0),
	}
}
