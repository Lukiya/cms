package cms

import (
	"sync"

	"github.com/CloudyKit/jet/v6"
	"github.com/syncfuture/go/spool"
)

type ICMS interface {
	GetHtml(key string, args ...interface{}) string
	Render(key string, args ...interface{}) (string, error)
}

const _err1 = "could not be found"

var (
	_paramPool = sync.Pool{
		New: func() interface{} {
			return make(jet.VarMap)
		},
	}
	_bufferPool = spool.NewSyncBufferPool(1024)
)

func MakeParams() jet.VarMap {
	return _paramPool.Get().(jet.VarMap)
}

func releaseParams(params jet.VarMap) {
	for k := range params {
		delete(params, k)
	}
	_paramPool.Put(params)
}
