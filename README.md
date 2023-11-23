# go-project-startup

A framework for quickly starting a Go project


## Use the framework to build your own project
1. Update the project information in the makefile
```makefile
APP_NAME = go-project-startup
APP_DESC = go project startup template
BASE_PKG = github.com/zunkk
APP_PKG = $(BASE_PKG)/$(APP_NAME)
```

2. Run update project information cmd
```shell
# This command will update the project name and go package
make reset-project-info
```

## Quick run
```shell
# This command will compile the program and deploy the binary copy to the deployment package
make dev-package

# Start in terminal
./deploy/app start

# Background startup
./deploy/start.sh
```

## Generate deploy package
```shell
# This command will compile the program and package the script and binary into a compressed package
make package

```
