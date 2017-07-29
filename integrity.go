package sri

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileSet struct {
	// A mapping of path name to SHA512 hash
	hashes map[string]string
	root   string
}

// We're not actually parsing anything
func ParseFiles(root string) (*FileSet, error) {
	if root == "" {
		root = "."
	}
	hashes := map[string]string{}
	return &FileSet{hashes, root}, filepath.Walk(root, func(p string, info os.FileInfo, other error) error {
		if other != nil {
			return other
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		hash := sha512.New()
		hash.Write([]byte(data))
		hashes[p] = "sha512-" + base64.StdEncoding.EncodeToString(hash.Sum(nil))
		return nil
	})
}

func (b *FileSet) Digest(p string) (string, error) {
	fp := filepath.Join(b.root, filepath.FromSlash(path.Clean("/"+p)))
	if h, ok := b.hashes[fp]; ok {
		return h, nil
	}
	return "", fmt.Errorf("no integrity value for subresource %s", p)
}

type Verifier struct {
	// Query parameter that contains the SHA512 integrity value. Defaults to
	// "sha".
	Param string

	// If an integrity value is provided and it doesn't match the file, return
	// a 404 using this handler. Defaults to http.NotFoundHandler
	NotFound http.Handler

	bundle *FileSet
	next   http.Handler
}

func (v *Verifier) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if sha := values.Get(v.Param); sha != "" {
		if integ, _ := v.bundle.Digest(r.URL.Path); integ != sha {
			v.NotFound.ServeHTTP(w, r)
			return
		}
	}
	v.next.ServeHTTP(w, r)
}

func Verify(b *FileSet, next http.Handler) *Verifier {
	return &Verifier{
		next:     next,
		bundle:   b,
		NotFound: http.NotFoundHandler(),
		Param:    "sha",
	}
}
