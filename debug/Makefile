ifneq (,$(wildcard ./.env))
	include .env
	export
endif

sync-sqlite: build
	cloudquery sync sharepoint.yml sqlite.yml

sync-postgresql: build
	CQ_DESTINATION=postgresql cloudquery sync sharepoint.yml postgresql.yml

build:
	cd .. && go build -o bin/cq-source-sharepoint -v

clean:
	rm db.sqlite cloudquery.log