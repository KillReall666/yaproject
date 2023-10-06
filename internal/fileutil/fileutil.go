package fileutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/storage"
)

type FileIoStruct struct {
	cfg        config.RunFileIo
	memStorage *storage.MemStorage
}

func NewFileIo(repo *storage.MemStorage, cfg config.RunFileIo) *FileIoStruct {
	return &FileIoStruct{
		memStorage: repo,
		cfg:        cfg,
	}
}

func (f FileIoStruct) SaveMetricsToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|syscall.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	metrics, err := f.memStorage.ToJSON()
	if err != nil {
		log.Fatal("ошибка при преобразовании в формат JSON:", err)
	}
	fmt.Println(string(metrics))

	file.Write(metrics)
	file.Write([]byte("\n"))
	return nil
}

func (f *FileIoStruct) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("ошибка при открытии файла")
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("ошибка при чтении из файла")
		return err
	}

	err = json.Unmarshal(data, &f.memStorage)
	if err != nil {
		log.Println("ошибка при анмаршале файла")
		return err
	}

	err = f.memStorage.UnmarshalJSONData(data)
	if err != nil {
		fmt.Println("Ошибка при распаковке JSON:", err)

	}
	return nil
}

func (f FileIoStruct) Run() {
	if f.cfg.Path == "" {
		return
	}

	if f.cfg.Restore {
		err := f.LoadFromFile(f.cfg.Path)
		if err != nil {
			log.Println("Ошибка при загрузке из файла:", err)
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
				log.Println("ошибка при сохранении текущих значений метрик:", err)
			}
			ticker.Reset(time.Duration(timeInterval) * time.Second)
		}
	}()
}
