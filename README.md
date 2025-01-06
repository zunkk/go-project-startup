# go-project-startup

A framework for quickly starting a Go project

Table of Contents
=================

* [go-project-startup](#go-project-startup)
    * [Use the framework to build your own project](#use-the-framework-to-build-your-own-project)
    * [Quick run](#quick-run)
    * [Generate deploy package](#generate-deploy-package)
    * [Tools](#tools)
        * [Generate db models code from db](#generate-db-models-code-from-db)

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
./deploy/bin_proxy start

# Background startup
./deploy/start.sh
```

## Generate deploy package

```shell
# This command will compile the program and package the script and binary into a compressed package
make package

```

## Tools

### Generate db models code from db

Use [sqlboiler](https://github.com/volatiletech/sqlboiler) to generate db models code.

1. Create/Update db tables


2. Update db information in `build/sqlboiler.toml`,
   default is `Postgres`

```toml
[psql]
dbname = "test"
host = "127.0.0.1"
port = 5432
user = "zunkk"
pass = "zunkk"
schema = "public"
sslmode = "disable"
```

3. Update generate config

```makefile
MODELS_PATH := ${PROJECT_PATH}/internal/core/model
DB_TYPE = psql
```

4. Generate db models code

```shell
make generate-models
```

