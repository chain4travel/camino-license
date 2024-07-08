// Copyright (C) 2022-2024, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package caminolicense

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	config "github.com/chain4travel/camino-license/pkg/config"
	"github.com/pkg/errors"
	"github.com/yargevad/filepathx"
)

type WrongLicenseHeader struct {
	File   string
	Reason string
}

type CaminoLicenseHeader struct {
	Config config.HeadersConfig
}

var CheckErr = errors.New("Some files has wrong License Header")

// public function to start checking for license headers in a list of files or directories
func (h CaminoLicenseHeader) CheckLicense(files []string) ([]WrongLicenseHeader, error) {
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
				// TODO: set license exclusions to be configured in the configuration file
				match, matchErr := filepath.Match("mock_*.go", filepath.Base(path))
				if strings.HasSuffix(path, ".pb.go") || matchErr != nil || match {
					continue
				}
				isWrong, wrongFile := h.checkFileLicense(path)
				if isWrong {
					wrongFiles = append(wrongFiles, wrongFile)
				}
			}
		} else {
			isWrong, wrongFile := h.checkFileLicense(f)
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
func (h CaminoLicenseHeader) checkFileLicense(f string) (bool, WrongLicenseHeader) {
	isCustomHeader, headerName, header := h.checkCustomHeader(f)
	if isCustomHeader {
		correctLicense, reason := verifyCustomLicenseHeader(f, headerName, header)
		if !correctLicense {
			return true, WrongLicenseHeader{f, reason}
		}
	} else {
		correctLicense, reason := h.verifyDefaultLicenseHeader(f)
		if !correctLicense {
			return true, WrongLicenseHeader{f, reason}
		}
	}
	return false, WrongLicenseHeader{}
}

// to check if the file is included in a custom header path
func (h CaminoLicenseHeader) checkCustomHeader(file string) (bool, string, string) {
	for _, customHeader := range h.Config.CustomHeaders {
		absFile, fileErr := filepath.Abs(file)
		if fileErr != nil {
			absFile = file
		}
		if slices.Contains(customHeader.AllFiles, absFile) && !slices.Contains(customHeader.ExcludedFiles, absFile) {
			return true, customHeader.Name, customHeader.Header
		}
	}
	return false, "", ""
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
func (h CaminoLicenseHeader) verifyDefaultLicenseHeader(file string) (bool, string) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Sprintf("Cannot read file: %s", err)
	}
	content := string(bytes)

	for _, defaultHeader := range h.Config.DefaultHeaders {
		header := defaultHeader.Header
		currentYear := time.Now().Format("2006")
		header = strings.ReplaceAll(header, "{YEAR}", currentYear)

		if strings.HasPrefix(content, header) {
			return true, ""
		}
	}

	return false, "File doesn't have the same License Header as any of the default headers defined in the configuration file"
}
