package fileutil

import (
	"encoding/json"
	"github.com/KillReall666/yaproject/internal/logger"
	"io"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/storage"
)

type FileIoStruct struct {
	cfg        config.RunConfig
	memStorage *storage.MemStorage
	logger     *logger.Logger
}

func NewFileIo(cfg config.RunConfig, store *storage.MemStorage, log *logger.Logger) *FileIoStruct {
	return &FileIoStruct{
		cfg:        cfg,
		memStorage: store,
		logger:     log,
	}
}

func (f *FileIoStruct) SaveMetricsToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|syscall.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	metrics, err := f.memStorage.ToJSON()
	if err != nil {
		f.logger.LogInfo("ошибка при преобразовании в формат JSON:", err)
	}
	file.Write(metrics)
	file.Write([]byte("\n"))
	return nil
}

func (f *FileIoStruct) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		f.logger.LogInfo("ошибка при чтении из файла")
		return err
	}

	err = json.Unmarshal(data, &f.memStorage)
	if err != nil {
		f.logger.LogInfo("ошибка при анмаршале файла")
		return err
	}

	err = f.memStorage.UnmarshalJSONData(data)
	if err != nil {
		f.logger.LogInfo("Ошибка при распаковке JSON:", err)
	}
	return nil
}

func (f *FileIoStruct) Run() {
	if f.cfg.Path == "" {
		return
	}

	if f.cfg.Restore {
		err := f.LoadFromFile(f.cfg.Path)
		if err != nil {
			f.logger.LogInfo("ошибка при загрузке из файла: ", err)
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
		for range ticker.C {
			err := f.SaveMetricsToFile(f.cfg.Path)
			if err != nil {
				f.logger.LogInfo("ошибка при сохранении текущих значений метрик:", err)
			}
			ticker.Reset(time.Duration(timeInterval) * time.Second)
		}
	}()
}
