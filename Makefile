dev:
	docker compose -f compose.yml -f compose.development.yml down
	docker compose -f compose.yml -f compose.development.yml up --build -d

live:
	docker compose down && docker compose up --build -d

down:
	docker compose -f compose.yml -f compose.development.yml down
