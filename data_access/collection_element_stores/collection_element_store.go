package collection_element_stores

import (
	"github.com/greg-chuchro/cli-golang/data_access"
)

// ICollectionElementStore provides key-value access to collection elements with collection-specific operations.
type ICollectionElementStore[TKey comparable, TValue any] interface {
	data_access.IKeyValueStore[TKey, TValue]
	Add(value TValue)
	Append(value TValue) any
	AddNew() TValue
	Remove(key TKey)
}
