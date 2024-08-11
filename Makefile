# LOCAL DEVELOPMENT UTILITY SHELL APPS
migrate-up:
	sql-migrate up -config=local.dbconfig.yml -env="local-aztebot"

migrate-up-dry:
	sql-migrate up -config=local.dbconfig.yml -env="local-aztebot" -dryrun

migrate-rollback:
	sql-migrate down -config=local.dbconfig.yml -env="local-aztebot"
	
up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

update-env:
	openssl base64 -A -in cmd/radio-service/.prod.env -out cmd/radio-service/.prod.env.out