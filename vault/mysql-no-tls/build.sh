#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

main() {
	eval "$(minikube docker-env)"
	docker build -t yanhan/vault-mysql-no-tls:0.1 .
}

main "$@"
# vim:noet
