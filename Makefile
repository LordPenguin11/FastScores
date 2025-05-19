migrate:
	docker cp internal/db/migrations.sql $$(docker-compose ps -q db):/migrations.sql
	docker-compose exec db psql -U postgres -d livescore -f /migrations.sql 
