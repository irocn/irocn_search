package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/static/_index.html": {
		local:   "../static/_index.html",
		size:    22049,
		modtime: 1519111881,
		compressed: `
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
84Ow/xsAAP//9y2EyiFWAAA=
`,
	},
