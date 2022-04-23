# Backend Intern Assignment : Retail Pulse

## Description

The task of receiving and scheduling the jobs to process the image has been achieved using the powerful goroutines. The entry point for our program is `main.go`. The APIs have been implemented using `net/http` module (present in `routes.go`), and `gorm` has been used as the choice for object relational mapper for the `postgres` database. On hitting  the submit job api, a `goroutine` is spawned (present in `job.go`) which iterates over the image urls and perform the necessary task of downloading and calculating the perimeter of the images. Database contains two tables, one for `Job` which contains `job_id` and its `status`. Second table is `Image`, which contains metadata of the images collected from store like `url`, `store_id`, `perimeter`, as well as `error_message` in case of any encountered error, e.g. network timeout during downloading the image, invalid image url and so on.
All `Image` objects belonging to a particular `Job` are linked to it using a `foreign_key` relation. Finally, error handling has been maintained everywhere as extensive as possible (present in `errors.go`).

## Installation and Setup

Make sure to have `go` installed in your system. 
Perform the follwoing operations to setup the program:

```
cd $GOPATH
mkdir src/github.com/pranjaldube
cd src/github.com/pranjaldube
git clone <repo-url>
cd <repo>
go mod download
```

To run the program without docker, you would need to setup the postgres database on your local system with the credentials present in `.env`, then simply run

```bash
go run .  
## or
go build
./<binary>
```

To run on docker, make sure to have `docker` and `docker-compose` installed in your system. Then run the following command:

```bash
docker-compose build compose up -d --build
```

To bring the containers down, run

```bash
docker-compose down -v
```

## Work Environment

OS : WSL2 Ubuntu on Windows 10
Text Editor : VS Code with Go extension

## Scope of improvement

`goroutines`, no doubt, is a lightweight yet powerful wrapper for threading processes in `go`. But it is not fully reliable. For something more scalable and managed, I would prefer to opt for some message broker and task queue like `redis`, `rabbitmq`, in which tasks get scheduled and are processed in a queue manner. The latter approach is more robust and tolerant towards failures. But considering the sanctity of time, it was not implemented in the current stage.