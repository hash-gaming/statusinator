# statusinator

an automated way to change name and profile picture based off of slack statuses.

## Installation

Once you have the repository pulled down, a simple `go get -t ./...` will install all the dependencies. Although, the dependency installation will also happen on `go build`.

## Initialization

### Automated way using AWS CloudFormation

To start running this application, there needs to be some resource creation for the authentication pieces of this application. There is a AWS CloudFormation stack included in the auth folder that can be used to automatically set up the right pieces. 

### Manual Setup

There are 4 resources that need to be created here. A role for the developers to assume during development, a user that the service can use to run in production, a policy that is shared between the two and an access key pair for the service user. 

## Running

The first thing you're going to have to do is create a `.env` file that has the following properties

```
# Required
ENV={development|production}
BUCKET_NAME=<S3_BUCKET_NAME>
REGION=<VALID_AWS_REGION>
SERVICE_ROLE_ARN=<ARN_OF_THE_ROLE_TO_ASSUME>

# Optional
AWS_ACCESS_KEY_ID=<ACCESS_KEY_ID>
AWS_SECRET_ACCESS_KEY=<SECRET_ACCESS_KEY>
```

The `SERVICE_ROLE_ARN` is the role that will give this application the permissions to read/write to the bucket denoted by `BUCKET_NAME`. The `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are used to provide credentials (either in development or production) to the Docker container.

There isn't much to run right now (the code just reads the objects in the bucket specified above) but once you run a `make build`, you can execute the compiled binary by running `./bin/statusinator`.
