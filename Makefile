SHELL=bash

setup:
	pre-commit install

test:
	cd examples && diff <(sed 's|CURDIR|$(CURDIR)|' expected_results.txt) <(go run .. 2>&1 | sed '/^go: downloading/d')

lint:
	pre-commit run --all-files

.PHONY: lint test
