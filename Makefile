SHELL=bash

setup:
	pre-commit install

test:
	cd examples && diff <(sed 's|CURDIR|$(CURDIR)|' expected_results.txt) <(go run .. 2>&1 | sed '/^go: downloading/d')
	cd examples && diff <(sed 's|CURDIR|$(CURDIR)|' expected_analyze_types_results.txt) <(go run .. -analyze_types '{p: "^sql\\.DB\\.Exec$$"}' -- 2>&1 | sed '/^go: downloading/d')

lint:
	pre-commit run --all-files

.PHONY: lint test
