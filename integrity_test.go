package integrity

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"
)

const (
	DigestCSS = "sha512-tXj6GsgkfpVKb1u4ktHSdZAn3B3CEqvcv4OgB2zlCQ2MFMYa/k8qc6JaS0BTe5kTpvs8ywuhGMsqyzzNEowyFA=="
	DigestJS  = "sha512-/9UYockPRzjO73c05tr9/jrzEbKGVSZzyk+rrt1e7siZS6IXlbjZYCegSetVVQ+fIdmIIfHbpgXVkJHuFxE8Jg=="
	DigestJPG = "sha512-AJMABkkop36C2qXMAsP6xFWjW2r7Q7M8bd6/1U1m+rg3ANlrcRewBZzfS15J3NtYvHBdnrWuYYCDyjaaUduzxQ=="
)

func TestDigests(t *testing.T) {
	fs, err := ParseFiles("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for path, digest := range map[string]string{
		"css/style.css":        DigestCSS,
		"js/godocs.js":         DigestJS,
		"img/fancygopher.jpg":  DigestJPG,
		"/css/style.css":       DigestCSS,
		"/js/godocs.js":        DigestJS,
		"/img/fancygopher.jpg": DigestJPG,
	} {
		t.Run(path, func(t *testing.T) {
			d, err := fs.Digest(path)
			if err != nil {
				t.Error(err)
			}
			if d != digest {
				t.Errorf("Mismatched digest: %s != %s", digest, d)
			}
		})
	}
}

func TestVerification(t *testing.T) {
	fs, err := ParseFiles("testdata")
	if err != nil {
		t.Fatal(err)
	}

	h := Verify(fs, http.FileServer(http.Dir("testdata")))
	h.Param = "h"
	handler := http.StripPrefix("/static/", h)

	js := url.Values{"h": []string{DigestJS}}
	css := url.Values{"h": []string{DigestCSS}}

	for _, pair := range []struct {
		req  *http.Request
		code int
	}{
		{httptest.NewRequest("GET", "/static/css/style.css", nil), 200},
		{httptest.NewRequest("GET", "/static/css/style.css?h=", nil), 200},
		{httptest.NewRequest("GET", "/static/css/style.css?"+css.Encode(), nil), 200},
		{httptest.NewRequest("GET", "/static/css/style.css?h=invalid", nil), 404},
		{httptest.NewRequest("GET", "/static/css/style.css?"+js.Encode(), nil), 404},
	} {
		t.Run(filepath.Base(pair.req.URL.String()), func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, pair.req)
			resp := w.Result()
			if resp.StatusCode != pair.code {
				t.Errorf("Expected status code of %d, not %d", pair.code, resp.StatusCode)
			}
		})
	}
}
