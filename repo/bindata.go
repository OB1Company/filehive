// Code generated for package repo by go-bindata DO NOT EDIT. (@generated)
// sources:
// sample-filehive.conf
package repo

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _sampleFilehiveConf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\xdf\x6a\xd5\x40\x10\x87\xef\xf3\x14\xf3\x00\xe1\xf4\xb4\x45\x2d\x0d\xb9\x28\x88\x50\x44\x5a\xa8\x88\xb7\x93\xec\x24\x19\xba\xd9\x49\x77\x36\x27\x06\xd1\x67\x97\xdd\xfc\xf1\x54\x11\xe4\xdc\x24\x67\xe6\xfb\xf6\x37\xb3\x29\xe0\x73\x47\x60\xd8\x53\x1d\xc4\xcf\x10\x04\x34\x88\x27\x30\x18\x10\x74\xac\x3b\x40\x85\xd0\x11\x34\x6c\xa9\xe3\xd3\x52\xa9\x50\xe9\x90\xad\x30\x35\x38\xda\x00\xac\xf0\xf3\xe2\xb0\xb7\x89\x83\xc7\x87\xa7\xfb\xaf\xf0\xf0\x44\x7a\xc8\xb2\x02\xde\x53\x35\xb6\x60\xa5\x6d\xd9\xb5\x60\xe9\x44\x36\x3a\xbe\xa0\x65\xb3\xbc\x2a\xa0\x27\xf8\x6e\x62\x63\x0e\xec\x1a\xc9\xc1\x49\xe0\x9a\x72\x98\xd0\x3b\x76\x6d\x0e\xe4\xbd\xf8\x1c\x6a\xcf\x81\x6b\xb4\x3f\xb2\x22\x3a\x13\x5f\x46\x24\xcb\xfe\x39\x94\x95\x36\xcd\xa1\x0b\x63\xd8\x97\x67\x91\x2f\xac\xb4\x1a\xe9\xbb\xd7\xec\xa8\x94\x14\x78\x22\xd0\x80\x81\xeb\x45\xb2\xaf\x87\x7b\x6c\x49\x63\x0f\x3a\x03\x4a\xfe\x44\x3e\xee\xac\x87\xc6\x4b\x1f\x67\x5c\xb0\x48\xfd\x79\xe6\x34\x4d\x7b\x60\xe9\x91\x5d\x5a\xf6\xea\x98\xd8\x5a\xf0\xa3\x03\x71\x59\xb1\xd6\x4b\xfa\x86\xfd\x60\xe9\x50\x4b\xbf\x91\xec\x02\xf9\x06\x6b\xba\x1d\xc4\x07\x68\x24\x1d\x0f\x13\x55\x7b\x1a\x81\x8a\x9d\x81\x20\x31\x8e\x65\x0d\xe4\xca\xe3\x21\xfd\x6e\x6f\x8e\x37\xc7\x45\xc5\x1a\xaf\x31\xb2\xa6\x82\x30\x0f\x74\x80\xfb\x00\x35\x3a\x20\x0e\x1d\x79\xa8\x08\xf4\xc5\x72\xa0\xeb\x1c\xfa\x59\x5f\x6c\x0e\xe2\x61\x10\x0d\xad\x27\xd5\x28\x37\x95\x61\xb4\x54\x87\x32\x35\x6c\x19\x3b\xd1\xb0\xc9\xe3\xf3\xeb\xa8\xa9\x15\xd6\x97\x5d\xb7\xa6\x5f\xac\x11\x2a\x2f\xaf\xde\xa5\xd0\x97\xb7\xd7\xd7\xc7\xb7\x9b\x7b\x54\xf2\x0e\x7b\xfa\x5b\xf7\x5b\x65\xaa\x45\x13\x7b\xcb\x0d\xd8\x04\x03\xaa\x4e\xe2\xcd\xff\x08\x62\x6f\xb9\x01\x9b\x00\x4d\x1f\xaf\x4e\x9e\xc9\x25\xc7\xa3\x4c\xe4\x5b\x0c\x94\x15\x30\x6c\xcf\xa9\x5c\x6e\xc8\x07\xb6\x54\x0b\x3b\x40\x63\xd2\x01\x91\x1b\x70\x96\x31\xc4\xef\xb3\x59\xcb\x6b\x35\x61\xbb\x35\x6d\x30\x0d\x70\xa6\x3f\x5b\xce\x9b\xe3\xf1\x2a\x02\x9f\x90\x6d\x3b\x3a\xb8\x7b\xbc\x87\x8f\x34\x67\x05\xf4\xcb\x3f\xcf\x34\x97\xbf\x02\x00\x00\xff\xff\xb6\x60\x85\x29\xfc\x03\x00\x00")

func sampleFilehiveConfBytes() ([]byte, error) {
	return bindataRead(
		_sampleFilehiveConf,
		"sample-filehive.conf",
	)
}

func sampleFilehiveConf() (*asset, error) {
	bytes, err := sampleFilehiveConfBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sample-filehive.conf", size: 1020, mode: os.FileMode(420), modTime: time.Unix(1612634376, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"sample-filehive.conf": sampleFilehiveConf,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"sample-filehive.conf": &bintree{sampleFilehiveConf, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
