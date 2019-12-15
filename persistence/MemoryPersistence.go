package persistence

import (
	"encoding/json"
	"github.com/pip-services3-go/pip-services3-commons-go/convert"
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-components-go/log"
	"reflect"
	"sync"
)

/*
Abstract persistence component that stores data in memory.

This is the most basic persistence component that is only
able to store data items of any type. Specific CRUD operations
over the data items must be implemented in child struct by
accessing _items property and calling Save method.

The component supports loading and saving items from another data source.
That allows to use it as a base struct for file and other types
of persistence components that cache all data in memory.

References

- *:logger:*:*:1.0       (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-go/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages

Example

    type MyMemoryPersistence struct {
        MemoryPersistence

    }
     func (c * MyMemoryPersistence) GetByName(correlationId string, name string)(item interface{}, err error) {
        for _, v := range c._items {
            if v.name == name {
                item = v
                break
            }
        }
        return item, nil
    });

    func (c * MyMemoryPersistence) Set(correlatonId: string, item: MyData, callback: (err) => void): void {
        c._items = append(c._items, item);
        c.Save(correlationId);
    }

    persistence := NewMyMemoryPersistence();
    err := persistence.Set("123", interface{}({ name: "ABC" }))
    item, err := persistence.GetByName("123", "ABC")
    fmt.Println(item)   // Result: { name: "ABC" }
*/
// implements IReferenceable, IOpenable, ICleanable
type MemoryPersistence struct {
	_logger    log.CompositeLogger
	_items     []interface{}
	_loader    ILoader
	_saver     ISaver
	_opened    bool
	_prototype reflect.Type
	_lockMutex sync.RWMutex
}

// Creates a new empty instance of the MemoryPersistence
// Return *MemoryPersistence
// empty MemoryPersistence
func NewEmptyMemoryPersistence(prototype reflect.Type) (mp *MemoryPersistence) {
	if prototype == nil {
		return nil
	}
	mp = &MemoryPersistence{}
	mp._prototype = prototype
	mp._logger = *log.NewCompositeLogger()
	mp._items = make([]interface{}, 0, 10)
	return mp
}

// Creates a new instance of the persistence.
// Parameters:
//    - loader ILoader
//    (optional) a loader to load items from external datasource.
//    - saver  ISaver
//    (optional) a saver to save items to external datasource.
// Return *MemoryPersistence
// MemoryPersistence
func NewMemoryPersistence(prototype reflect.Type, loader ILoader, saver ISaver) (mp *MemoryPersistence) {
	if prototype == nil {
		return nil
	}
	mp = &MemoryPersistence{}
	mp._items = make([]interface{}, 0, 10)
	mp._loader = loader
	mp._saver = saver
	mp._logger = *log.NewCompositeLogger()
	return mp
}

//  Sets references to dependent components.
//  Parameters:
// 	- references refer.IReferences
//	references to locate the component dependencies.
func (c *MemoryPersistence) SetReferences(references refer.IReferences) {
	c._logger.SetReferences(references)
}

//  Checks if the component is opened.
//  Returns true if the component has been opened and false otherwise.
func (c *MemoryPersistence) IsOpen() bool {
	return c._opened
}

// Opens the component.
// Parameters:
// 		- correlationId  string
// 		(optional) transaction id to trace execution through call chain.
// Returns  error or null no errors occured.
func (c *MemoryPersistence) Open(correlationId string) error {
	c._lockMutex.Lock()
	defer c._lockMutex.Unlock()
	err := c.load(correlationId)
	if err == nil {
		c._opened = true
	}
	return err
}

func (c *MemoryPersistence) load(correlationId string) error {
	if c._loader == nil {
		return nil
	}

	items, err := c._loader.Load(correlationId)
	if err == nil && items != nil {
		c._items = make([]interface{}, len(items))
		for i, v := range items {
			item := convert.MapConverter.ToNullableMap(v)
			jsonMarshalStr, errJson := json.Marshal(item)
			if errJson != nil {
				panic("MemoryPersistence.Load Error can't convert from Json to any type")
			}
			value := reflect.New(c._prototype).Interface()
			json.Unmarshal(jsonMarshalStr, value)
			c._items[i] = reflect.ValueOf(value).Elem().Interface()
		}
		length := len(c._items)
		c._logger.Trace(correlationId, "Loaded %d items", length)
	}
	return err
}

// Closes component and frees used resources.
// Parameters:
// 	- correlationId 	(optional) transaction id to trace execution through call chain.
// Retruns: error or null no errors occured.
func (c *MemoryPersistence) Close(correlationId string) error {
	err := c.Save(correlationId)
	c._opened = false
	return err
}

// Saves items to external data source using configured saver component.
//    - correlationId string
//     (optional) transaction id to trace execution through call chain.
// Return error or null for success.
func (c *MemoryPersistence) Save(correlationId string) error {
	c._lockMutex.RLock()
	defer c._lockMutex.RUnlock()

	if c._saver == nil {
		return nil
	}

	err := c._saver.Save(correlationId, c._items)
	if err == nil {
		length := len(c._items)
		c._logger.Trace(correlationId, "Saved %d items", length)
	}
	return err
}

// Clears component state.
// 	- correlationId 	(optional) transaction id to trace execution through call chain.
//  Returns error or null no errors occured.
func (c *MemoryPersistence) Clear(correlationId string) error {
	c._lockMutex.Lock()
	defer c._lockMutex.Unlock()

	c._items = make([]interface{}, 0, 5)
	c._logger.Trace(correlationId, "Cleared items")
	return c.Save(correlationId)
}
