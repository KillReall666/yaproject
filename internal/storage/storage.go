package storage

type Metrics struct {
	Count int64
	Gauge float64
}
type MemStorage struct {
	storage map[string]*Metrics
	//mtx     *sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]*Metrics),
	}
}

func (ms *MemStorage) Count(name string, count int64) {
	//ms.mtx.Lock()
	_, ok := ms.storage[name]
	if !ok {
		ms.storage[name] = &Metrics{}
		ms.storage[name].Count += count
		return
	}
	ms.storage[name].Count = count
	//ms.mtx.Unlock()
}

func (ms *MemStorage) Gauge(name string, gauge float64) {
	//ms.mtx.Lock()
	_, ok := ms.storage[name]
	if !ok {
		ms.storage[name] = &Metrics{}
		ms.storage[name].Gauge = gauge
		return
	}
	ms.storage[name].Gauge = gauge
	//ms.mtx.Unlock()
}
