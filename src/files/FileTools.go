package files

import (
	"bike/models"
	"os"

	log "github.com/sirupsen/logrus"
)

func GetFileList(dataDirectory string) ([]*models.FileDetails, error) {
	dirEntries, err := os.ReadDir(dataDirectory)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	filenames := []*models.FileDetails{}
	for _, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			log.Error(err)
			continue
		}
		filenames = append(filenames,
			&models.FileDetails{
				Filename:   entry.Name(),
				ModifyDate: info.ModTime().String(),
				Size:       info.Size(),
			})
	}
	return filenames, nil
}
