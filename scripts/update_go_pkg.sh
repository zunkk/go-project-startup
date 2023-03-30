#! /bin/bash
set -e

dir=$1
old_pkg=$2
new_pkg=$3

for file in `grep "${old_pkg}" -rl ${dir}`; do
  if [[ "${file}" == *.go ]] || [[ "${file}" == *go.mod ]]; then
    echo ${file}
    sed -i "" "s|${old_pkg}|${new_pkg}|g" ${file}
  fi
done

