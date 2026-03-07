.PHONY: db-backup

db-backup:
	@echo "Dumping current PostgreSQL database..."
	@docker exec thyngo-postgres-1 pg_dump -U thyngo_user -d thyngo_db > backups/$$(date +%Y%m%d_%H%M%S)_backup.sql
	@echo "Database successfully backed up!"
