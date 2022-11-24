package service

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"team-task/internal/dto"
	"team-task/internal/storage"
	"team-task/pkg/logger"
)

type BackupperService struct {
	storage storage.UserGrade
}

func NewBackupperService(storage storage.UserGrade) *BackupperService {
	return &BackupperService{storage: storage}
}

func (b *BackupperService) Backup() (string, error) {
	fileName := "./backup.csv"
	compressedFilename := "./backup.csv.gz"

	store := b.storage.GetAll()
	data := b.convertToSlices(store)
	if err := b.generateCSVDump(data, fileName); err != nil {
		return "", fmt.Errorf("error generation .csv dump file: %s", err.Error())
	}

	if err := b.compressFile(fileName, compressedFilename); err != nil {
		return "", fmt.Errorf("error compressing .csv dump file to .csv.gz: %s", err.Error())
	}

	return compressedFilename, nil
}

func (b *BackupperService) convertToSlices(m map[string]dto.UserGrade) [][]string {
	data := make([][]string, len(m))
	i := 0
	for _, value := range m {
		data[i] = append(data[i], value.UserId)
		data[i] = append(data[i], strconv.Itoa(value.PostpaidLimit))
		data[i] = append(data[i], strconv.Itoa(value.Spp))
		data[i] = append(data[i], strconv.Itoa(value.ShippingFee))
		data[i] = append(data[i], strconv.Itoa(value.ReturnFee))
		i++
	}

	return data
}

func (b *BackupperService) generateCSVDump(data [][]string, name string) (error) {
	csvFile, err := os.Create(name)
	if err != nil {
		logger.Errorf("failed creating file: %s", err)
	}

	csvWriter := csv.NewWriter(csvFile)

	for _, empRow := range data {
		if err := csvWriter.Write(empRow); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	if err := csvFile.Close(); err != nil {
		return err
	}

	return nil
}

func (b *BackupperService) compressFile(filename string, compressedFilename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	read := bufio.NewReader(file)
	data, err := ioutil.ReadAll(read)
	if err != nil {
		return err
	}

	file, _ = os.Create(compressedFilename)
	w := gzip.NewWriter(file)
	if _, err := w.Write(data); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}
