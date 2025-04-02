down:
	@echo "stopping docker images ..."
	docker-compose down
	@echo "... done!"

up_build:
	@echo "stopping docker images (if running) ..."
	docker-compose down
	@echo "building (if required) and starting docker images ..."
	docker-compose up --build -d
	@echo "... done!"