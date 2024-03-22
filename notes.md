### Run a docker images

docker run --name pg-go -v /mydata/volumes/pggo:/var/lib/postgresql/data -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password -e POSTGRES_DB=pggo -d -p 5432:5432 postgres
