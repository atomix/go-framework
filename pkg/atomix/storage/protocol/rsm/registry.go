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

package rsm

// NewServiceFunc is a function for creating a primitive service
type NewServiceFunc func(scheduler Scheduler, context ServiceContext) Service

// Registry is the default primitive registry
type Registry struct {
	services map[string]NewServiceFunc
}

func (r *Registry) Register(primitiveType string, primitive NewServiceFunc) {
	r.services[primitiveType] = primitive
}

func (r *Registry) GetServices() []NewServiceFunc {
	services := make([]NewServiceFunc, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}
	return services
}

func (r *Registry) GetService(primitiveType string) NewServiceFunc {
	return r.services[primitiveType]
}

// NewRegistry creates a new primitive registry
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]NewServiceFunc),
	}
}
