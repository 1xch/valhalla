package image

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	_ "image/jpeg"
	_ "image/png"
)

func basicErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", e)
		os.Exit(-1)
	}
}

type xrror struct {
	base string
	vals []interface{}
}

func (x *xrror) Error() string {
	return fmt.Sprintf("%s", fmt.Sprintf(x.base, x.vals...))
}

func (x *xrror) Out(vals ...interface{}) *xrror {
	x.vals = vals
	return x
}

func Xrror(base string) *xrror {
	return &xrror{base: base}
}

func openImage(path string) (image.Image, error) {
	file, err := openFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	i, err := decodeImage(file)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func exist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir|0755)
	}
}

var openFileError = Xrror("unable to find or open file %s, provided %s").Out

//
func openFile(path string) (*os.File, error) {
	p := filepath.Clean(path)

	dir, name := filepath.Split(p)

	var fp string
	var err error
	switch dir {
	case "":
		fp, err = filepath.Abs(name)
	default:
		exist(dir)
		fp, err = filepath.Abs(p)
	}
	if err != nil {
		return nil, err
	}

	if file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0660); err == nil {
		return file, nil
	}

	return nil, openFileError(fp, path)
}

func decodeImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}
