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

var _sampleFilehiveConf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x41\x6b\xdc\x3c\x10\x86\xef\xfe\x15\xf3\x03\x8c\xb3\xf9\x02\x5f\xc3\x2e\x3e\x14\x4a\x21\xa7\x04\x52\x4a\xaf\x63\x6b\x6c\x0f\x95\x35\x8e\x66\xbc\xee\x52\xda\xdf\x5e\x24\xdb\xdb\xa6\xa5\x50\xf6\x62\x31\xef\xf3\xe8\x95\xb4\x27\xf8\x30\x10\x38\x8e\xd4\x9a\xc4\x0b\x98\x80\x9a\x44\x02\x87\x86\xa0\x73\x3b\x00\x2a\xd8\x40\xd0\xb1\xa7\x81\xcf\xeb\xa4\x41\xa5\xaa\xd8\x60\xea\x70\xf6\x06\xac\xf0\xfd\xa6\xba\xc6\x24\xc0\xd3\xe3\xf3\xc3\x27\x78\x7c\x26\xad\x8a\xe2\x04\xef\xa8\x99\x7b\xf0\xd2\xf7\x1c\x7a\xf0\x74\x26\x9f\x1c\x1f\xd1\xb3\x5b\x97\x0a\x18\x09\xbe\xba\x14\x2c\x81\x43\x27\x25\x04\x31\x6e\xa9\x84\x05\x63\xe0\xd0\x97\x40\x31\x4a\x2c\xa1\x8d\x6c\xdc\xa2\xff\x56\x9c\x92\x33\xf3\x75\x42\x8a\xe2\xaf\x87\xf2\xd2\xe7\x73\xe8\xca\x38\x8e\xf5\x2f\x95\x6f\xbc\xf4\x9a\xe8\xb7\xaf\xd9\x59\x29\x2b\xf0\x4c\xa0\x86\xc6\xed\x2a\xb9\x5e\x0f\x8f\xd8\x93\xa6\x0c\x06\x07\x4a\xf1\x4c\x31\xdd\xd9\x08\x5d\x94\x31\x9d\x71\xc5\x12\xf5\xfb\x9e\xcb\xb2\x5c\x0b\xcb\x88\x1c\xf2\x65\x6f\x8e\x85\xbd\x87\x38\x07\x90\x50\x9c\xb6\x79\x4d\x5f\x70\x9c\x3c\x55\xad\x8c\x3b\xc9\xc1\x28\x76\xd8\xd2\x71\x92\x68\xd0\x49\xde\x1e\x16\x6a\xae\x6d\x04\x1a\x0e\x0e\x4c\x52\x1d\xcf\x6a\x14\xea\x43\x95\x7f\xc7\xfb\xc3\xfd\x61\x55\xb1\xa6\x67\x4c\xac\x6b\xc0\x2e\x13\x55\xf0\x60\xd0\x62\x00\x62\x1b\x28\x42\x43\xa0\x2f\x9e\x8d\xee\x4a\x18\x2f\xfa\xe2\x4b\x90\x08\x93\xa8\xf5\x91\x54\x93\xdc\x35\x8e\xd1\x53\x6b\x75\x0e\xec\x1d\x07\x51\xdb\xe5\xe9\xfb\x75\xd5\x1c\x85\x6d\x71\xd5\x6d\xed\x57\x6b\x82\xea\xdb\xff\xde\xe4\xd2\xb7\xc7\xbb\xbb\xc3\xff\xbb\x7b\x56\x8a\x01\x47\xfa\x53\xf7\x53\xe5\x9a\x55\x93\xb2\xf5\x0e\xec\x82\x09\x55\x17\x89\xee\x5f\x04\x29\x5b\xef\xc0\x2e\x40\x37\xa6\xa7\x93\xcf\x14\xb2\xe3\x49\x16\x8a\x3d\x1a\x15\x27\x98\xf6\xef\x3c\xae\x77\xe4\x3d\x7b\x6a\x85\x03\xa0\x73\x79\x83\xc4\x4d\x78\x91\xd9\xd2\xff\xb3\xdb\xc6\xdb\xb4\xfe\x11\x00\x00\xff\xff\x3b\x4a\x42\xc7\xaa\x03\x00\x00")

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

	info := bindataFileInfo{name: "sample-filehive.conf", size: 938, mode: os.FileMode(420), modTime: time.Unix(1612031081, 0)}
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
