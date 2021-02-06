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

var _sampleFilehiveConf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x5f\x6b\xd5\x4c\x10\x87\xef\xf3\x29\xe6\x03\x84\x34\x6d\x79\x5f\x4b\x43\x2e\x0a\x22\x14\x91\x16\x2a\xe2\xed\x24\x3b\x49\x06\x37\x3b\xe9\xce\xe6\xc4\x83\xe8\x67\x97\xdd\xfc\xb1\x55\x04\x39\x37\xc9\x99\x79\x9e\xfd\xcd\x6c\x2a\xf8\x38\x10\x18\xf6\xd4\x06\xf1\x67\x08\x02\x1a\xc4\x13\x18\x0c\x08\x3a\xb7\x03\xa0\x42\x18\x08\x3a\xb6\x34\xf0\x69\xad\x34\xa8\x54\x64\x1b\x4c\x1d\xce\x36\x00\x2b\xfc\xb8\x28\x8e\x36\x71\xf0\xf8\xf0\x74\xff\x19\x1e\x9e\x48\x8b\x2c\xab\xe0\x2d\x35\x73\x0f\x56\xfa\x9e\x5d\x0f\x96\x4e\x64\xa3\xe3\x13\x5a\x36\xeb\xab\x02\x7a\x82\x6f\x26\x36\xe6\xc0\xae\x93\x1c\x9c\x04\x6e\x29\x87\x05\xbd\x63\xd7\xe7\x40\xde\x8b\xcf\xa1\xf5\x1c\xb8\x45\xfb\x3d\xab\xa2\x33\xf1\x75\x44\xb2\xec\xaf\x43\x59\xe9\xd3\x1c\xba\x32\x86\x7d\xfd\x22\xf2\x85\x95\x5e\x23\x7d\xf7\x9a\x9d\x95\x92\x02\x4f\x04\x1a\x30\x70\xbb\x4a\x8e\xf5\xf0\x88\x3d\x69\xec\x41\x67\x40\xc9\x9f\xc8\xc7\x9d\x8d\xd0\x79\x19\xe3\x8c\x2b\x16\xa9\xdf\xcf\x5c\x96\xe5\x08\x2c\x23\xb2\x4b\xcb\xde\x1c\x0b\x5b\x0b\x7e\x76\x20\x2e\xab\xb6\x7a\x4d\x5f\x71\x9c\x2c\x15\xad\x8c\x3b\xc9\x2e\x90\xef\xb0\xa5\xdb\x49\x7c\x80\x4e\xd2\xf1\xb0\x50\x73\xa4\x11\x68\xd8\x19\x08\x12\xe3\x58\xd6\x40\xae\x2e\x8b\xf4\xbb\xbd\x29\x6f\xca\x55\xc5\x1a\xaf\x31\xb2\xa6\x81\x70\x9e\xa8\x80\xfb\x00\x2d\x3a\x20\x0e\x03\x79\x68\x08\xf4\xd9\x72\xa0\xeb\x1c\xc6\xb3\x3e\xdb\x1c\xc4\xc3\x24\x1a\x7a\x4f\xaa\x51\x6e\x1a\xc3\x68\xa9\x0d\x75\x6a\xd8\x33\x0e\xa2\x61\x97\xc7\xe7\xd7\x51\x53\x2b\x6c\x2f\x87\x6e\x4b\xbf\x5a\x23\x54\x5f\x5e\xbd\x49\xa1\x2f\x6f\xaf\xaf\xcb\xff\x77\xf7\xac\xe4\x1d\x8e\xf4\xa7\xee\x97\xca\x34\xab\x26\xf6\xd6\x3b\xb0\x0b\x26\x54\x5d\xc4\x9b\x7f\x11\xc4\xde\x7a\x07\x76\x01\x9a\x31\x5e\x9d\x7c\x21\x97\x1c\x8f\xb2\x90\xef\x31\x50\x56\xc1\xb4\x3f\xa7\x72\xbd\x23\xef\xd8\x52\x2b\xec\x00\x8d\x49\x07\x44\x6e\xc2\xb3\xcc\x21\x7e\x9f\xdd\x56\xde\xaa\x09\x3b\xac\x69\x83\x69\x80\x17\xfa\x17\xcb\xf9\xaf\x2c\xaf\x22\xf0\x01\xd9\xf6\xb3\x83\xbb\xc7\x7b\x78\x4f\xe7\xac\x82\x71\xfd\xa7\xfe\x19\x00\x00\xff\xff\xce\x92\x4d\x40\xf9\x03\x00\x00")

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

	info := bindataFileInfo{name: "sample-filehive.conf", size: 1017, mode: os.FileMode(420), modTime: time.Unix(1612629854, 0)}
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
