app:
	docker compose up -d --build
	docker run --rm -v "$(shell pwd)/migrations:/migrations" \
    		--network host \
    		migrate/migrate \
    		-path=/migrations \
    		-database "postgresql://root:pass@localhost:5432/news?sslmode=disable" \
    		up
