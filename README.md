# statusinator

an automated way to change name and profile picture based off of slack statuses.

## Installation

Once you have the repository pulled down, a simple `go get -t ./...` will install all the dependencies. Although, the dependency installation will also happen on `go build`.

## Running

The first thing you're going to have to do is create a `.env` file that has the following properties

```
BUCKET_NAME=<S3_BUCKET_NAME>
REGION=<VALID_AWS_REGION>
SERVICE_ROLE_ARN=<ARN_OF_THE_ROLE_TO_ASSUME>
```

The `SERVICE_ROLE_ARN` is the role that will give this application the permissions to read/write to the bucket denoted by `BUCKET_NAME`. 

There isn't much to run right now (the code just reads the objects in the bucket specified above) but once you run a `make build`, you can execute the compiled binary by running `./bin/statusinator`.
