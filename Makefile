# Daftar service yang ada
SERVICES := order product shop user warehouse

# Jalankan mockgen untuk setiap file _service.go dan _repository.go
generate-mocks:
	@echo "Generating mocks for all services..."
	@for service in $(SERVICES); do \
		for file in micro-services/$$service/service/*_service.go; do \
			echo "Generating mock for $$file..."; \
			mockgen -source=$$file -destination=micro-services/$$service/mocks/mock_$$service_$$file -package=mocks; \
		done; \
		for file in micro-services/$$service/repository/*_repository.go; do \
			echo "Generating mock for $$file..."; \
			mockgen -source=$$file -destination=micro-services/$$service/mocks/mock_$$service_$$file -package=mocks; \
		done; \
	done

# Clean mocks: hapus semua folder mocks di setiap service
clean-mocks:
	@echo "Cleaning mocks for all services..."
	@for service in $(SERVICES); do \
		rm -rf micro-services/$$service/mocks; \
	done

.PHONY: generate-mocks clean-mocks