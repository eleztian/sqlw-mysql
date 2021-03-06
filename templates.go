// Code generated by "esc -o templates.go templates"; DO NOT EDIT.

package main

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
		_ = f.Close()
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

	"/templates/default/interface.go.tmpl": {
		local:   "templates/default/interface.go.tmpl",
		size:    2238,
		modtime: 1535858993,
		compressed: `
H4sIAAAAAAAC/6SVTW/cNhCG7/oVA18iuRup5wI9pHYOBlzHddz0EBQFVxpJbCiOTI6iXSz2vxf80Jdj
owhy2yXfeeadGZIqCng3MEGDGo1grGB/BPukxrfd0T4pSFvm3v5SFI3kdtjnJXVFOwjd/DvoEXWxSDOo
sBaDYmDseiUY86Qo4PrDHdx9eIT31zePeZL0ovwiGoTTCfL78PtOdAjnc5LIrifDkCYAFyVpxgNfuN+V
YLEXFl2yiyRLHPb9AUs0IC1wi1BS15EGqRlNLUoEJsADlgMjCHga0BxhlNzSwGCQB6OlbkDoIxga84SP
Pc7EmXFKwC9eBStpyQeItvK4totsy0bqZgfCNBbyPJ8hp3MGqX1S+QPaQfEO0BgyWXL2Rfzhgr+visW9
oTG1WTQ/kzbu/eoP27/0/mm0s/tIfqDxR+ETOzbEtXsq5ddpIj9N1S1zerncoJ/coYnQK1JDpz+WQms0
YLA3aFGzBQElKYUlS9JAtfs3dNrNQjCUQsMewfqwymUiqA11IFZnZoveuikK0Id7g24XRN+jriwIpZY0
BBVahppMSCN1kyewRKV++/Lz3+uWJTOaLHt2N1h2VkuhFFYgasYF6LMQWOoQ8MBGQG+oRGvtki2C0izM
N7btUewVPtC47Ri71Td21YNZty1/05po2kt/RxbxGIdjPyE7tyF1TaYTbiTO3RyRZnA5/4m4T0LJakGZ
AUHWwK207kYJ+Or3Pd75hRTzJgdNDHd/3t5mju8RaQZ7IhWpwfgnoQbc2JRvuY2zc+QB87lKL06l60AG
q2FtiPfkd15jvrHQB8XCjSEvkrdD+ktye29kJ/wTMc/LVU11bIF7AKGPoi94fDa/NWI7ykkRy4mibzu0
YocO+dxuHP60hCPQTjFYRZG07uXrhfEu9zHWuoB1qnTbgFD/RyXLZ5d6umyeZt1+LDRqX7ilt6g3hSjU
DbfB/cyAIAw25lt4w9h9O1LpVjfhIMG2NKjK3VQ3Canh88+7iZntgLhFM0rr3vpeaFnCKJVyciOkxSo4
cPleP2n68M7Xvrw3OnhhWnu5qcOqtKCl2rktPd8X7Ho+hv3FceBNJkKSNDDWb9M5Sb4K47/g/zx7zMNn
5Pq3LNVSZa8Krkjr/5E8HqIgS/4LAAD//7nL3CG+CAAA
`,
	},

	"/templates/default/manifest.json": {
		local:   "templates/default/manifest.json",
		size:    271,
		modtime: 1535858993,
		compressed: `
H4sIAAAAAAAC/2yNQcqDMBCF9znFkPVPDuAZfrtoXRRKCVOZFouJg3kuRLx7UdsKJZuBee/7eJMhsqnm
WI0qJastttdjVPGB1T1TF+3fgqn0xyHagi6GiMg2EdLfuRb36ByCtitGZIOAf7MBTZvjPCThWxii62er
4lsr+xqW10+TW/PtHjjIPGftEwLO5f/uJwQs+rvIqGY2rwAAAP//zYpLFg8BAAA=
`,
	},

	"/templates/default/meta.go.tmpl": {
		local:   "templates/default/meta.go.tmpl",
		size:    10577,
		modtime: 1535858993,
		compressed: `
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
zB2pLTiPvX8FAAD//3PwqTxRKQAA
`,
	},

	"/templates/default/meta_test.go.tmpl": {
		local:   "templates/default/meta_test.go.tmpl",
		size:    9170,
		modtime: 1535858993,
		compressed: `
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
IwAA
`,
	},

	"/templates/default/scan_type_map.json": {
		local:   "templates/default/scan_type_map.json",
		size:    664,
		modtime: 1535858993,
		compressed: `
H4sIAAAAAAAC/4TRP8+CMBDH8Z1X0TATEv6kecL4DCbO6tQ4gKI5U9pEymR876a0dz0XZPvBJ98beGVC
5IO1Ou/E+qiwCpGbRevy349z4RUY95fUulDt/QhqYUzFhewEyYFxlSQXFstVkvUQKlw8SBKMa2pebGpe
bGpWRKhw8SJJME62vChbXpQtKyJUuHiR5E3bHo8rWkh3cTMbs4rWl6Wug2lM/8av8uhfoV5HoNfxAlOv
Y3Z2TzB3cocwgxzAYfOHfMzWpPNbMn7rNmX2zj4BAAD//+zO34KYAgAA
`,
	},

	"/templates/default/stmt_{{.StmtXMLName}}.go.tmpl": {
		local:   "templates/default/stmt_{{.StmtXMLName}}.go.tmpl",
		size:    9667,
		modtime: 1535858993,
		compressed: `
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
/w8zo1fDJQAA
`,
	},

	"/templates/default/table_{{.Table.TableName}}.go.tmpl": {
		local:   "templates/default/table_{{.Table.TableName}}.go.tmpl",
		size:    9516,
		modtime: 1536726154,
		compressed: `
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
`,
	},

	"/templates/default/util.go.tmpl": {
		local:   "templates/default/util.go.tmpl",
		size:    953,
		modtime: 1535858993,
		compressed: `
H4sIAAAAAAAC/3RSTU/rSBC8z68ockBjATZ7hc1KrOCAtAK0IN4BcRjsdjIv88W4/ZIo+L8/jT8ScngX
yzNdXVXdNUWBm5Y9FuQoKqYKH1s0n2Z9YbfNp4FcMofmqigWmpftR156Wyxb5RY/W7cmVxygGSqqVWsY
TDYYxZSLosDt4wMeHl9wd3v/kgsRVLlSC8Juh/xp+H9QltB1QmgbfGRIAcxqyzORCVG3rsQi+jb8u5VN
LPFsdEnxHA0ZKtlHJITUjinWqqRdl+FFfRj6369/aF4+RW1V3J6DY7Pv7fmmY4adEICuEyR3m//IyQwn
c1zi6wtyxJ7M4bTB6enYfARMFAAQlNOlrC3ndzH6WMvZ3SZQySAbeIsm6TWzLBNAlzRD1LbB1RxWrUha
Fd6+zfGuHWcoCowjvCrTEi7+gXYVbeB7v8XgRgC1j9CJ6/IaGn+jieXk8Rr67GyYEtBMNsGG+j2TlTrr
Kxz7+3GvMgGHQv8pCjyvdED06wZrzUu41hj8Sq7ygbnGCcf8VRldyWklQOkda9dSf+wGEyEpccy/jyYH
MV1tzuFXCdCv5y2879n9as+a3i0crY+3U/vWVfkIGeK8CYFcJTlm47WucZzpRIlDtGOT02bq6qbuaoP5
0Uu5wF9jbfKLeYJ9n/ePmnvFIYhqk+VyfJYHG4ckerpOiE6I3wEAAP//WFCx1rkDAAA=
`,
	},

	"/templates/graphviz/db.dot.tmpl": {
		local:   "templates/graphviz/db.dot.tmpl",
		size:    1020,
		modtime: 1535858993,
		compressed: `
H4sIAAAAAAAC/3xSPW/bMBDd+SseCI+10wKdbJEomqBLiwxBOhRFBso8K0QZSqAYNAXB/15QlC05/th0
fB+H93TaNF51z9A1ImNAmX4zAPDK/dHGC/7jgTPgacMA12oa4f5ZdSR4Z5Vxgd7CyGFAjFnbEFZ3X1eP
qrbUIyU2qHiMKG/36oWQEh/tAKtqsqKqQkZRt16TF/wjx5as3Y+fyth3amtck2E5yoEqeFkFjW2bcSf4
Zy5jXM73LVOqboKW1U3wJ0KZKQUOWj7+66bhp+tN40gfHu5frc2uZ9zyyjH/bWtfX9yQHphvm435QaNr
fRBDOUUztvOFQ1nTOMEt7UKO844w7D9xO6dZ3amgatVTDnZNnQVmVzT73EjpV4wgpy/uvJCA7+32jV21
Om2y8NgBzhayUJ42KEc1Hdy377OyY8TC0+62tT3WAquH4fv4hxykC/MBi21rB+aZ/zaZZYZxmt4m94WZ
U3MJ2arc3fG1rw/gcUtrwlIWZTG9LB7xd3ey/ruZIs1KmyZ2NCT2PwAA//88MZDf/AMAAA==
`,
	},

	"/templates/graphviz/manifest.json": {
		local:   "templates/graphviz/manifest.json",
		size:    40,
		modtime: 1535858993,
		compressed: `
H4sIAAAAAAAC/6rmUlBQKkgtCirNU7JSiOZSUFBQUEpJ0kvJL9EryS3IUeJSUIjlquUCBAAA//+7ac17
KAAAAA==
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/templates": {
		isDir: true,
		local: "templates",
	},

	"/templates/default": {
		isDir: true,
		local: "templates/default",
	},

	"/templates/graphviz": {
		isDir: true,
		local: "templates/graphviz",
	},
}
