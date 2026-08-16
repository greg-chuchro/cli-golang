package class_member_stores

import (
	"github.com/greg-chuchro/cli-golang/data_access"
)

// PropertyStore provides key-value access to object properties by name.
type PropertyStore struct {
	*PropertyStoreBase
}

func NewPropertyStore(obj any, bindingFlags BindingFlags) *PropertyStore {
	return &PropertyStore{PropertyStoreBase: NewPropertyStoreBase(obj, bindingFlags, nil)}
}

var _ data_access.IKeyValueStoreMediated[string, any, string, any] = (*PropertyStore)(nil)
