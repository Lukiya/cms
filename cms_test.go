package cms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/go/sconfig"
)

func Test_defaultCMS_Render(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	a := NewJetCMS(cp)

	params := GetParams()
	params.Set("s", []string{"zzz", "yyyy"})

	b, err := a.Render("/", params)
	assert.NoError(t, err)
	t.Log(b)
}

func Test_defaultCMS_GetContent(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	a := NewJetCMS(cp)

	b := a.GetContent("/", nil, false, true)
	t.Log(b)
}
