package save

import (
	"fyne.io/fyne/v2/storage"
)

var folderFilter = storage.NewMimeTypeFileFilter([]string{"application/x-directory"})

func (f *SaveDialog) isDirectory() bool {
	return f.filter == folderFilter
}
