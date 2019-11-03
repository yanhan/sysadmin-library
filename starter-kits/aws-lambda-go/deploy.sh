#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

main() {
	if [ "${#}" -ne 1 ]; then
		printf >&2 "Usage: %s <function_name>\n"  "${0}"
		exit 1
	fi

	declare -r zipfile_name=function.zip
	rm -f "${zipfile_name}"

	CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o main
	zip "${zipfile_name}" ./main
	aws lambda update-function-code \
		--function-name "${1}" \
		--zip-file fileb://"${zipfile_name}" \
		--publish
}

main "$@"
# vim:noet
