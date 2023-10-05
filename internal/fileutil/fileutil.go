package fileutil

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/storage"
)

type FileIoStruct struct {
	cfg        config.RunFileIo
	memStorage *storage.MemStorage
}

var flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND

func NewFileIo(repo *storage.MemStorage, cfg config.RunFileIo) *FileIoStruct {
	return &FileIoStruct{
		memStorage: repo,
		cfg:        cfg,
	}
}

func (f FileIoStruct) SaveMetricsToFile(filePath string) error {
	file, err := os.OpenFile(filePath, flag, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	metrics := f.memStorage.MetricsReturner()
	file.Write([]byte(metrics))
	file.Write([]byte("\n"))
	return nil
}

func ClearFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte{})
	if err != nil {
		return err
	}

	return nil
}


func (f FileIoStruct) Run() {
	if f.cfg.Path == "" {
		return
	}
	if !f.cfg.Restore {
		err := ClearFile(f.cfg.Path)
		if err != nil {
			log.Println(err)
		}
	}

	var timeInterval int
	if f.cfg.Interval > 0 {
		timeInterval = f.cfg.Interval
	} else if f.cfg.Interval == 0 {
		timeInterval = 10
	}
		go func() {
		defer runtime.Goexit()
		ticker := time.NewTicker(time.Duration(timeInterval) * time.Second)
		defer ticker.Stop()
		for range ticker.C{
				err := f.SaveMetricsToFile(f.cfg.Path)
				if err != nil {
					log.Println("ошибка при сохранении текущих значений метрик:", err)
				}
				ticker.Reset(time.Duration(timeInterval) * time.Second)

		}
	}()
}
