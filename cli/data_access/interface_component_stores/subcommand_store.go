package interface_component_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/data_access/class_member_stores"
	"github.com/greg-chuchro/cli-golang/cli/interface_components"
)

// SubcommandStore provides read-only access to CLI subcommands (methods) with metadata.
type SubcommandStore struct {
	store *class_member_stores.MethodInfoStore
}

func NewSubcommandStore(store *class_member_stores.MethodInfoStore) *SubcommandStore {
	return &SubcommandStore{store: store}
}

func (this *SubcommandStore) Get(key string) reflect.Method {
	return this.store.Get(key)
}

func (this *SubcommandStore) ContainsKey(key string) bool {
	return this.store.ContainsKey(key)
}

func (this *SubcommandStore) GetValueType(key string) reflect.Type {
	return this.store.GetValueType(key)
}

func (this *SubcommandStore) GetSubcommands(createSubcommandExecutor func(reflect.Method, any) ISubcommandExecutor) []interface_components.Subcommand {
	result := make([]interface_components.Subcommand, 0)
	for _, pair := range this.store.GetAccessorKeysPairs() {
		method := this.store.Get(pair.Keys[0])
		executor := createSubcommandExecutor(method, nil)
		result = append(result, interface_components.Subcommand{
			Keys:       pair.Keys,
			MethodInfo: method,
			ReturnType: this.GetValueType(pair.Keys[0]),
			Parameters: executor.GetParameters(),
		})
	}
	return result
}
