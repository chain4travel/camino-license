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
)

var headersConfig = caminolicense.HeadersConfig{
	[]caminolicense.DefaultHeader{
		{
			Name:   "l1",
			Header: "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L1\n",
		},
		{
			Name:   "l2",
			Header: "// Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.\n// L2\n",
		},
	}, []caminolicense.CustomHeader{
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
	wrongFiles, err := caminolicense.CheckLicense([]string{"./test_correct_default_1.go"}, headersConfig)
	require.NoError(t, err)
	require.Equal(t, 0, len(wrongFiles))
	require.NoError(t, os.Remove("./test_correct_default_1.go"))
	require.NoError(t, os.WriteFile("./test_correct_default_2.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L2\n\n package caminolicense", time.Now().Year())), 0o600))
	wrongFiles2, err2 := caminolicense.CheckLicense([]string{"./test_correct_default_2.go"}, headersConfig)
	require.NoError(t, err2)
	require.Equal(t, 0, len(wrongFiles2))
	require.NoError(t, os.Remove("./test_correct_default_2.go"))
}

func TestWrongDefaultLicense(t *testing.T) {
	require.NoError(t, os.WriteFile("./test_wrong_default.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// Wrong License\n\n package caminolicense", time.Now().Year())), 0o600))
	wrongFiles, err := caminolicense.CheckLicense([]string{"./test_wrong_default.go"}, headersConfig)
	require.ErrorIs(t, err, caminolicense.CheckErr)
	require.Equal(t, 1, len(wrongFiles))
	require.NoError(t, os.Remove("./test_wrong_default.go"))
}

func TestCorrectCustomLicense(t *testing.T) {
	require.NoError(t, os.WriteFile("./camino_test_correct_custom.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L3\n\n package caminolicense", time.Now().Year())), 0o600))
	require.NoError(t, os.WriteFile("./camino_test_exclude.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L1\n\n package caminolicense", time.Now().Year())), 0o600))
	headersConfig2, _ := caminolicense.GetHeadersConfig("./config_test.yaml")
	wrongFiles, err := caminolicense.CheckLicense([]string{"./camino_test_correct_custom.go", "./camino_test_exclude.go"}, headersConfig2)
	require.NoError(t, err)
	require.Equal(t, 0, len(wrongFiles))
	require.NoError(t, os.Remove("./camino_test_correct_custom.go"))
	require.NoError(t, os.Remove("./camino_test_exclude.go"))
}

func TestWrongCustomLicense(t *testing.T) {
	require.NoError(t, os.WriteFile("./camino_test_exclude.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L3\n\n package caminolicense", time.Now().Year())), 0o600))
	require.NoError(t, os.WriteFile("./camino_test_wrong_custom.go", []byte(fmt.Sprintf("// Copyright (C) 2022-%d, Chain4Travel AG. All rights reserved.\n// L1\n\n package caminolicense", time.Now().Year())), 0o600))
	headersConfig2, _ := caminolicense.GetHeadersConfig("./config_test.yaml")
	wrongFiles, err := caminolicense.CheckLicense([]string{"./camino_test_wrong_custom.go", "./camino_test_exclude.go"}, headersConfig2)
	require.ErrorIs(t, err, caminolicense.CheckErr)
	require.Equal(t, 2, len(wrongFiles))
	require.NoError(t, os.Remove("./camino_test_wrong_custom.go"))
	require.NoError(t, os.Remove("./camino_test_exclude.go"))
}
