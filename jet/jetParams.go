package jet

import (
	"sync"

	"github.com/CloudyKit/jet/v6"
)

var (
	_paramPool = sync.Pool{
		New: func() interface{} {
			v := make(jet.VarMap)
			return &JetParams{
				data: &v,
			}
		},
	}
)

type JetParams struct {
	data *jet.VarMap
}

func (o *JetParams) Reset() {
	*o.data = make(jet.VarMap)
}

func (o *JetParams) Set(name string, v interface{}) {
	o.data.Set(name, v)
}

func GetParams() *JetParams {
	return _paramPool.Get().(*JetParams)
}

func ReleaseParams(params *JetParams) {
	params.Reset()
	_paramPool.Put(params)
}
