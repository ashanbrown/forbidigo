SHELL=bash

test:
	cd examples && diff <(sed 's|CURDIR|$(CURDIR)|' expected_results.txt) <(go run .. 2>&1 | sed '/^go: downloading/d')
	diff <(sed 's|CURDIR|$(CURDIR)|' examples/expected_results_singlechecker.txt) <(go run ./cmd/forbidigo ./examples 2>&1)
