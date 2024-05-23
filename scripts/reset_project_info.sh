#! /bin/bash
set -e

dir=$1
old_pkg=$2
new_pkg=$3
new_app_name=$4

mv ${dir}/cmd/go-project-startup ${dir}/cmd/${new_app_name}

function x_replace() {
  system=$(uname)

  if [ "${system}" = "Linux" ] || [[ "${system}" =~ "MINGW" ]]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

for file in `grep "${old_pkg}" -rl ${dir}`; do
  if [[ "${file}" == *.go ]] || [[ "${file}" == *go.mod ]]; then
    echo ${file}
    x_replace "s|${old_pkg}/cmd/go-project-startup|${new_pkg}/cmd/${new_app_name}|g" ${file}
    x_replace "s|${old_pkg}|${new_pkg}|g" ${file}
  fi
done

x_replace "s|app_name=go-project-startup|app_name=${new_app_name}|g" ${dir}/deploy/tools/control.sh
x_replace "s|go-project-startup|${new_app_name}|g" ${dir}/deploy/bin_proxy