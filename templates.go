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

	"/templates/default/group.tmpl": {
		local:   "templates/default/group.tmpl",
		size:    830,
		modtime: 1530326882,
		compressed: `
H4sIAAAAAAAC/3RSTWvcMBC961e83UOQycZOr023kEIOhVJCU9JDyEG15XSw9cFYZr1s/N+LLHs/DrkY
S/PmvXlv5FXZqDeNwwH5Y/r/qYzGOApBxjsOkAJY1yasRSZE3dsSb+x6/20vOy7x1FKpeYNOt7oMjhER
kmzQXKtSH8YMv9XfVv9yuz8U/j0yGcX7DQJ3x96JbzlmOAgBUB0huR1+aCszrLa4xfs75IxdbWGpxdXV
3HwBjBQA4JWlUtYm5A/Mjmu5fhi8LgO08WGPLup16ywTwBg1PZPp8HkLoxotjfIvZz5eyYYMRYHZwrNq
e42bryBb6QFumrdI0wigdgyKXLd3IHxBx+Uy4x3o+jq5BChoE2Gp/j1oIymbKoGn+zlXGYGpMH2KAk8N
edi+bcFu1+WJrsYqcP6sWqrkkgNQOhvI9no6jknZR/rA+bkfmRSoGjZwTQRMmbz41yO7a46sRYF7WL27
jKR2va3yGZJ2eO+9tpUMnM3XVONykQslTvucmyy1S9e4dFcDthfP4waf5toyL7YRdu73Q82jYkq/GrJc
zm/xNMYp/oluFGIU4n8AAAD//8prifY+AwAA
`,
	},

	"/templates/default/headnote": {
		local:   "templates/default/headnote",
		size:    103,
		modtime: 1529743910,
		compressed: `
H4sIAAAAAAAC/0TKsRHCMAwF0J4pfgcUWD01I7CAQ4QlzpEJki6X7SnTPyI8RR1v7Qx11IyBxsa/Gjxj
2uFr327L7mvHRSK+fidqGpJTeY2FJKu1T9rGRge9FuAx7BzgWaOc/gEAAP//80UxpWcAAAA=
`,
	},

	"/templates/default/intf.tmpl": {
		local:   "templates/default/intf.tmpl",
		size:    1931,
		modtime: 1530328445,
		compressed: `
H4sIAAAAAAAC/6SUTW/bPAzH7/4URC+1izzO8xWGYocCXZG1W3codlBkOhYqi65EzwmCfPdBL/FLXw5D
b7ZI/vgnKaoT8lnsEI5HKDfx+060CKdTlqm2I8uQZwAXkgzjni/8dyVYbIXDtXvRF1mRZes1fN2jRAvK
ATcIktqWDCjDaGshEZgA9yh7RhDw0qM9wKC4oZ7BIvfWKLMDYQ5gaSgzPnQ4EkfGMYNweB2l5JL3kGSV
6WyV2I6tMrsVCLtzUJblCDmeCsjdiy7v0fWaV4DWki2yUyjiuw/+tyom9ZaG3BVJ/EhaqA+nn5Z/FfTT
4Eb1iXxPw2fhZ3ZqyDXpvjUPUhiDFix2Fh0adiBAktYoWZEBqv1f3xrfNsEghYEtggthle8BQW2pBTEb
7xK97NN6DWa/seitILoOTeVAaD2lIajQMdRkYxpldmUGU1QezFdPv+fVZSOaHAd22zv2UqXQGisQNeME
DFkIHLUIuGcroLMk0Tk3ZUugvIijSG37IbYa72lYdoz96aWb9WD0W5a/aE0SHVy/IYt04+INPSNbb1Cm
JtsKPxKvbozIC7gafxLuUWhVTSjbI6gauFHOX34Bf4I94L1eyLHclWCI4e7n7W3h+QGRF7Al0okahT8K
3eNCpvqPmzQ7T+6xHKsMzrnyHShgNqwFcUPB8hHz0kEXPSZuCnmXvBzSL8XNxqpWhG0e5+Wrpjq1wL9V
0CWnZzy8mt8csRzl2SOVk5zedmjGjh0Kuf04wm2JV6A5x2CVnJTzj1QnbFC5TbHOB8xT5csGxPoftJKv
lvq8bIHmvD0Vmnzf2dJbNItCNJodN1H9yIDoGGWMW3jD2L4dqfKni3BQ4BrqdeU31U9CGXj6f3VmFisg
btAOyvlnuRNGSRiU1t7dCuWwigp8vo9vmtl/CbVP742JWpjmWm7qeKocGKVX3mTGfcG240O0T4oj7ywi
JskjY/42nbK/AQAA//+OXtwViwcAAA==
`,
	},

	"/templates/default/manifest.json": {
		local:   "templates/default/manifest.json",
		size:    258,
		modtime: 1529741733,
		compressed: `
H4sIAAAAAAAC/2yNSwoCMRBE9zlF0ethDjBXEQlxbH9MPphyITJ3l7TgJ7ireq+SfjhA6hyS572oj6HI
1IHxUnOSoQ1PGvYpU9vmnc0wlkUmtO9aC7vFRhZGk8NLVUbaCUb+COUsEzZWADknHr49IFEZOnS85lvp
GLXyz0v/4Ya3Dljd6p4BAAD//wb8EYkCAQAA
`,
	},

	"/templates/default/meta.tmpl": {
		local:   "templates/default/meta.tmpl",
		size:    10122,
		modtime: 1530328858,
		compressed: `
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
AAA=
`,
	},

	"/templates/default/meta.tmpl.bak": {
		local:   "templates/default/meta.tmpl.bak",
		size:    10793,
		modtime: 1530328853,
		compressed: `
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
`,
	},

	"/templates/default/meta_test.tmpl": {
		local:   "templates/default/meta_test.tmpl",
		size:    6246,
		modtime: 1529627530,
		compressed: `
H4sIAAAAAAAC/8xYX2/iuhJ/Tj7FbB5WySobGlqhVaUKsSW9t1KX9lK6K93N6mCCoVGdmHWc/jmI735k
OymBGCjbHqkvwYxnJjPz+3lsZ4aiOzTFMJ+Dd6XGPZRgWCxMM05mlHGwTcOKaMrxI7dMw+I443E6lcM4
wZZpGtY05rf5yIto0sg4wzy6ZQ2pN3lqoCzDTFpO6exu6sVpo9sZdD53L//TmNLP2W+S0OjOu/dXVO4p
QTwmmOPotpHmhHj3Lct0THOSpxEMcMYHaETwN8yRzeFTEZQ3cGBumoZ6JxyfgBp5Pfxgc8c0jblpGAnm
SMzZnwYZ9x07jYnjLd05plE48ILfOSK2xTPuWy4IO48LPVGimtrPXxlncTqdW5HQtmZN+ZRj5FuLwkFE
SZ6kwkOmcRGnfN50wS+1ZyxOEHs6lUbZFa3bHBaaKOf0PI2Upk7x568RpWQ+QSTDLhQ/nOXFs3zjLcq6
eIJywmseEjT7qXKUcYo8j+FApXoMvsr2GJoq4WM4XM052x3VMpwiwtJDnF2pStTRGc784UkbOr0uDGfN
4UnbWq9dOq5bXQcXwekAhpE/dKWdeMox8odw1r/8BkOB+7D0lmGCI/6/HOti6AYXwSComsGP/wb9ADYG
N8aC3F+firRKvwsdR5sv42hzf45GzR28vMCprSeiCwc195/9V5DxNQSMJAH3pttLKfYiRr2IKDvhX1Ra
3HkqTAf9PTqcMR65IBqqC5gxMV82WKHgmEY8kRMfTiCNCQiqce8McURszJiMwBjjCWYwHnmnhGbYFm2z
0YAOIfA3ZhTuEclx5hU8Fa6DxxmOePCII9s6710H/QGc9waX5UoIQ3u5zMLQge+di5vgWsjDsO1CGLbD
0LEczzQM40fMbztsmslSHZQyQvqY5yzt4ywn3K7kVEj8Qxd8R9Sby6xFZ58vluXv0YAxyuy4LGmxo3lf
UXQ3ZTRPx7bjgqjeR84cp25ZyRTxmKbZD8zEYrSdGsry5ZF/DCKHWVNSdObLH+Qfgy/6ImfFam80IHic
kTjiIFYNxGnEcIJTrgr9p3WutjN9zZePDdUX+5ALh6/DQNiK3H1XjppiexAjUYhDMXo3IDUVSL4CqYJR
dUnezMaI43eyJHsUilYHnEIuQyv4srYKUvwwkBLO6kXLy5y2ltuFj9KLs+TtZYqLALa9fxX3Zy7Uw5L/
vMiHE7DQKLJ0tL+56nYGQcn460Dt4Cdh2K5uuOLv85Ybhu0au4V7ye7mi9l9UJL71dXbh7Sy1lUCduXG
8R4I+HKYayhuPy7tQu9NYBuXdXzTHrPeLq7fDVqNBshDBoyeoDjGeH+IofSzxwn61djSh2wFWfF/123L
8TrjcZ8+2Gq+3Ms0eKrz2u4FLA+L2/YR+aKtO4kWBniI+S0Qkdw7AATOLvuguuw7wkacuHfho26NlvWW
EEV5xmkCEU3HsVjr+t31DaCIqlAgf/hBu23J8A//PQyOXotBca+yZD7VbNrytqnC3x+fowo+Kxek7P+Y
0b0+AE0og79cEPqnKMNChaF0ikEWLI+47KjfEQEAiFOO2QRFWMKs3gbiAmkaC6E2nxCK+GHTPnDKjyhV
qe88XzGfxa0jnXLrqKa8/CCiBCsfRYQgTvmXNV9StOYoTrnfquv5LY1iLRMlqyvWslCyNcVcE2KuizHX
BZlro8x1YebaOHNdoLk2Uh4n2BvECZ4vVmouxD36YK9qixvDUskSnK3Myq+VZ4oCq+6qM2eMJutRLOdb
R5ssW0cbLb9SSjRmQixt1AZWtztP+ReNnRBvfNe5wEZv5Le2WWlrIuXbrLT1kPKNVjexPi0p32qlTUxN
bLXTpqYmttppk1MTG+3qXH0WS5slczXG12pnqJurCenAsqqGC9kVV5p22UE91RZdiItmXMq/I+IU9+d/
AgAA//9wgb16ZhgAAA==
`,
	},

	"/templates/default/scan_type_map.json": {
		local:   "templates/default/scan_type_map.json",
		size:    622,
		modtime: 1528164482,
		compressed: `
H4sIAAAAAAAC/3zRv2rDMBDH8d1PITQbg1sjiscOBc9NJpEhgSRcUE4QS1PIuwf9udNl8fi1P/yG07NT
Sl+cP4bvLz0rpSxXrzRG54a/2oeerZmkNdOnNVO1J+9dhtnmIviboijA8NNULlJLiqKiYLYWsT00BxhG
w66UmBuN2CNoqeQgS0C6TV0Ut1mwXSYKaKnkIktAumBdFBdcsN0vCmip5CLLAPdzu2GqYZc+Ec5RHwUC
SWX1Gh6AV3b/JYu8rR7b6Jas/+ZN2b26dwAAAP//G8mEhG4CAAA=
`,
	},

	"/templates/default/stmt.tmpl": {
		local:   "templates/default/stmt.tmpl",
		size:    8676,
		modtime: 1530329395,
		compressed: `
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
798BAAD//5V9NJHkIQAA
`,
	},

	"/templates/default/table.tmpl": {
		local:   "templates/default/table.tmpl",
		size:    8397,
		modtime: 1530329044,
		compressed: `
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
//9QMzPZzSAAAA==
`,
	},

	"/templates/default/test.tmpl": {
		local:   "templates/default/test.tmpl",
		size:    2130,
		modtime: 1529627540,
		compressed: `
H4sIAAAAAAAC/7RU32/TMBB+jv+KIxKVXZWO3N4m9QEBDzwwKpjGwzQhz3U2a6ljORfKVOV/R3YSmm3p
BGPdw+qzv/vx3d0XJ9WtvNaw3cJ82Z5P5VpD0zBm1q70BJwlab6mlAnGjo7grKIMTAWkK4IzeVXor+UG
NoZugDblG1UW9dqC82Yt/R3c6rs5ozunW7+KfK0ItixR0TD2miUOwVg6Rpa4rD/J/tT0SXE8aVnTeDIc
T6ZwFzhCqaJs2QY4l0WtB267ctzA66f0oSc/dnUM/hbAp4Gp4NYUYoj6buimy/MEaiQWdijBWF5bBZw8
tN5gfy29/qak5avQmOnFpbGkfS6V3jYiUJjGhwVI57Rd8WjOYEJ+rrL212H329kyE4HlSKqyophLgPa+
9CG811R7C9YUg86Ejn7WJGEBp3oTmQWTsyRJw1s6Y0lycdlOZJuqLJ1B6jD+j2eZpU3AfHHUdex93KqK
dwCHqegA3Uvo7gedy7ogvovSg97VVH6yqsXy/mW0pbtyBUz/GAOuPbuxJp3Lwqy4gKuyLIYuHl4t+iY9
9GmLiqvHTdgyAYMhhijVxpC6ARM3WVYa3p6wZBd7rrLuPrt/77C7xwf3Pf74/r3MWLJqexgenLRG8XxN
849h3DlPX6+CBoPiyhy8tNc6nYERgiXNfmbLMtJ5DrfJPnKTfewm++hNDsBv+NXgj5iZPMx9Eece7KFW
kubeOg0DBaTLTto5hfV12Bo4GysDX+YjoFDsCf0PoscnRI+jolcYhf5YiPh3QsRRIeIzhIiHFKLCF1w8
PLSw/rPa3wEAAP//N0x8k1IIAAA=
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
}
