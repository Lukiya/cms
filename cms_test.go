package cms

import (
	"testing"

	"github.com/syncfuture/go/sconfig"
)

func Test_defaultCMS_GetHtml(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	a := NewJetHtmlCMS(cp)

	params := MakeParams()
	params.Set("s", []string{"zzz", "yyyy"})

	b := a.GetHtml("/", params)
	t.Log(b)
}
