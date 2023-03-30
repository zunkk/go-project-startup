#! /bin/bash
set -e

base_dir=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
export ROOT_PATH=${base_dir}

${base_dir}/tools/control.sh restart