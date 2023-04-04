run:
	go run main.go

docker: 
	docker build -t obiex-finance .

compose:
	docker-compose up