# Generate RSA keys for JWT signing
gen-keys:
	@echo "Generating RSA private and public keys..."
	openssl genrsa -out svc/auth/private.pem 2048
	openssl rsa -in svc/auth/private.pem -pubout -out svc/auth/public.pem

# Clean up keys
clean-keys:
	rm -f svc/auth/private.pem svc/auth/public.pem

run-auth:
	go run ./svc/auth

run-product:
	go run ./svc/product

run-frontend:
	cd frontend && npm run dev

migrate-identity:
	go run ./cmd/migrate identity

docker-up:
	docker compose up --build

docker-down:
	docker compose down

test-auth:
	go test ./svc/auth/...

fmt:
	go fmt ./...

lint:
	golangci-lint run

dev:
	make -j3 run-auth run-product run-frontend