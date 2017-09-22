package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import "unsafe"

var stringBuffer4096 = fixedString(4096)

func vipsForeignFindLoad(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_load(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindLoadBuffer(bytes []byte) (string, error) {
	cOperationName := C.vips_foreign_find_load_buffer(
		byteArrayPointer(bytes),
		C.size_t(len(bytes)))
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindSave(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_save(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindSaveBuffer(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_save_buffer(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsInterpolateNew(name string) (*C.VipsInterpolate, error) {
	cName := C.CString(name)
	defer freeCString(cName)

	interp := C.vips_interpolate_new(cName)
	if interp == nil {
		return nil, ErrInvalidInterpolator
	}
	return interp, nil
}

func vipsOperationNew(name string) *C.VipsOperation {
	cName := C.CString(name)
	defer freeCString(cName)
	return C.vips_operation_new(cName)
}

func vipsCall(name string, options *Options) error {
	operation := vipsOperationNew(name)
	return vipsCallOperation(operation, options)
}

func vipsCallOperation(operation *C.VipsOperation, options *Options) error {
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(operation)))
	defer options.Release()

	// TODO(d): Unref the outputs

	// Set the inputs
	if options != nil {
		for _, option := range options.Options {
			if option.IsOutput {
				continue
			}

			cName := C.CString(option.Name)
			defer freeCString(cName)

			C.SetProperty(
				(*C.VipsObject)(unsafe.Pointer(operation)),
				cName,
				&option.GValue)
		}
	}

	if ret := C.vips_cache_operation_buildp(&operation); ret != 0 {
		return handleVipsError()
	}

	// Write back the outputs
	if options != nil {
		for _, option := range options.Options {
			if !option.IsOutput {
				continue
			}

			cName := C.CString(option.Name)
			defer freeCString(cName)

			C.g_object_get_property(
				(*C.GObject)(unsafe.Pointer(operation)),
				(*C.gchar)(cName),
				&option.GValue)
			option.Deserialize()
		}
	}

	return nil
}
