package fork

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func extracttoken(from *httptest.ResponseRecorder, by string) string {
	b := bytes.Fields(from.Body.Bytes())

	for _, x := range b {
		if bytes.Contains(x, []byte(by)) {
			return strings.Split(string(x), `"`)[1]
		}
	}

	return ""
}

func TestXSRF(t *testing.T) {
	field := XSRF("testXSRF", "SECRET")
	xsrffield, _ := field.(*xsrf)

	f := NewForm("XSRF", Fields(field))

	ts := testServe()
	ts.handlers["GET"] = formHandlerGet(t, f)
	ts.handlers["POST"] = formHandlerPost(t, f)

	w1 := PerformGet(ts)

	expect := `<input type="hidden" name="testXSRF" value="`

	if !strings.Contains(w1.Body.String(), expect) {
		t.Errorf("\n%s XSRF GET Error\ngot\n\n%s\nshould contain:\n\n%s\n\n", w1.Body, expect)
	}

	token := extracttoken(w1, "value=")

	if !validTokenAtTime(token, xsrffield.Secret, xsrffield.Key, time.Now()) {
		t.Errorf("\nInvalid xsrf token: %s\n", token)
	}

	sendtoken := fmt.Sprintf(`testXSRF=%s`, token)

	w2 := PerformPost(ts, sendtoken)

	posttoken := extracttoken(w2, "value=")

	if posttoken == token {
		t.Errorf("\n %s == %s ; Indicative of an invalid token\n", token, posttoken)
	}

	w3 := PerformPost(ts, `testXSRF=invalidtoken`)

	invalidtoken := extracttoken(w3, "value=")

	if invalidtoken != "invalidtoken" {
		t.Errorf(`\nInvalid token should be "invalidtoken", not %s\n`, invalidtoken)
	}

	var invalidresult string = `<input type="hidden" name="testXSRF" value="invalidtoken" ><div class="field-errors"><ul><li>Invalid XSRF Token</li></ul></div>`

	if !strings.Contains(w3.Body.String(), invalidresult) {
		t.Errorf("\n%s POST Error\ngot %s\nshould contain %s\n\n", w3, invalidresult)
	}
}
