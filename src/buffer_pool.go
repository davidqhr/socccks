package socks5

var BUFFER_SIZE = 1024
var freeList = make(chan []byte, 100)

type BufferPool struct {
	freeList chan []byte
}

func (pool *BufferPool) Get() []byte {
	var buffer []byte
	select {
	case buffer = <-pool.freeList:
	default:
		buffer = make([]byte, BUFFER_SIZE)
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
