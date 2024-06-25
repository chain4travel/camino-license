// Copyright (C) 2022-2024, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

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

var CheckErr = errors.New("Some files has wrong License Header")

// public function to start checking for license headers in a list of files or directories
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
		return wrongFiles, CheckErr
	}

	return wrongFiles, nil
}

// To check if a file should have a custom license header or one of the default ones
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

// to check if the file is included in a custom header path
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

// Checks if a file exists in a list of files
func exists(filename string, files []string) bool {
	for _, f := range files {
		if f == filename {
			return true
		}
	}
	return false
}

// to verify if a custom license header from the configuration is similar to the one in the file.
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

// to verify if any of the default license headers from the configuration is similar to the one in the file.
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
