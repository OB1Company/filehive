// Code generated by go-bindata.
// sources:
// sample-filehive.conf
// DO NOT EDIT!

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _sampleFilehiveConf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x90\x41\x8b\x9c\x40\x10\x85\xef\xfe\x8a\xf7\x03\xc4\x9d\xcd\x40\xb2\xcc\xe0\x2d\x97\x3d\x6d\x60\x43\xc8\xb5\xb5\x4b\x2d\x68\xbb\x9c\xaa\x76\x26\x43\x48\x7e\x7b\x68\x1d\x27\x09\x21\xb0\x78\xb1\xf0\x7d\x1f\xef\x79\xc4\xe7\x81\xe0\x59\xa9\x4d\xa2\x57\x24\x81\x25\x51\x82\x77\xc9\xc1\xe6\x76\x80\x33\xa4\x81\xd0\x71\xa0\x81\xcf\xeb\x97\xc6\x19\x55\xc5\x0d\xa6\xce\xcd\x21\x81\x0d\x3f\x1f\xaa\x7b\x4c\x22\x3e\xbd\xbc\x3e\x7f\xc5\xcb\x2b\x59\x55\x14\x47\x7c\xa4\x66\xee\x11\xa4\xef\x39\xf6\x08\x74\xa6\x90\x1d\x5f\x5c\x60\xbf\x9e\x06\xa7\x84\xef\x3e\x07\x4b\x70\xec\xa4\x44\x94\xc4\x2d\x95\xb8\x38\x8d\x1c\xfb\x12\xa4\x2a\x5a\xa2\x55\x4e\xdc\xba\xf0\xa3\x38\x66\xe7\xc2\xd7\x19\x29\x8a\xff\x8e\x0a\xd2\x2f\x3b\x6c\x65\x3c\x6b\xfd\x47\xe5\x87\x20\xbd\xdd\x69\x19\x1d\xc7\x65\xb9\x91\x9e\x49\x71\xe1\x10\xa0\x73\x84\xc4\xe2\x08\xfa\xe6\xc6\x29\x50\xd5\xca\xb8\x21\x1c\x13\x69\xe7\x5a\x3a\x4c\xa2\x09\x9d\xe8\x82\x5f\xa8\xd9\x14\x49\xd0\x70\xf4\x48\x92\x87\xef\xaa\xe5\x39\x3c\xed\x9e\x76\xab\x83\x2d\xff\xc5\x0c\xf9\x06\xe9\x3a\x51\x85\xe7\x84\xd6\x45\x10\xa7\x81\x14\x0d\xc1\x4e\x81\x13\xed\x4b\x8c\x57\x3b\x85\x12\xa2\x98\xc4\x52\xaf\x64\x96\xad\xbe\xf1\xec\x02\xb5\xa9\x5e\x02\x5b\xb9\x41\x2c\x6d\xf2\xfc\xfe\x77\xc7\x25\x8a\xdb\x71\xd7\xdd\x6a\xaf\xd6\x0c\xd5\x8f\xef\x3e\x2c\xa5\x1f\x0f\xfb\xfd\xee\xfd\xe6\x9e\x8d\x34\xba\x91\xfe\xd5\xfd\x56\xf9\x66\xd5\xe4\x6c\xbd\x01\x9b\x60\x72\x66\x17\x51\xff\x16\x41\xce\xd6\x1b\x50\xfc\x0a\x00\x00\xff\xff\xad\xcd\xf3\xe3\xc0\x02\x00\x00")

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

	info := bindataFileInfo{name: "sample-filehive.conf", size: 704, mode: os.FileMode(420), modTime: time.Unix(1605116031, 0)}
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
