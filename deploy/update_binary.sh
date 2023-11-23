#! /bin/bash
set -e

base_dir=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
export REPO_PATH=${base_dir}

if [ ! -f $1 ]; then
  echo "error: new binary($1) does not exist"
  exit 1
fi

${base_dir}/tools/control.sh update-binary $1

