// Copyright (C) 2022-2024, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/yargevad/filepathx"
	"gopkg.in/yaml.v2"
)

type DefaultHeader struct {
	Name   string `yaml:"name"`
	Header string `yaml:"header"`
}

type CustomHeader struct {
	Name          string   `yaml:"name"`
	Header        string   `yaml:"header"`
	IncludePaths  []string `yaml:"include-paths"`
	ExcludePaths  []string `yaml:"exclude-paths"`
	AllFiles      []string
	ExcludedFiles []string
}

type HeadersConfig struct {
	DefaultHeaders       []DefaultHeader `yaml:"default-headers"`
	CustomHeaders        []CustomHeader  `yaml:"custom-headers"`
	ExcludedHeadersPaths []string        `yaml:"headers-excluded-paths"`
	ExcludedFiles        []string
}

// read configuration file
func GetHeadersConfig(configPath string) (HeadersConfig, error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return HeadersConfig{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}
	headersConfig := &HeadersConfig{}
	err = yaml.Unmarshal(yamlFile, headersConfig)
	if err != nil {
		return HeadersConfig{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}
	configAbsPath, err := filepath.Abs(configPath)
	if err != nil {
		fmt.Println("Error: Couldn't get the absolute path for the config file:", configPath)
		configAbsPath = configPath
	}

	for i, customHeader := range headersConfig.CustomHeaders {
		includedFiles, err := getCustomHeaderIncludedFiles(customHeader, filepath.Dir(configAbsPath))
		if err != nil {
			return HeadersConfig{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}
		headersConfig.CustomHeaders[i].AllFiles = includedFiles

		excludedFiles, err := getCustomHeaderExcludedFiles(customHeader, filepath.Dir(configAbsPath))
		if err != nil {
			return HeadersConfig{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}
		headersConfig.CustomHeaders[i].ExcludedFiles = excludedFiles
	}
	excludedFiles, err := getExcludedFiles(headersConfig.ExcludedHeadersPaths, filepath.Dir(configAbsPath))
	if err != nil {
		return HeadersConfig{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}
	headersConfig.ExcludedFiles = excludedFiles
	return *headersConfig, nil
}

// walk through directories of include-paths to get all possible files that matches the pattern
func getCustomHeaderIncludedFiles(customHeader CustomHeader, dir string) ([]string, error) {
	var files []string
	for _, path := range customHeader.IncludePaths {
		absPath := path
		if !filepath.IsAbs(path) {
			absPath = filepath.Join(dir, path)
		}
		pathFiles, err := filepathx.Glob(absPath)
		if err != nil {
			return files, errors.New("Cannot get file matches of the custom header included path: " + path)
		}
		files = append(files, pathFiles...)
	}
	return files, nil
}

// walk through directories of exclude-paths to get all possible files that matches the pattern
func getCustomHeaderExcludedFiles(customHeader CustomHeader, dir string) ([]string, error) {
	var files []string
	for _, path := range customHeader.ExcludePaths {
		absPath := path
		if !filepath.IsAbs(path) {
			absPath = filepath.Join(dir, path)
		}
		pathFiles, err := filepathx.Glob(absPath)
		if err != nil {
			return files, errors.New("Cannot get file matches of the custom header excluded path: " + path)
		}
		files = append(files, pathFiles...)
	}
	return files, nil
}

// walk through directories of headers-excluded-paths to get all possible files that matches the pattern
func getExcludedFiles(paths []string, dir string) ([]string, error) {
	var files []string
	for _, path := range paths {
		absPath := path
		if !filepath.IsAbs(path) {
			absPath = filepath.Join(dir, path)
		}
		pathFiles, err := filepathx.Glob(absPath)
		if err != nil {
			return files, errors.New("Cannot get file matches of the custom header excluded path: " + path)
		}
		files = append(files, pathFiles...)
	}
	return files, nil
}
