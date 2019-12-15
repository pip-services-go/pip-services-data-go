package persistence

import (
	"github.com/pip-services3-go/pip-services3-commons-go/config"
	"reflect"
)

/*
Abstract persistence component that stores data in flat files
and implements a number of CRUD operations over data items with unique ids.
The data items must implement
 IIdentifiable interface

In basic scenarios child classes shall only override GetPageByFilter,
GetListByFilter or DeleteByFilter operations with specific filter function.
All other operations can be used out of the box.

In complex scenarios child classes can implement additional operations by
accessing cached items via IdentifiableFilePersistence._items property and calling Save method
on updates.

See JsonFilePersister
See MemoryPersistence

Configuration parameters

- path:                    path to the file where data is stored
- options:
    - max_page_size:       Maximum number of items returned in a single page (default: 100)

 References

- *:logger:*:*:1.0      (optional)  ILogger components to pass log messages

 Examples
type MyFilePersistence  struct {
	IdentifiableFilePersistence
}
    func NewMyFilePersistence(path string)(mfp *MyFilePersistence) {
		mfp = MyFilePersistence{}
		mfp.IdentifiableFilePersistence = *NewJsonPersister(path)
		return mfp
    }

    func composeFilter(filter cdata.FilterParams)(func (item interface{})bool) {
		if &filter == nil {
			filter = NewFilterParams()
		}
        name := filter.GetAsNullableString("name");
        return func (item interface) bool {
            dummy, ok := item.(Dummy)
			if *name != "" && ok && dummy.Name != *name {
				return false
			}
            return true
        }
    }

    func (c *MyFilePersistence ) GetPageByFilter(correlationId string, filter FilterParams, paging PagingParams)(page cdata.DataPage, err error){
        return c.GetPageByFilter(correlationId, this.composeFilter(filter), paging, nil, nil)
    }

    persistence := NewMyFilePersistence("./data/data.json")

	_, errc := persistence.Create("123", { Id: "1", Name: "ABC" })
	if (errc != nil) {
		panic()
	}
    page, errg := persistence.GetPageByFilter("123", NewFilterParamsFromTuples("Name", "ABC"), nil)
    if errg != nil {
		panic("Error")
	}
    fmt.Println(page.Data)         // Result: { Id: "1", Name: "ABC" )
    persistence.DeleteById("123", "1")
*/
type IdentifiableFilePersistence struct {
	IdentifiableMemoryPersistence
	_persister JsonFilePersister
}

// Creates a new instance of the persistence.
// - persister    (optional) a persister component that loads and saves data from/to flat file.
// Return *IdentifiableFilePersistence
// pointer on new IdentifiableFilePersistence
func NewIdentifiableFilePersistence(prototype reflect.Type, persister JsonFilePersister) *IdentifiableFilePersistence {
	var c = &IdentifiableFilePersistence{}
	if &persister == nil {
		persister = *NewJsonFilePersister("")
	}
	c.IdentifiableMemoryPersistence = *NewIdentifiableMemoryPersistence(prototype, &persister, &persister)
	c._persister = persister
	return c
}

// Configures component by passing configuration parameters.
// - config    configuration parameters to be set.
func (c *IdentifiableFilePersistence) Configure(config config.ConfigParams) {
	c.Configure(config)
	c._persister.Configure(config)
}