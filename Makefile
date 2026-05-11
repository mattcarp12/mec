# Generate RSA keys for JWT signing
gen-keys:
	@echo "Generating RSA private and public keys..."
	openssl genrsa -out backend/svc/auth/private.pem 2048
	openssl rsa -in backend/svc/auth/private.pem -pubout -out backend/svc/auth/public.pem

# Clean up keys
clean-keys:
	rm -f backend/svc/auth/private.pem backend/svc/auth/public.pem
