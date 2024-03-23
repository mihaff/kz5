package filestorage

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileStorage struct {
	UploadsDir   string
	DownloadsDir string
}

func (fs *FileStorage) NewFileStorage(uploadsDir, downloadsDir string) error {
	_, err := os.Stat(uploadsDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			return err
		}
		log.Println("Directory", uploadsDir, "created successfully.")
	} else if err != nil {
		return err
	}

	_, err = os.Stat(downloadsDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(downloadsDir, 0755); err != nil {
			return err
		}
		log.Println("Directory", downloadsDir, "created successfully.")
	} else if err != nil {
		return err
	}

	fs.UploadsDir = uploadsDir
	fs.DownloadsDir = downloadsDir
	return nil
}

func (fs *FileStorage) saveFile(filePath string, file io.Reader) error {
	log.Printf("Saving physical file %s", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Panic(err)
		}
	}()

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) SaveDownloadedFile(fileName string, file io.Reader) error {
	filePath := filepath.Join(fs.DownloadsDir, fileName)
	return fs.saveFile(filePath, file)
}

func (fs *FileStorage) SaveUploadedFile(fileName string, file io.Reader) error {
	filePath := filepath.Join(fs.UploadsDir, fileName)
	return fs.saveFile(filePath, file)
}

func (fs *FileStorage) deleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) DeleteUploadedFile(fileName string) error {
	filePath := filepath.Join(fs.UploadsDir, fileName)
	return fs.deleteFile(filePath)
}

func (fs *FileStorage) DeleteDownloadedFile(fileName string) error {
	filePath := filepath.Join(fs.DownloadsDir, fileName)
	return fs.deleteFile(filePath)
}

func (fs *FileStorage) GetUploadedFilePath(fileName string) string {
	return filepath.Join(fs.UploadsDir, fileName)
}

func (fs *FileStorage) GetDownloadedFilePath(fileName string) string {
	return filepath.Join(fs.DownloadsDir, fileName)
}
