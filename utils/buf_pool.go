package utils

type bufferPool struct {
	freeList chan []byte
	size     int
}

func (pool *bufferPool) Get() (buffer []byte) {
	select {
	case buffer = <-pool.freeList:
	default:
		buffer = make([]byte, pool.size)
	}
	return
}

func (pool *bufferPool) Put(buffer []byte) {
	select {
	case pool.freeList <- buffer:
	default:
	}
}

var Pool32K = &bufferPool{freeList: make(chan []byte, 100), size: 32 * 1024}
var Pool33K = &bufferPool{freeList: make(chan []byte, 100), size: 33 * 1024}
