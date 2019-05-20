package chanMonitor

import (
	"fmt"
	"reflect"
	"sync"
)

var chans = make(map[key]interface{})
var chmu sync.RWMutex

// AddNamed adds a channel to be monitor and associates the channel
// with this name and, optionally, the instance of this named channel (there may be many)
func AddNamed(name, instance string, channel interface{}) error {

	//reflect on the input to get the correct channel type.
	if reflect.TypeOf(channel).Kind() != reflect.Chan {
		return fmt.Errorf("invalid input type %v for input param channel, must be of type chan", channel)
	}

	chmu.Lock()
	defer chmu.Unlock()

	k := key{name: name, instance: instance}

	if _, found := chans[k]; found {
		return fmt.Errorf("channel with name: %s already being monitored.", name)
	}
	chans[k] = channel

	return nil
}

// ChanState struct holding Length and Capacity.
type ChanState struct {
	Len      int    `json:"length"`
	Cap      int    `json:"capacity"`
	Instance string `json:"instance"`
}

type key struct {
	name     string
	instance string
}

// Get returns the channel state for a give channel name.
func Get(name, instance string) *ChanState {

	chmu.RLock()
	defer chmu.RUnlock()

	k := key{name: name, instance: instance}

	ch, found := chans[k]
	if !found {
		return nil
	}

	return &ChanState{
		Len:      reflect.ValueOf(ch).Len(),
		Cap:      reflect.ValueOf(ch).Cap(),
		Instance: k.instance,
	}

}

// Get the channel states map[string]*ChanState of all the monitored channels. Keyed by channel name.
func GetAll() map[string]*ChanState {

	results := make(map[string]*ChanState)

	chmu.RLock()
	defer chmu.RUnlock()
	for k, _ := range chans {
		results[k.name] = Get(k.name, k.instance)
	}

	return results

}
