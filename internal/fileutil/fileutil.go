package fileutil

import (
	"log"
	"os"
	"time"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/storage"
)

type FileIoStruct struct {
	cfg        config.RunFileIo
	memStorage *storage.MemStorage
	cfgAgent   config.RunConfig
}

func NewFileIo(repo *storage.MemStorage, cfg config.RunFileIo) *FileIoStruct {
	return &FileIoStruct{
		memStorage: repo,
		cfg:        cfg,
	}
}

func (f FileIoStruct) SaveMetricsToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	metrics := f.memStorage.MetricsReturner()
	file.Write([]byte(metrics))
	file.Write([]byte("\n"))
	return nil
}

func LoadMetricsFromFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		// Если файл не найден, просто возвращаемся без ошибки
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	return nil
}

func (f FileIoStruct) Run() {
	// Загружаем предыдущие значения метрик из файла при необходимости
	if f.cfg.Restore {
		err := LoadMetricsFromFile(f.cfg.Path)
		if err != nil {
			log.Println("ошибка при загрузке предыдущих значений метрик:", err)
		}
	}

	var timeInterval int
	if f.cfg.Interval > 0 {
		timeInterval = f.cfg.Interval
	} else if f.cfg.Interval == 0 {
		timeInterval = f.cfgAgent.DefaultReportInterval
	}
	// Запускаем горутину для периодического сохранения метрик на диск
	go func() {
		ticker := time.NewTicker(time.Duration(timeInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := f.SaveMetricsToFile(f.cfg.Path)
				if err != nil {
					log.Println("ошибка при сохранении текущих значений метрик:", err)
				}
				ticker.Reset(time.Duration(timeInterval) * time.Second)
			case <-f.cfg.ShutdownChan:
				// При завершении сервера сохраняем данные
				err := f.SaveMetricsToFile(f.cfg.Path)
				if err != nil {
					log.Println("ошибка при сохранении текущих значений метрик:", err)
				}
				return
			}
		}
	}()

}
