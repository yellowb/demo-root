.PHONY: dev test reset-db

dev:
	./scripts/dev.sh

test:
	./scripts/test.sh

reset-db:
	./scripts/reset-db.sh
