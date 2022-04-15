package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCookie(t *testing.T) {
	r := ParseCookie("_MHYUUID=adfdas-adfa-adfad-afd-adfasf; ltoken=adsfadsfadsfadsf; ltuid=213123; cookie_token=adfadf; account_id=2121; login_uid=1212; login_ticket=dfadfasf")
	if v, ok := r["login_ticket"]; !ok {
		t.Fatal(v)
	}
	assert.Equal(t, "dfadfasf", r["login_ticket"])
}
