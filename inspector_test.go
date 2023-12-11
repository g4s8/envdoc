package main

import (
	"fmt"
	"testing"
)

func TestTagParsers(t *testing.T) {
	type testCase struct {
		tag    string
		expect docItem
		fail   bool
	}
	for i, c := range []testCase{
		{tag: "", fail: true},
		{tag: " ", fail: true},
		{tag: `env:"FOO"`, expect: docItem{envName: "FOO"}},
		{tag: ` env:FOO `, fail: true},
		{tag: `json:"bar"   env:"FOO"   qwe:"baz"`, expect: docItem{envName: "FOO"}},
		{tag: `env:"SECRET,file"`, expect: docItem{envName: "SECRET", flags: docItemFlagFromFile}},
		{
			tag:    `env:"PASSWORD,file"           envDefault:"/tmp/password"   json:"password"`,
			expect: docItem{envName: "PASSWORD", flags: docItemFlagFromFile, envDefault: "/tmp/password"},
		},
		{
			tag:    `env:"CERTIFICATE,file,expand" envDefault:"${CERTIFICATE_FILE}"`,
			expect: docItem{envName: "CERTIFICATE", flags: docItemFlagFromFile | docItemFlagExpand, envDefault: "${CERTIFICATE_FILE}"},
		},
		{
			tag:    `env:"SECRET_KEY,required" json:"secret_key"`,
			expect: docItem{envName: "SECRET_KEY", flags: docItemFlagRequired},
		},
		{
			tag:    `json:"secret_val" env:"SECRET_VAL,notEmpty"`,
			expect: docItem{envName: "SECRET_VAL", flags: docItemFlagNonEmpty | docItemFlagRequired},
		},
		{
			tag: `fooo:"1" env:"JUST_A_MESS,required,notEmpty,file,expand" json:"just_a_mess" envDefault:"${JUST_A_MESS_FILE}" bar:"2"`,
			expect: docItem{
				envName:    "JUST_A_MESS",
				flags:      docItemFlagRequired | docItemFlagNonEmpty | docItemFlagFromFile | docItemFlagExpand,
				envDefault: "${JUST_A_MESS_FILE}",
			},
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out docItem
			ok := parseTag(c.tag, &out)
			if ok != !c.fail {
				t.Error("parseTag returned false")
			}
			if out != c.expect {
				t.Errorf("parseTag of `%s` returned wrong result: %+v; expected: %+v", c.tag, out, c.expect)
			}
		})
	}
}
