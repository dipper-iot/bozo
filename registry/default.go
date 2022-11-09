package registry

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	sendEventTime = 10 * time.Millisecond
	ttlPruneTime  = time.Second
)

type defaultRegistry struct {
	options Options

	sync.RWMutex
	records  map[string]map[string]*record
	watchers map[string]*memWatcher
}

func NewDefaultRegistry(opts ...Option) Registry {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	records := getServiceRecords(options.Context)
	if records == nil {
		records = make(map[string]map[string]*record)
	}

	reg := &defaultRegistry{
		options:  options,
		records:  records,
		watchers: make(map[string]*memWatcher),
	}

	return reg
}

func (m *defaultRegistry) Init(opts ...Option) error {
	for _, o := range opts {
		o(&m.options)
	}

	// add services
	m.Lock()
	defer m.Unlock()

	records := getServiceRecords(m.options.Context)
	for name, record := range records {
		// add a whole new service including all of its versions
		if _, ok := m.records[name]; !ok {
			m.records[name] = record
			continue
		}
		// add the versions of the service we dont track yet
		for version, r := range record {
			if _, ok := m.records[name][version]; !ok {
				m.records[name][version] = r
				continue
			}
		}
	}

	return nil
}

func (m *defaultRegistry) Options() Options {
	return m.options
}

func (m *defaultRegistry) Register(s *Service, opts ...RegisterOption) error {

	return nil
}

func (m *defaultRegistry) Deregister(s *Service, opts ...DeregisterOption) error {

	return nil
}

func (m *defaultRegistry) GetService(name string, opts ...GetOption) ([]*Service, error) {
	m.RLock()
	defer m.RUnlock()

	o := NewGetOption(opts...)

	records, ok := m.records[name]
	if !ok {
		return nil, ErrNotFound
	}

	services := make([]*Service, len(m.records[name]))
	i := 0
	for _, record := range records {
		// check version
		if len(o.Version) > 0 && record.Version != o.Version {
			continue
		}

		services[i] = recordToService(record)
		i++
	}

	return services, nil
}

func (m *defaultRegistry) ListServices(opts ...ListOption) ([]*Service, error) {
	m.RLock()
	defer m.RUnlock()

	var services []*Service
	for _, records := range m.records {
		for _, record := range records {
			services = append(services, recordToService(record))
		}
	}

	return services, nil
}

func (m *defaultRegistry) Watch(opts ...WatchOption) (Watcher, error) {
	var wo WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	w := &memWatcher{
		exit: make(chan bool),
		res:  make(chan *Result),
		id:   uuid.New().String(),
		wo:   wo,
	}

	m.Lock()
	m.watchers[w.id] = w
	m.Unlock()

	return w, nil
}

func (m *defaultRegistry) String() string {
	return "default"
}
