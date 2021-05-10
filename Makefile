SHELL=bash

test:
	cd examples && diff <(sed 's|CURDIR|$(CURDIR)|' expected_results.txt) <(go run .. 2>&1 | sed '/^go: downloading/d')
