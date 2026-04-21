.PHONY: setup dev lint test validate-specs reset-db

setup:
	./scripts/setup.sh

dev:
	./scripts/dev.sh

lint:
	./scripts/lint.sh

test:
	./scripts/test.sh

validate-specs:
	./scripts/validate-specs.sh

reset-db:
	./scripts/reset-db.sh
