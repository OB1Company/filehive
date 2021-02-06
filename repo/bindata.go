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

var _sampleFilehiveConf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\xdf\x6a\xdc\x3a\x10\x87\xef\xfd\x14\xf3\x00\x66\xb3\x49\x38\xa7\x21\xc6\x17\x81\xb6\x10\x4a\x49\x20\xa5\xf4\x76\x6c\x8d\xed\x21\xb2\xc6\xd1\xc8\xeb\x9a\xd2\x3e\x7b\x91\xfc\xa7\x9b\x96\x42\xd9\x1b\x6b\x67\xbe\x4f\x3f\x8d\x54\xc0\xa7\x8e\xc0\xb0\xa7\x3a\x88\x9f\x21\x08\x68\x10\x4f\x60\x30\x20\xe8\x58\x77\x80\x0a\xa1\x23\x68\xd8\x52\xc7\xa7\xa5\x52\xa1\xd2\x21\x5b\x61\x6a\x70\xb4\x01\x58\xe1\xc7\xc5\x61\x6f\x13\x07\x8f\x0f\x4f\xf7\x5f\xe0\xe1\x89\xf4\x90\x65\x05\xbc\xa5\x6a\x6c\xc1\x4a\xdb\xb2\x6b\xc1\xd2\x89\x6c\x74\x7c\x46\xcb\x66\x59\x2a\xa0\x27\xf8\x66\x62\x63\x0e\xec\x1a\xc9\xc1\x49\xe0\x9a\x72\x98\xd0\x3b\x76\x6d\x0e\xe4\xbd\xf8\x1c\x6a\xcf\x81\x6b\xb4\xdf\xb3\x22\x3a\x13\x5f\x46\x24\xcb\xfe\x7a\x28\x2b\x6d\x3a\x87\x2e\x8c\x61\x5f\x9e\x45\xbe\xb0\xd2\x6a\xa4\xef\x5e\xb3\xa3\x52\x52\xe0\x89\x40\x03\x06\xae\x17\xc9\x3e\x1e\xee\xb1\x25\x8d\x3d\xe8\x0c\x28\xf9\x13\xf9\x38\xb3\x1e\x1a\x2f\x7d\x3c\xe3\x82\x45\xea\xf7\x3d\xa7\x69\xda\x03\x4b\x8f\xec\xd2\xb0\x57\xc7\xc4\xd6\x82\x1f\x1d\x88\xcb\x8a\xb5\x5e\xd2\x57\xec\x07\x4b\x87\x5a\xfa\x8d\x64\x17\xc8\x37\x58\xd3\xed\x20\x3e\x40\x23\x69\x7b\x98\xa8\xda\xd3\x08\x54\xec\x0c\x04\x89\x71\x2c\x6b\x20\x57\x1e\x0f\xe9\x77\x7b\x73\xbc\x39\x2e\x2a\xd6\x78\x8d\x91\x35\x15\x84\x79\xa0\x03\xdc\x07\xa8\xd1\x01\x71\xe8\xc8\x43\x45\xa0\x2f\x96\x03\x5d\xe7\xd0\xcf\xfa\x62\x73\x10\x0f\x83\x68\x68\x3d\xa9\x46\xb9\xa9\x0c\xa3\xa5\x3a\x94\xa9\x61\xcb\xd8\x89\x86\x4d\x1e\xbf\x5f\x47\x4d\xad\xb0\x2e\x76\xdd\x9a\x7e\xb1\x46\xa8\xbc\xbc\x7a\x93\x42\x5f\xde\x5e\x5f\x1f\xff\xdf\xdc\xa3\x92\x77\xd8\xd3\x9f\xba\x5f\x2a\x53\x2d\x9a\xd8\x5b\x6e\xc0\x26\x18\x50\x75\x12\x6f\xfe\x45\x10\x7b\xcb\x0d\xd8\x04\x68\xfa\x78\x75\xf2\x4c\x2e\x39\x1e\x65\x22\xdf\x62\xa0\xac\x80\x61\xfb\x4e\xe5\x72\x43\xde\xb3\xa5\x5a\xd8\x01\x1a\x93\x36\x88\xdc\x80\xb3\x8c\x21\xbe\xcf\x66\x2d\xaf\xd5\x84\xed\xd6\x34\xc1\x74\x80\x33\xfd\xd9\x70\xfe\x3b\x1e\xaf\x22\xf0\x11\xd9\xb6\xa3\x83\xbb\xc7\x7b\xf8\x40\x73\x56\x40\xbf\xfc\xf3\x4c\x73\x32\xbe\x8b\xeb\xf5\x65\xad\xd5\xf5\x99\xfd\x0c\x00\x00\xff\xff\xe6\x8a\x11\xed\x1a\x04\x00\x00")

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

	info := bindataFileInfo{name: "sample-filehive.conf", size: 1050, mode: os.FileMode(420), modTime: time.Unix(1612645556, 0)}
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
