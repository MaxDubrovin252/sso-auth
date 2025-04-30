run:
	go run cmd/main.go
docker:
	sudo docker run --name=auth -e POSTGRES_PASSWORD='qwerty' -p 5434:5432 --rm -d postgres


	
 down:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5434/postgres?sslmode=disable' down

 up:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5434/postgres?sslmode=disable' up