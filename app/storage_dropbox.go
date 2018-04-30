package gomoney

import (
	godropbox "github.com/joaosoft/go-dropbox/app"
	goerror "github.com/joaosoft/go-error/app"
)

// storageDropbox ...
type storageDropbox struct {
	conn *godropbox.Dropbox
}

// newStoragePostgres ...
func newStorageDropbox(config *godropbox.DropboxConfig) *storageDropbox {
	var dropbox *godropbox.Dropbox
	if config == nil {
		dropbox = godropbox.NewDropbox()
	} else {
		dropbox = godropbox.NewDropbox(godropbox.WithConfiguration(config))
	}
	return &storageDropbox{
		conn: dropbox,
	}
}

func (storage *storageDropbox) upload(path string, data []byte) *goerror.ErrorData {
	_, err := storage.conn.File().Upload(path, data)
	return err
}

func (storage *storageDropbox) download(path string) ([]byte, *goerror.ErrorData) {
	return storage.conn.File().Download(path)
}

func (storage *storageDropbox) delete(path string) *goerror.ErrorData {
	_, err := storage.conn.Folder().DeleteFolder(path)
	return err
}
