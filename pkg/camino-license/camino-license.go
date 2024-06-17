package caminolicense

import (
	"os"
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
		if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
			wrongFiles = append(wrongFiles, WrongLicenseHeader{f, "File doesn't exist"})
		} else {
			isCustomHeader, headerName, header := checkCustomHeader(f, headersConfig)
			if isCustomHeader {
				correctLicense, reason := verifyCustomLicenseHeader(f, headerName, header)
				if !correctLicense {
					wrongFiles = append(wrongFiles, WrongLicenseHeader{f, reason})
				}
			} else {
				correctLicense, reason := verifyDefaultLicenseHeader(f, headersConfig.PossibleHeaders)
				if !correctLicense {
					wrongFiles = append(wrongFiles, WrongLicenseHeader{f, reason})
				}
			}
		}
	}

	if len(wrongFiles) > 0 {
		return wrongFiles, errors.New("Some files has wrong License Header")
	}

	return wrongFiles, nil

}

func UpdateLicense(files []string, headersConfig HeadersConfig) error {

	var wrongFiles []WrongLicenseHeader
	for _, f := range files {
		if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
			wrongFiles = append(wrongFiles, WrongLicenseHeader{f, "File doesn't exist"})
		} else {
			isCustomHeader, headerName, header := checkCustomHeader(f, headersConfig)
			if isCustomHeader {
				correctLicense, reason := verifyUpdateCustomLicenseHeader(f, headerName, header)
				if !correctLicense {
					wrongFiles = append(wrongFiles, WrongLicenseHeader{f, reason})
				}
			} else {
				correctLicense, reason := verifyUpdateDefaultLicenseHeader(f, headersConfig.PossibleHeaders)
				if !correctLicense {
					wrongFiles = append(wrongFiles, WrongLicenseHeader{f, reason})
				}
			}
		}
	}

	if len(wrongFiles) > 0 {
		return errors.New("Some files has wrong License Header. Please run check command first, solve the issues then run the update command again")
	}

	return nil

}

func checkCustomHeader(file string, headersConfig HeadersConfig) (bool, string, string) {

	// check Custome Headers
	headerName := ""
	longestPath := ""
	header := ""
	for _, customHeader := range headersConfig.CustomHeaders {

		for _, path := range customHeader.Paths {
			pathFiles, err := filepathx.Glob(path)
			if err != nil {
				panic(err)
			}

			file = strings.Replace(file, "./", "", 1)

			if exists(file, pathFiles) {
				if len(longestPath) < len(path) {
					longestPath = path
					headerName = customHeader.Name
					header = customHeader.Header
				}
			}
		}
	}

	if len(headerName) == 0 {
		return false, "", ""
	}

	return true, headerName, header

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
	os.ReadFile(file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, "Cannot read file"
	}
	content := string(bytes)
	currentYear := time.Now().Format("2006")

	header = strings.ReplaceAll(header, "{YEAR}", currentYear)

	if strings.HasPrefix(content, header) {
		return true, ""
	}
	return false, "File doesn't have the same License Header as Custom Header: " + headerName
}

func verifyDefaultLicenseHeader(file string, defaultHeaders []PossibleHeader) (bool, string) {

	os.ReadFile(file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, "Cannot read file"
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

	return false, "File doesn't have the same License Header as any of the possible headers defined in the configuration file"

}

func verifyUpdateCustomLicenseHeader(file string, headerName string, header string) (bool, string) {
	os.ReadFile(file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, "Cannot read file"
	}
	content := string(bytes)
	currentYear := time.Now().Format("2006")
	lastYear := time.Now().AddDate(-1, 0, 0).Format("2006")

	currentYearheader := strings.ReplaceAll(header, "{YEAR}", currentYear)
	lastYearheader := strings.ReplaceAll(header, "{YEAR}", lastYear)

	if strings.HasPrefix(content, currentYearheader) {
		return true, ""
	} else if strings.HasPrefix(content, lastYearheader) {
		newContent := strings.Replace(content, lastYearheader, currentYearheader, 1)
		err := os.WriteFile(file, []byte(newContent), 0666)
		if err != nil {
			return false, "error changing year"
		}

		return true, ""

	}
	return false, "File doesn't have the same License Header as Custom Header: " + headerName
}

func verifyUpdateDefaultLicenseHeader(file string, defaultHeaders []PossibleHeader) (bool, string) {

	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, "Cannot read file"
	}
	content := string(bytes)

	for _, defaultHeader := range defaultHeaders {
		header := defaultHeader.Header
		currentYear := time.Now().Format("2006")
		lastYear := time.Now().AddDate(-1, 0, 0).Format("2006")
		currentYearheader := strings.ReplaceAll(header, "{YEAR}", currentYear)
		lastYearheader := strings.ReplaceAll(header, "{YEAR}", lastYear)

		if strings.HasPrefix(content, currentYearheader) {
			return true, ""
		} else if strings.HasPrefix(content, lastYearheader) {
			newContent := strings.Replace(content, lastYearheader, currentYearheader, 1)
			err := os.WriteFile(file, []byte(newContent), 0666)
			if err != nil {
				return false, "error changing year"
			}

			return true, ""

		}

	}

	return false, "File doesn't have the same License Header as any of the possible headers defined in the configuration file"

}
