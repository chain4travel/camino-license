# camino-license
A go package to check license headers.


# Usage
`camino-license check --config=config.yaml FILES/DIRS` 
It checks license headers in the specified files or directories according to the given configuration file. `FILES/DIRS` are space-separated strings of the path (absolute or relative). 
Example:
`camino-license check --config=config.yaml camino-license/cmd camino-license/pkg/camino-license/config.go ` 


# Configuration
This is an example of a configuration file:

```yml
default-headers: 
  - name: avax
    header: |
      // Copyright (C) 2019-{YEAR}, Ava Labs, Inc. All rights reserved.
      // See the file LICENSE for licensing terms.

  - name: avax-c4t
    header: |
      // Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.
      //
      // This file is a derived work, based on ava-labs code whose
      // original notices appear below.
      //
      // It is distributed under the same license conditions as the
      // original code from which it is derived.
      //
      // Much love to the original authors for their work.
      // **********************************************************
      // Copyright (C) 2019-{YEAR}, Ava Labs, Inc. All rights reserved.
      // See the file LICENSE for licensing terms.

custom-headers:
  - name: c4t
    header: |
      // Copyright (C) 2022-{YEAR}, Chain4Travel AG. All rights reserved.
      // See the file LICENSE for licensing terms.
    
    include-paths:
      - "./**/camino*.go"
      - "./test"

    exclude-paths:
      - "./**/camino_visitor2.go"

```
`default-headers`: If the file is not specified in the custom-headers paths, then it should contain one of the default headers

`name`: Name to the header object. It must be unique. **(Required)**

`header`: License header to be used. **(Required)**

`custom-headers`: Headers to be used in certain files according to the `include-paths` and `exclude-paths`.

`include-paths`: list of path pattern to identify files which should contain that custom license header. Use `**`to include all sub directories. Use `*` to include all matching pattern in the same directory. **(Required if custom-headers is specified)**

`exclude-paths`: list of path pattern to identify files which will be excluded from `include-paths`. Use `**`to include all sub directories. Use `*` to include all matching pattern in the same directory.

`{YEAR}`: It will be automatically replaced with current year when the check runs.
