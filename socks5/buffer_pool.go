package socks5

var bufferSize = 1024
var freeList = make(chan []byte, 100)

type BufferPool struct {
	freeList chan []byte
}

func (pool *BufferPool) Get() []byte {
	var buffer []byte
	select {
	case buffer = <-pool.freeList:
	default:
		buffer = make([]byte, bufferSize)
	}
	return buffer
}

func (pool *BufferPool) Put(buffer []byte) {
	select {
	case pool.freeList <- buffer:
	default:
	}
}

var bufferPool = &BufferPool{freeList: make(chan []byte, 100)}
