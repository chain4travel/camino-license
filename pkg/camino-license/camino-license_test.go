// Copyright (C) 2022-2024, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package caminolicense_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	caminolicense "github.com/chain4travel/camino-license/pkg/camino-license"
	config "github.com/chain4travel/camino-license/pkg/config"
)

var headersConfig = config.HeadersConfig{
	DefaultHeaders: []config.DefaultHeader{
		{
			Name:   "l1",
			Header: "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L1\n",
		},
		{
			Name:   "l2",
			Header: "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L2\n",
		},
	}, CustomHeaders: []config.CustomHeader{
		{
			Name:         "l3",
			Header:       "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L3\n",
			IncludePaths: []string{"./**/camino*.go"},
			ExcludePaths: []string{"./**/camino_test_exclude.go"},
		},
	},
}

func TestCorrectDefaultLicense(t *testing.T) {
	require.NoError(t, os.WriteFile("./test_correct_default_1.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L1\n\n package caminolicense", time.Now().Year())), 0o600))
	h := caminolicense.CaminoLicenseHeader{Config: headersConfig}
	wrongFiles, err := h.CheckLicense([]string{"./test_correct_default_1.go"})
	require.NoError(t, err)
	require.Empty(t, wrongFiles)
	require.NoError(t, os.Remove("./test_correct_default_1.go"))
	require.NoError(t, os.WriteFile("./test_correct_default_2.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L2\n\n package caminolicense", time.Now().Year())), 0o600))
	wrongFiles2, err2 := h.CheckLicense([]string{"./test_correct_default_2.go"})
	require.NoError(t, err2)
	require.Empty(t, wrongFiles2)
	require.NoError(t, os.Remove("./test_correct_default_2.go"))
}

func TestWrongDefaultLicense(t *testing.T) {
	require.NoError(t, os.WriteFile("./test_wrong_default.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// Wrong License\n\n package caminolicense", time.Now().Year())), 0o600))
	h := caminolicense.CaminoLicenseHeader{Config: headersConfig}
	wrongFiles, err := h.CheckLicense([]string{"./test_wrong_default.go"})
	require.ErrorIs(t, err, caminolicense.CheckErr)
	expectedWrongFiles := []caminolicense.WrongLicenseHeader{
		{
			File:   "./test_wrong_default.go",
			Reason: "File doesn't have the same License Header as any of the default headers defined in the configuration file",
		},
	}
	require.Equal(t, expectedWrongFiles, wrongFiles)
	require.NoError(t, os.Remove("./test_wrong_default.go"))
}

func TestCorrectCustomLicense(t *testing.T) {
	require.NoError(t, os.WriteFile("./camino_test_correct_custom.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L3\n\n package caminolicense", time.Now().Year())), 0o600))
	require.NoError(t, os.WriteFile("./camino_test_exclude.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L1\n\n package caminolicense", time.Now().Year())), 0o600))
	headersConfig2, _ := config.GetHeadersConfig("../config_test.yaml")
	h := caminolicense.CaminoLicenseHeader{Config: headersConfig2}
	wrongFiles, err := h.CheckLicense([]string{"./camino_test_correct_custom.go", "./camino_test_exclude.go"})
	require.NoError(t, err)
	require.Empty(t, wrongFiles)
	require.NoError(t, os.Remove("./camino_test_correct_custom.go"))
	require.NoError(t, os.Remove("./camino_test_exclude.go"))
}

func TestWrongCustomLicense(t *testing.T) {
	require.NoError(t, os.WriteFile("./camino_test_exclude.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L3\n\n package caminolicense", time.Now().Year())), 0o600))
	require.NoError(t, os.WriteFile("./camino_test_wrong_custom.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L1\n\n package caminolicense", time.Now().Year())), 0o600))
	headersConfig2, _ := config.GetHeadersConfig("../config_test.yaml")
	h := caminolicense.CaminoLicenseHeader{Config: headersConfig2}
	wrongFiles, err := h.CheckLicense([]string{"./camino_test_wrong_custom.go", "./camino_test_exclude.go"})
	require.ErrorIs(t, err, caminolicense.CheckErr)
	expectedWrongFiles := []caminolicense.WrongLicenseHeader{
		{
			File:   "./camino_test_wrong_custom.go",
			Reason: "File doesn't have the same License Header as Custom Header: l3",
		},
		{
			File:   "./camino_test_exclude.go",
			Reason: "File doesn't have the same License Header as any of the default headers defined in the configuration file",
		},
	}
	require.Equal(t, expectedWrongFiles, wrongFiles)
	require.NoError(t, os.Remove("./camino_test_wrong_custom.go"))
	require.NoError(t, os.Remove("./camino_test_exclude.go"))
}
