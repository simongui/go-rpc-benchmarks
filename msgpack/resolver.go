package msgpack

import "reflect"

// Resolver Resolves types.
type Resolver map[string]reflect.Value

// Resolve Resolves a type from the specified name.
func (resolver Resolver) Resolve(name string, arguments []reflect.Value) (reflect.Value, error) {
	return resolver[name], nil
}
