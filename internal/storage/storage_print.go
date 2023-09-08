package storage

import (
	"bytes"
	"fmt"
)

// Prints metricsstorage
func (ms *MemStorage) String() string {
	buf := &bytes.Buffer{}
	for _, v := range ms.storage {
		fmt.Fprintf(buf, " %f:  %v:", v.Gauge, v.Count)
	}
	return buf.String()
}
