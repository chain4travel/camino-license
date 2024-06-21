package caminolicense

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yargevad/filepathx"
)

type WrongLicenseHeader struct {
	File   string
	Reason string
}

func CheckLicense(files []string, headersConfig HeadersConfig) ([]WrongLicenseHeader, error) {
	var wrongFiles []WrongLicenseHeader
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			wrongFiles = append(wrongFiles, WrongLicenseHeader{f, "File doesn't exist"})
			continue
		}

		if info.IsDir() {
			pathFiles, filePathErr := filepathx.Glob(f + "/**/*.go")
			if filePathErr != nil {
				wrongFiles = append(wrongFiles, WrongLicenseHeader{f, "Cannot find .go files under this directory"})
				continue
			}
			for _, path := range pathFiles {
				match, matchErr := filepath.Match("mock_*.go", filepath.Base(path))
				if strings.HasSuffix(path, ".pb.go") || matchErr != nil || match {
					continue
				}
				isWrong, wrongFile := checkFileLicense(path, headersConfig)
				if isWrong {
					wrongFiles = append(wrongFiles, wrongFile)
				}
			}
		} else {
			isWrong, wrongFile := checkFileLicense(f, headersConfig)
			if isWrong {
				wrongFiles = append(wrongFiles, wrongFile)
			}
		}
	}

	if len(wrongFiles) > 0 {
		return wrongFiles, errors.New("Some files has wrong License Header")
	}

	return wrongFiles, nil
}

func checkCustomHeader(file string, headersConfig HeadersConfig) (bool, string, string) {
	for _, customHeader := range headersConfig.CustomHeaders {
		absFile, fileErr := filepath.Abs(file)
		if fileErr != nil {
			absFile = file
		}
		if exists(absFile, customHeader.AllFiles) && !exists(absFile, customHeader.ExcludedFiles) {
			return true, customHeader.Name, customHeader.Header
		}
	}
	return false, "", ""
}

func exists(filename string, files []string) bool {
	for _, f := range files {
		if f == filename {
			return true
		}
	}
	return false
}

func verifyCustomLicenseHeader(file string, headerName string, header string) (bool, string) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Sprintf("Cannot read file: %s", err)
	}
	content := string(bytes)
	currentYear := time.Now().Format("2006")

	header = strings.ReplaceAll(header, "{YEAR}", currentYear)

	if strings.HasPrefix(content, header) {
		return true, ""
	}
	return false, "File doesn't have the same License Header as Custom Header: " + headerName
}

func verifyDefaultLicenseHeader(file string, defaultHeaders []DefaultHeader) (bool, string) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Sprintf("Cannot read file: %s", err)
	}
	content := string(bytes)

	for _, defaultHeader := range defaultHeaders {
		header := defaultHeader.Header
		currentYear := time.Now().Format("2006")
		header = strings.ReplaceAll(header, "{YEAR}", currentYear)

		if strings.HasPrefix(content, header) {
			return true, ""
		}
	}

	return false, "File doesn't have the same License Header as any of the default headers defined in the configuration file"
}

func checkFileLicense(f string, headersConfig HeadersConfig) (bool, WrongLicenseHeader) {
	isCustomHeader, headerName, header := checkCustomHeader(f, headersConfig)
	if isCustomHeader {
		correctLicense, reason := verifyCustomLicenseHeader(f, headerName, header)
		if !correctLicense {
			return true, WrongLicenseHeader{f, reason}
		}
	} else {
		correctLicense, reason := verifyDefaultLicenseHeader(f, headersConfig.DefaultHeaders)
		if !correctLicense {
			return true, WrongLicenseHeader{f, reason}
		}
	}
	return false, WrongLicenseHeader{}
}
