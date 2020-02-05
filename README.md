# statusinator

an automated way to change name and profile picture based off of slack statuses.

## Installation

Once you have the repository pulled down, a simple `go get -t ./...` will install all the dependencies. Although, the dependency installation will also happen on `go build`.

## Running

The first thing you're going to have to do is create a `.env` file that has the following properties

```
BUCKET_NAME=<S3_BUCKET_NAME>
```

There isn't much to run right now (the code just reads the objects in the bucket specified above) but once you run a `make build`, you can execute the compiled binary by running `./bin/statusinator`.
