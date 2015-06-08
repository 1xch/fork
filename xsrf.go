package fork

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type xsrf struct {
	Secret  string
	Key     string
	current string
	*baseField
	*processor
}

func (x *xsrf) New() Field {
	var newfield xsrf = *x
	newfield.baseField = x.baseField.Copy()
	newfield.current = ""
	return &newfield
}

func (x *xsrf) Get() *Value {
	return NewValue(x.Token())
}

func (x *xsrf) Set(r *http.Request) {
	v := x.Filter(x.Name(), r)
	x.current = v.String()
	err := x.Validate(x)
	if err == nil {
		x.current = ""
	}
}

func (x *xsrf) Token() string {
	if x.current == "" {
		return generateTokenAtTime(x.Secret, x.Key, time.Now())
	}
	return x.current
}

func clean(s string) string {
	return strings.Replace(s, ":", "_", -1)
}

func generateTokenAtTime(secret string, key string, now time.Time) string {
	h := hmac.New(sha1.New, []byte(key))
	fmt.Fprintf(h, "%s:%s:%d", clean(secret), clean(key), now.UnixNano())
	tok := fmt.Sprintf("%s:%d", h.Sum(nil), now.UnixNano())
	return base64.URLEncoding.EncodeToString([]byte(tok))
}

func ValidateXsrf(x *xsrf) error {
	if x.current != "" {
		valid := validTokenAtTime(x.current, x.Secret, x.Key, time.Now())
		if !valid {
			return fmt.Errorf("Invalid XSRF Token")
		}
	}
	return nil
}

const Timeout = 12 * time.Hour

func validTokenAtTime(token string, secret string, key string, now time.Time) bool {
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	sep := bytes.LastIndex(data, []byte{':'})
	if sep < 0 {
		return false
	}
	nanos, err := strconv.ParseInt(string(data[sep+1:]), 10, 64)
	if err != nil {
		return false
	}
	issueTime := time.Unix(0, nanos)

	if now.Sub(issueTime) >= Timeout {
		return false
	}

	if issueTime.After(now.Add(1 * time.Minute)) {
		return false
	}

	expected := generateTokenAtTime(secret, key, issueTime)

	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}

const xsrfWidget = `<input type="hidden" name="{{ .Name }}" value="{{ .Token }}" >`

var keybase = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randomSequence(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = keybase[rand.Intn(len(keybase))]
	}
	return string(b)
}

func XSRF(name string, secret string) Field {
	ret := &xsrf{
		Secret:    secret,
		Key:       randomSequence(12),
		baseField: newBaseField(name),
		processor: NewProcessor(
			NewWidget(xsrfWidget),
			NewValidater(ValidateXsrf),
			NewFilterer(),
		),
	}
	ret.validateable = true
	return ret
}
