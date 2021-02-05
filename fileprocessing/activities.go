package fileprocessing

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.temporal.io/sdk/activity"
)

/**
 * Sample activities used by file processing sample workflow.
 */

type Activities struct{}

func (a *Activities) DownloadFileActivity(ctx context.Context, fileURL string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Downloading file...", "URL", fileURL)
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tmpFile, err := saveToTmpFile(data)
	if err != nil {
		logger.Error("downloadFileActivity failed to save tmp file.", "Error", err)
		return "", err
	}
	fileName := tmpFile.Name()
	logger.Info("downloadFileActivity succeed.", "SavedFilePath", fileName)
	return fileName, nil
}

func (a *Activities) ProcessFileActivity(ctx context.Context, fileName string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("processFileActivity started.", "FileName", fileName)

	defer func() { _ = os.Remove(fileName) }() // cleanup temp file

	// read downloaded file
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error("processFileActivity failed to read file.", "FileName", fileName, "Error", err)
		return "", err
	}

	return digest(ctx, data), nil
}

func digest(ctx context.Context, data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func saveToTmpFile(data []byte) (f *os.File, err error) {
	tmpFile, err := ioutil.TempFile("", "temporal_sample")
	if err != nil {
		return nil, err
	}
	_, err = tmpFile.Write(data)
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return nil, err
	}

	return tmpFile, nil
}
