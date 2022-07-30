export

RUNNER=docker-compose exec app sql-migrate

ifeq ($(p),host)
 	RUNNER=sql-migrate
endif

MIGRATE=$(RUNNER)

migrate-status:
	$(MIGRATE) status

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

redo:
	@read -p  "Are you sure to reapply the last migration? [y/n]" -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; \
	then \
		$(MIGRATE) redo; \
	fi

create:
	@read -p  "What is the name of migration? " NAME; \
	${MIGRATE} new $$NAME

.PHONY: status migrate-up migrate-down redo create