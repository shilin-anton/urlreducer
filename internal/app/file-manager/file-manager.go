package filemanager

import (
	"bufio"
	"encoding/json"
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
	"log"
	"os"
	"strconv"
)

type record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileWriter struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

type FileReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewWriter() (*FileWriter, error) {
	file, err := os.OpenFile(config.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file:    file,
		scanner: bufio.NewScanner(file),
		writer:  bufio.NewWriter(file),
	}, nil
}

func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

func NewReader() (*FileReader, error) {
	file, err := os.OpenFile(config.FilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (fr *FileReader) Close() {
	fr.file.Close()
}

func ReadFromFile(storage *storage.Storage) {
	if localStorageDisabled() {
		return
	}

	reader, err := NewReader()
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer reader.Close()

	for reader.scanner.Scan() {
		line := reader.scanner.Bytes()

		rec := &record{}
		if err := json.Unmarshal(line, &rec); err != nil {
			log.Fatal("Error decoding data from file:", err)
		}
		storage.Add(rec.ShortURL, rec.OriginalURL)
	}

	if err := reader.scanner.Err(); err != nil {
		log.Fatal("Error scanning from file:", err)
	}
}

func AddRecord(short string, url string, uuid int) error {
	if localStorageDisabled() {
		return nil
	}

	writer, err := NewWriter()
	if err != nil {
		return err
	}
	defer writer.Close()

	newRecord := record{
		UUID:        strconv.Itoa(uuid),
		ShortURL:    short,
		OriginalURL: url,
	}
	recordJSON, err := json.Marshal(newRecord)
	if err != nil {
		return err
	}
	if _, err := writer.file.WriteString(string(recordJSON) + "\n"); err != nil {
		return err
	}

	return nil
}

func localStorageDisabled() bool {
	return config.FilePath == ""
}
