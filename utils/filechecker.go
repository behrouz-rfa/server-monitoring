package utils

import (
	"io"
	"mime/multipart"
	"net/http"
)

var FileChecker fileCheckerInterface = &fileChecker{}

type Sizer interface {
	Size() int64
}
type fileCheckerInterface interface {
	GetFileContentType(multipart.File) (string, error)

}
type fileChecker struct{}


func (f fileChecker) Size() int64 {
	panic("implement me")
}

func (f fileChecker) GetFileContentType(out multipart.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType,  nil
}
