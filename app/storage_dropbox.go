package gomoney

import (
	"github.com/joaosoft/dropbox"
)

// storageDropbox ...
type storageDropbox struct {
	conn *dropbox.Dropbox
}

// newStoragePostgres ...
func newStorageDropbox(config *dropbox.DropboxConfig) *storageDropbox {
	var conn *dropbox.Dropbox
	if config == nil {
		conn, _ = dropbox.NewDropbox()
	} else {
		conn, _ = dropbox.NewDropbox(dropbox.WithConfiguration(config))
	}
	return &storageDropbox{
		conn: conn,
	}
}

func (storage *storageDropbox) upload(path string, data []byte) error {
	_, err := storage.conn.File().Upload(path, data)
	return err
}

func (storage *storageDropbox) download(path string) ([]byte, error) {
	return storage.conn.File().Download(path)
}

func (storage *storageDropbox) delete(path string) error {
	_, err := storage.conn.Folder().DeleteFolder(path)
	return err
}
