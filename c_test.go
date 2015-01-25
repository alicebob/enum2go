package main

import (
	"reflect"
	"strings"
	"testing"
)

var (
	someC = `# 28 "/usr/include/stdint.h" 2 3 4
# 36 "/usr/include/stdint.h" 3 4
typedef signed char int8_t;
typedef short int int16_t;

enum
  {
          SCM_RIGHTS = 0x01





};
`
	someCwithStr = `static struct pnp_card_device_id sb_pnp_card_table[] = {
{.id = "CTL0024", .driver_data = 0, .devs = { {.id="CTL0031"}, } },`
)

func TestTokenize(t *testing.T) {
	for code, keywords := range map[string]string{
		someC:        `typedef signed char int8_t ; typedef short int int16_t ; enum { SCM_RIGHTS = 0x01 } ;`,
		someCwithStr: `static struct pnp_card_device_id sb_pnp_card_table[] = { { . id = "CTL0024" , . driver_data = 0 , . devs = { { . id = "CTL0031" } , } } ,`,
	} {
		have, err := cTokenize(strings.NewReader(code))
		if err != nil {
			t.Error(err)
		}
		want := strings.Fields(keywords)
		if !reflect.DeepEqual(have, want) {
			t.Errorf("have: %#v, want: %#v", have, want)
		}
	}
}
