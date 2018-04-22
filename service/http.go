package gomoney

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joaosoft/go-error/service"
	"github.com/labstack/echo"
)

func download(key string, ctx echo.Context) ([]bytes.Buffer, *goerror.ErrorData) {
	writers := make([]bytes.Buffer, 0)

	// multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	files := form.File[key]

	for _, file := range files {
		// source
		src, err := file.Open()
		if err != nil {
			return nil, goerror.NewError(err)
		}
		defer src.Close()

		// destination
		var dst bytes.Buffer

		// copy
		dst.ReadFrom(src)

		writers = append(writers, dst)
	}
	log.Infof("uploaded successfully %d %s", len(files), key)

	return writers, nil
}

func upload(client *http.Client, url string, values map[string]io.Reader) *goerror.ErrorData {
	var b bytes.Buffer
	var err error
	writer := multipart.NewWriter(&b)

	for key, value := range values {
		var fw io.Writer
		if x, ok := value.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if file, ok := value.(*os.File); ok {
			if fw, err = writer.CreateFormFile(key, file.Name()); err != nil {
				return goerror.NewError(err)
			}
		} else {
			// Add other fields
			if fw, err = writer.CreateFormField(key); err != nil {
				return goerror.NewError(err)
			}
		}
		if _, err = io.Copy(fw, value); err != nil {
			return goerror.NewError(err)
		}

	}
	// close the multipart
	writer.Close()

	// now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return goerror.NewError(err)
	}
	// don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// submit the request
	res, err := client.Do(req)
	if err != nil {
		return goerror.NewError(err)
	}

	// check the response
	if res.StatusCode != http.StatusOK {
		return goerror.FromString(fmt.Sprintf("bad status: %s", res.Status))
	}

	return nil
}
