// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/nlpodyssey/gopickle/types"
	"io"
	"math/big"
	"os"
	"path"
)

const hexMagicNumber = "1950a86a20f9469cfc6c"
const protocolVersion = 1001

var ErrInvalidMagicNumber = errors.New("invalid pytorch magic number")
var ErrInvalidProtocolVersion = errors.New("invalid pytorch protocol version")

func Load(filename string) (interface{}, error) {
	newUnpickler := func(r io.Reader) pickle.Unpickler {
		return pickle.NewUnpickler(r)
	}
	return LoadWithUnpickler(filename, newUnpickler)
}

// LoadWithUnpickler is like Load, but it accepts a newUnpickler function which
// is used to create new customized pickle.Unpickler instances.
func LoadWithUnpickler(filename string, newUnpickler func(r io.Reader) pickle.Unpickler) (interface{}, error) {
	if !isZipFile(filename) {
		return loadLegacyFile(filename, newUnpickler)
	}
	return loadZipFile(filename, newUnpickler)
}

func loadZipFile(filename string, newUnpickler func(r io.Reader) pickle.Unpickler) (interface{}, error) {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	fileRecords := make(map[string]*zip.File, len(r.File))
	for _, f := range r.File {
		_, recordName := path.Split(f.Name)
		fileRecords[recordName] = f
	}

	if _, isTorchScript := fileRecords["constants.pkl"]; isTorchScript {
		return nil, fmt.Errorf("TorchScript is not supported")
	}

	dataFile, hasDataFile := fileRecords["data.pkl"]
	if !hasDataFile {
		return nil, fmt.Errorf("data.pkl not found in zip file")
	}
	df, err := dataFile.Open()
	if err != nil {
		return nil, err
	}
	defer df.Close()

	loadedStorages := make(map[string]StorageInterface)

	u := newUnpickler(df)
	u.FindClass = makePickleFindClass(u.FindClass)
	u.PersistentLoad = func(savedId interface{}) (interface{}, error) {
		tuple, tupleOk := savedId.(*types.Tuple)
		if !tupleOk || tuple.Len() == 0 {
			return nil, fmt.Errorf("PersistentLoad: non-empty tuple expected, got %#v", savedId)
		}
		typename, typenameOk := tuple.Get(0).(string)
		if !typenameOk {
			return nil, fmt.Errorf("PersistentLoad: cannot get typename")
		}
		if typename != "storage" {
			return nil, fmt.Errorf("unknown typename for PersistentLoad, expected 'storage' but got '%s'", typename)
		}
		if tuple.Len() < 5 {
			return nil, fmt.Errorf("PersistentLoad: unexpected storage data length")
		}
		dataType, dataTypeOk := tuple.Get(1).(StorageClassInterface)
		key, keyOk := tuple.Get(2).(string)
		location, locationOk := tuple.Get(3).(string)
		size, sizeOk := tuple.Get(4).(int)
		if !dataTypeOk || !keyOk || !locationOk || !sizeOk {
			return nil, fmt.Errorf("PersistentLoad: unexpected data types")
		}
		storage, storageExists := loadedStorages[key]
		if !storageExists {
			storage, err = loadTensor(dataType, size, location, key, fileRecords)
			if err != nil {
				return nil, err
			}
			loadedStorages[key] = storage
		}
		return storage, nil
	}
	return u.Load()
}

func loadTensor(
	dataType StorageClassInterface,
	size int,
	location, key string,
	zipFileRecords map[string]*zip.File,
) (StorageInterface, error) {
	file, fileOk := zipFileRecords[key]
	if !fileOk {
		return nil, fmt.Errorf("cannot find zip record '%s'", key)
	}
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	storage := dataType.New(size, location)
	err = storage.SetFromFileWithSize(f, size)
	return storage, err
}

func loadLegacyFile(filename string, newUnpickler func(r io.Reader) pickle.Unpickler) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tr := tar.NewReader(f)
	for {
		_, err := tr.Next()
		switch err {
		case nil:
			// TODO: ...
			panic("legacy load from tar not implemented")
		case io.EOF:
			break // End of archive
		case tar.ErrHeader, io.ErrUnexpectedEOF:
			_, err = f.Seek(0, io.SeekStart)
			if err != nil {
				return nil, err
			}
			return loadLegacyNoTar(f, newUnpickler)
		default:
			return nil, err
		}
	}
}

func loadLegacyNoTar(f *os.File, newUnpickler func(r io.Reader) pickle.Unpickler) (interface{}, error) {
	if err := readAndCheckMagicNumber(f); err != nil {
		return nil, err
	}
	if err := readAndChecProtocolVersion(f); err != nil {
		return nil, err
	}
	if _, err := unpickle(f); err != nil { // sys info
		return nil, err
	}

	deserializedObjects := make(map[string]StorageInterface)

	u := newUnpickler(f)
	u.FindClass = makePickleFindClass(u.FindClass)
	u.PersistentLoad = func(savedId interface{}) (interface{}, error) {
		tuple, tupleOk := savedId.(*types.Tuple)
		if !tupleOk || tuple.Len() == 0 {
			return nil, fmt.Errorf("PersistentLoad: non-empty tuple expected, got %#v", savedId)
		}
		typename, typenameOk := tuple.Get(0).(string)
		if !typenameOk {
			return nil, fmt.Errorf("PersistentLoad: cannot get typename")
		}

		switch typename {
		case "storage":
			if tuple.Len() < 6 {
				return nil, fmt.Errorf(
					"PersistentLoad: unexpected storage data length")
			}
			dataType, dataTypeOk := tuple.Get(1).(StorageClassInterface)
			rootKey, rootKeyOk := tuple.Get(2).(string)
			location, locationOk := tuple.Get(3).(string)
			size, sizeOk := tuple.Get(4).(int)
			viewMetadata := tuple.Get(5)
			if !dataTypeOk || !rootKeyOk || !locationOk || !sizeOk {
				return nil, fmt.Errorf("PersistentLoad: unexpected data types")
			}
			storage, storageExists := deserializedObjects[rootKey]
			if !storageExists {
				storage = dataType.New(size, location)
				deserializedObjects[rootKey] = storage
			}
			switch vm := viewMetadata.(type) {
			case nil:
				return storage, nil
			case []interface{}:
				if len(vm) != 3 {
					return nil, fmt.Errorf(
						"PersistentLoad: unexpected view metadata length")
				}
				panic("viewMetadata not implemented")
				// TODO: ...
				// view_key, offset, view_size = view_metadata
				// if view_key not in deserialized_objects:
				//     deserialized_objects[view_key] = storage[offset:offset + view_size]
				// return deserialized_objects[view_key]
			default:
				return nil, fmt.Errorf("PersistentLoad: unexpected view metadata type")
			}
		case "module":
			if tuple.Len() < 2 {
				return nil, fmt.Errorf("PersistentLoad: unexpected module data length")
			}
			return tuple.Get(1), nil
		default:
			return nil, fmt.Errorf("Unexpected saved ID type: %s", typename)
		}
	}
	result, err := u.Load()
	if err != nil {
		return nil, err
	}

	rawStorageKeys, err := unpickle(f)
	if err != nil {
		return nil, err
	}
	storageKeys, err := makeStorageKeys(rawStorageKeys)
	if err != nil {
		return nil, err
	}

	for _, key := range storageKeys {
		storageObj, ok := deserializedObjects[key]
		if !ok {
			return nil, fmt.Errorf("storage object not found for key '%s'", key)
		}
		err = storageObj.SetFromFile(f)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func makeStorageKeys(obj interface{}) ([]string, error) {
	list, ok := obj.(*types.List)
	if !ok {
		return nil, fmt.Errorf("invalid storage keys data")
	}
	keys := make([]string, len(*list))
	for i, rawKey := range *list {
		key, keyOk := rawKey.(string)
		if !keyOk {
			return nil, fmt.Errorf("invalid storage key")
		}
		keys[i] = key
	}
	return keys, nil
}

func readAndCheckMagicNumber(r io.Reader) error {
	obj, err := unpickle(r)
	if err != nil {
		return err
	}
	if n, ok := obj.(*big.Int); !ok || n.Text(16) != hexMagicNumber {
		return ErrInvalidMagicNumber
	}
	return nil
}

func readAndChecProtocolVersion(r io.Reader) error {
	obj, err := unpickle(r)
	if err != nil {
		return err
	}
	if n, ok := obj.(int); !ok || n != protocolVersion {
		return ErrInvalidProtocolVersion
	}
	return nil
}

func unpickle(r io.Reader) (interface{}, error) {
	u := pickle.NewUnpickler(r)
	return u.Load()
}

func isZipFile(filename string) bool {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return false
	}
	r.Close()
	return true
}

func makePickleFindClass(fallback func(module, name string) (interface{}, error)) func(module, name string) (interface{}, error) {
	return func(module, name string) (interface{}, error) {
		switch module + "." + name {
		case "torch._utils._rebuild_tensor_v2":
			return &RebuildTensorV2{}, nil
		case "torch.FloatStorage":
			return &FloatStorageClass{}, nil
		case "torch.HalfStorage":
			return &HalfStorageClass{}, nil
		case "torch.DoubleStorage":
			return &DoubleStorageClass{}, nil
		case "torch.CharStorage":
			return &CharStorageClass{}, nil
		case "torch.ShortStorage":
			return &ShortStorageClass{}, nil
		case "torch.IntStorage":
			return &IntStorageClass{}, nil
		case "torch.LongStorage":
			return &LongStorageClass{}, nil
		case "torch.ByteStorage":
			return &ByteStorageClass{}, nil
		case "torch.BoolStorage":
			return &BoolStorageClass{}, nil
		case "torch.nn.backends.thnn._get_thnn_function_backend":
			// this is for historical pickle deserilaization, it is not used otherwise
			return getThnnFunctionBackend{}, nil
		default:
			if fallback == nil {
				return nil, fmt.Errorf("class not found: %s %s", module, name)
			}
			return fallback(module, name)
		}
	}
}

// getThnnFunctionBackend is for historical pickle deserilaization, it is not used otherwise
type getThnnFunctionBackend struct{}

var _ types.Callable = &getThnnFunctionBackend{}

func (getThnnFunctionBackend) Call(_ ...interface{}) (interface{}, error) {
	return nil, nil
}
