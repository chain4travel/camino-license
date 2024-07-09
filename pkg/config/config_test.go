// Copyright (C) 2022-2024, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	headersConfig, err := GetHeadersConfig("../config_test.yaml")
	require.NoError(t, err)
	expectedHeadersConfig := HeadersConfig{
		[]DefaultHeader{
			{
				Name:   "l1",
				Header: "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L1\n",
			},
			{
				Name:   "l2",
				Header: "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L2\n",
			},
		},
		[]CustomHeader{
			{
				Name:         "l3",
				Header:       "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L3\n",
				IncludePaths: []string{"./**/camino*.go"},
				ExcludePaths: []string{"./**/camino_*exclude.go"},
			},
		},
	}
	require.Equal(t, expectedHeadersConfig.DefaultHeaders, headersConfig.DefaultHeaders)
	require.Equal(t, expectedHeadersConfig.CustomHeaders[0].Name, headersConfig.CustomHeaders[0].Name)
	require.Equal(t, expectedHeadersConfig.CustomHeaders[0].Header, headersConfig.CustomHeaders[0].Header)
	require.Equal(t, expectedHeadersConfig.CustomHeaders[0].IncludePaths, headersConfig.CustomHeaders[0].IncludePaths)
	require.Equal(t, expectedHeadersConfig.CustomHeaders[0].ExcludePaths, headersConfig.CustomHeaders[0].ExcludePaths)
}

func TestNoConfig(t *testing.T) {
	_, err := GetHeadersConfig("../config2_test.yaml")
	require.ErrorIs(t, err, os.ErrNotExist)
}
