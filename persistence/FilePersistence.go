package persistence

import (
	"reflect"

	"github.com/pip-services3-go/pip-services3-commons-go/config"
)

/*
Abstract persistence component that stores data in flat files
and caches them in memory.

FilePersistence is the most basic persistence component that is only
able to store data items of any type. Specific CRUD operations
over the data items must be implemented in child structs by
accessing fp._items property and calling Save method.

see MemoryPersistence
see JsonFilePersister

Configuration parameters

- path - path to the file where data is stored

References

- *:logger:*:*:1.0  (optional) ILogger components to pass log messages

Example
type MyJsonFilePersistence struct {
	FilePersistence
}
    func NewMyJsonFilePersistence(path string) mjfp* NewMyJsonFilePersistence {
		return NewFilePersistence(NewJsonPersister(path))
    }

	func (fp * FilePersistence) GetByName(correlationId string, name string) (item interface{}, err error){
		for _,v := range fp._items {
			if v.name == name {
				item = v
				break
			}
		}
        return item, nil
    }

    func (fp *FilePersistence) Set(correlatonId string, item MyData) error {
		for i,v:=range fp._items {
			if v.name == item.name {
				fp._items = append(fp._items[:i], fp._items[i+1:])
			}
		}
		fp._items = append(fp._items, item)
        retrun fp.save(correlationId)
    }
}
*/
//extends MemoryPersistence implements IConfigurable
type FilePersistence struct {
	MemoryPersistence
	_persister JsonFilePersister
}

// Creates a new instance of the persistence.
// - persister    (optional) a persister component that loads and saves data from/to flat file.
// Return *FilePersistence
// Pointer on new FilePersistence instance
func NewFilePersistence(prototype reflect.Type, persister JsonFilePersister) *FilePersistence {
	var c = &FilePersistence{}
	if &persister == nil {
		persister = *NewJsonFilePersister("")
	}
	c._persister = persister
	c.MemoryPersistence = *NewMemoryPersistence(prototype, &persister, &persister)
	return c
}

// Configures component by passing configuration parameters.
// - config    configuration parameters to be set.
func (c *FilePersistence) Configure(conf config.ConfigParams) {
	c._persister.Configure(conf)
}