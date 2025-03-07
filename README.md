# edgex-snap-hooks
[![Go Reference](https://pkg.go.dev/badge/github.com/canonical/edgex-snap-hooks.svg)](https://pkg.go.dev/github.com/canonical/edgex-snap-hooks/v2)

Snap hooks library used by [EdgeX Foundry](https://docs.edgexfoundry.org/) Go service snaps.  
It provides utilites to implement snap hooks, including some wrappers for the [`snapctl`](https://snapcraft.io/docs/using-snapctl) commands.

### Usage
Download or upgrade to the latest version:
```
go get github.com/canonical/edgex-snap-hooks/v2
```
Please refer to [go get docs](https://pkg.go.dev/cmd/go#hdr-Add_dependencies_to_current_module_and_install_them) for details.

#### Example

```go
package main

import (
	"fmt"
	"os"

	hooks "github.com/canonical/edgex-snap-hooks/v2"
)

func main() {
	var err error

	if err = hooks.Init(false, "edgex-device-example"); err != nil {
		fmt.Printf("initialization failure: %s", err)
		os.Exit(1)
	}

	// copy file from $SNAP to $SNAP_DATA
	if err = hooks.CopyFile(hooks.Snap+"/config.json", hooks.SnapData+"config.json"); err != nil {
		hooks.Error(err.Error())
		os.Exit(1)
	}
  
	// read env var override configuration
	cli := hooks.NewSnapCtl()
	envJSON, err := cli.Config(hooks.EnvConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'env' failed: %v", err))
		os.Exit(1)
	}
	hooks.Debug(fmt.Sprintf("envJSON: %s", envJSON))
}

```

### Testing (WIP)
The tests need to run in a snap environment:

```bash
snapcraft build
```
