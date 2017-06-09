package utils

type BufferPool struct {
	freeList chan []byte
	size     int
}

func (pool *BufferPool) Get() (buffer []byte) {
	select {
	case buffer = <-pool.freeList:
	default:
		buffer = make([]byte, pool.size)
	}
	return
}

func (pool *BufferPool) Put(buffer []byte) {
	select {
	case pool.freeList <- buffer:
	default:
	}
}

var Pool32K = &BufferPool{freeList: make(chan []byte, 100), size: 32 * 1024}
var Pool33K = &BufferPool{freeList: make(chan []byte, 100), size: 33 * 1024}
