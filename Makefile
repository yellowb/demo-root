.PHONY: setup dev test reset-db

setup:
	./scripts/setup.sh

dev:
	./scripts/dev.sh

test:
	./scripts/test.sh

reset-db:
	./scripts/reset-db.sh
