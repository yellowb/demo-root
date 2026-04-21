.PHONY: setup dev lint test reset-db

setup:
	./scripts/setup.sh

dev:
	./scripts/dev.sh

lint:
	./scripts/lint.sh

test:
	./scripts/test.sh

reset-db:
	./scripts/reset-db.sh
