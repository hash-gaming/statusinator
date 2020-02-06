# Statusinator design

## Purpose

The purpose of this project is to automate the changing of profile pictures and names in a Slack community once a user changes their status emoji. This can be used to more visibly alert other users of a status change; for example, switching from working to vacation. Additionally, users gain the ability to have preset name/picture pairs that they can easily toggle by just switching their emoji, giving them easier personalization options.

## User stories

Given the purpose statement above, there are some user stories that we can derive.

- As a user, I want my picture and name changed automatically so that I can alert other members that my availability has changed
- As a user, I want to configure name and picture pairs that go with certain status emoji so that it is easier for me to customize my profile

## High level design

The high level design for this project comprises of 3 parts - the event listener server that checks if a onboarded user changed their status, the utilities to handle storing and retrieval of images and names, and the Slack API client that handles changes to users' pictures and names. Technically, the code will be written in Golang mostly because Golang has a well-built standard library and running it under Docker doesn't require a runtime. The assets (emoji and the corresponding names and pictures) will be stored in AWS Simple Storage Service (S3).

## Detailed design

As mentioned above, there are three considerations for this design

- Event listener server
- Storage helper
- Slack client

Added below are detailed designs for each of the components and their interactions with each other.

### Event listener server

The only purpose of the event listener server will be to check for the `user_change` event coming from the Slack Events API. The documentation for this event as well as the general guidelines for using the events API is linked below. The payload of the event will contain the user that was responsible for this change and we can use the `users.profile.get` API to retrieve the newly changed status emoji.

At startup, the entire JSON file that stores all onboarded users, as it stands, will be read and will be used to hydrate the in-memory cache maintained by the event listener. On request, if a user is not found in the event listener cache, the master JSON will be checked and only reread if it has been updated since the last read. Otherwise, we will do nothing.

Before investing the effort to get the changed profile, we will check to see if the user is onboarded with this project. This will be done by either checking the in-memory cache or asking the storage helper if this user is onboarded.

Once we have the status emoji and the user, we will invoke the storage helper to return to us the name and the picture that corresponds to that emoji status. Once we have a link and a name string, we will pass that information along to the Slack API client which will commit the information to Slack.

### Storage helper

The organization of the datastore is as follows

```
bucket/
  |- users.json
  |- <user_id_1>/
  | |- <uuid>.png
  | |- <uuid>.png
  |- <user_id_2>/
    |- <uuid>.png
    |- <uuid>.png
    |- <uuid>.png
```

Each image stored will be given an UUID and stored in the `users.json` datastore as the UUID. The S3 link for the picture will be determined by the bucket name, the folder name and the UUID of the picture.

The storage helper has 2 functions

1. validate that we have details on the user that was changed
1. retrieve the name and picture pair for an emoji given the user ID

Both these functions will depend on S3 to be their datastore. The event listener will maintain a cache of users and whether they are onboarded or not backed by the storage helper. The S3 bucket will only store users that are onboarded in a flat JSON file.

```json
{
  "lastUpdatedAt": "<unix_epoch_time>",
  "users": {
    "<user_id_1>": {
      "<emoji_colon_code>": {
        "name": "<some_name>",
        "picture": "<an_s3_link>"
      }
    }
  }
}
```

This file will serve as the primary method of storing details about onboarded users and will only be written to by using a single admin level storage helper API called `/onboard`. This API will expect a JSON request body very similar to the structure of the `users` key above. This will serve as the way to onboard users to this system, for now. This API will also update the `lastUpdatedAt` time so that the event listener can detect changes. This covers the first function.

Once the storage helper is invoked in the context of the second function, it will asynchronously go and retrieve the name and picture pair that corresponds to the newly changed status emoji for that specific user. There are two error cases exposed here

1. There was no name and picture pair found for that emoji
1. There was only a name (or only a picture) found for that emoji

These error cases will be discussed in more detail in the next section. The only responsiblity of the storage helper will be to return the data that it has found.

### Slack Client

The only charter of the Slack client is to construct requests and send them to the Slack API. It will mostly deal with the `users.profile.set` and the `users.setPhoto` APIs.

Given that the storage helper returns any data found for a user and emoji status pair, the Slack client will create a new request for changing the user's name as well as downloading the picture and then sending it to the Slack APi via the `users.setPhoto` method.

There is relatively little functionality here and that's by design. Having previous experience working with Slack's APIs before, we have found it useful to have consistent behavior across our codebase even though the Slack API might change. It is also nice to have all of the potentially uncontrollable breaking points contained within a single "module".

## Performance

Since this is a request based system and users don't send requests in a periodic, synchronous manner, we will have to consider queueing, throttling and caching to improve the performance of the system. Each of these is talked about below.

### Request queuing

A request queue will live under the event listener component. For the purposes of this project, each time an onboarded user changes their status emoji, they are basically making a request into the system. This request may lead to no action from the system but it is still worth considering as a request. Therefore, we are going to start with the assumption that at any given time, there can only be 10 standby requests on the server.

Once this capacity of 10 is exhausted, any additional requests that are made will divert to a dead letter queue. The dead letter queue will empty itself into the request queue once the request queue reaches an empty state. The dead letter queue will also have space for 10 requests and we will monitor the size of the dead letter queue to find the right size for the primary request queue. Too many requests in the dead letter queue would mean that the primary request queue needs to be scaled up.

### Throttling

Since the infrastructure supporting this project still needs to be approachable for individuals, we will have to have strict throttling rules. Throttling will be triggered in three cases

1. if a user updates their emoji in quick succession resulting in two requests in the queue
1. the request queue is full
1. both the request queue and the dead letter queue are full

In the first case, if we have already found a request for the user that has made a new request, we will drop the old request and add a new request for the user. This is to ascertain that a singular user is not able to become the noisy neighbor and DoS the service impacting other users.

In the second case, the event listener will return an HTTP 503 to the API once we have determined that it is a request made by an onboarded user. This will cause the Slack Events API to retry the request again immediately, 1 minute later, and 5 minutes later. The first retry will put the request into the dead letter queue and stop the retry loop. It is expected that we won't need the 1 minute and 5 minute retries but we will have to look at the metrics and tune the system. We can detect a retry by reading the `X-Slack-Retry-Num` header in the incoming request.

In the third case, the service will start returning HTTP 503 with the `X-Slack-No-Retry: 1` header so that Slack doesn't retry sending events at all. This would be considered a total failure case for the system.

### Caching

The main cache maintained by the sytem will be the cache of users and their specific emoji maintained by the event listener component. This cache will be hydrated on startup using the data from the the storage helper and will be kept in memory. Rehydrating the cache will happen when there is a cache miss for a particular user and the master JSON stored with the storage helper has been updated since it was last read.

The event listener will be the only reader of this cache and the storage helper will be the only writer for the master JSON that the cache will be based upon.

## Architecture

The architecture of the application will be modeled as microservices but contained within the same application for now. Additional details on scaling this system can be found in the scaling section. We will leverage goroutines to help create separation of concerns in the application as well as cater to the asynchronous and parallelized nature of most of the processing done by this system.

Each component listed above will act as an orchestrator for running and managing goroutines. There will be an input and output queue (modeled by go channels) that will plumb information between the three components, as well as, between the user and the system and between the Slack and the system. As it happens, the request queue and the dead letter queue will also be go channels. There will be separate goroutines for each of the actions below and quite possibly more.

- processing each event
- processing each request from an onboarded user
- getting the name and the picture associated with an emoji
- setting the name of the user through the Slack API
- setting the picture for the user through the Slack API

## Observability

Since each of the components above have been modeled as a microservice, we will have to make sure that the observability of the lifecycle of a request will still be preserved as it would be in a monolithic system. To this end, we will leverage both `stdout` and file based logging to be able to trace requests through the system. Each request coming in will be given a UUID and will have the same UUID until the request lifecycle has ended. This UUID will also be useful as we log the path that the request has taken throughout the system.

In the interest of keeping the infrastructural complexity and cost low, we will most likely run this system on a small host with limited resources. This means that we can't store a lot of logs on the host itself and the logging will need to be moved out to some storage solution long term to enable the system owners to get a historical sense of the data as well. We will use S3 for this and will export logs from the host to S3 periodically. Depending on the data and the frequency of use of this system, it could be as quick as every hour and as slow as every day.

The other aspect of observability are metrics. System metrics are time-series data that we want to capture to make sure that the system is performing to the promised SLA as well as it is right-size scaled to meet the users' demands. To that end, we will capture the following data with each metric having the ability to be disabled so as to not overload the senses.

- onboarded user request lifecycle duration in milliseconds
- Slack event parsing latency in milliseconds
- inbound onboarded user request queue size as count
- dead letter queue size as count
- S3 data fetch latency in milliseconds
- Slack API client latency in milliseconds
- number of onboarded users as count
- the size of the JSON datastore as bytes

These metrics will be stored as an unsorted, append-only CSV file on the host and will follow the same export rules as the log files above. The format of the file is listed below

```csv
metric_name,timestamp,value,units
request_lifecycle_duration,<unix_epoch_timestamp>,500,ms
inbound_request_queue_size,<unix_epoch_timestamp>,3,count
```

Once we have these metrics captured, we will then be able to load this data up any analysis tool that can ingest CSV data and analyze scaling and performance of the system.

One thing to note is that the logs and metrics are being captured in this way to keep the total cost of running the system low. As this system scales out, we will probably switch to pushing logs and metrics to something like AWS CloudWatch or Datadog. This is discussed further in the scaling section.

## Infrastructure

This system requires, at bare minimum, some cloud storage and some policies for access. At full bore, the system can encompass some compute capacity, more storage, some database tables, access and execution roles as well as networking and monitoring stacks. None of this infrastructure will be built manually. Instead we will leverage the AWS Cloud Development Kit (CDK) or Terraform to build these resources for us. These libraries allow users to use programming constructs like loops and conditional statements to make decisions about their infrastructure. We will create another project to hold and deploy this infrastructure to the cloud as well as maintain the infrastructure under version control.

As an additional note, if only AWS resources are in use, it is recommended that we use the AWS CDK because it works directly with AWS CloudFormation to orchestrate the creation and modification of resources. It also allows users to build higher level constructs that comprise of AWS resources coupled together for ease of duplication and deployment.

## Security

There are 2 major areas of security concern with this system

1. access to the contents of the bucket
1. access to the JSON datastore of user IDs and emoji

The access to the contents of the bucket issue will be solved by using AWS Security Token Service (STS) to assume a role that only has access to that bucket specifically every time access to the bucket is needed. This access will then be logged to AWS CloudTrail and be retained for 90 days. This role will only have read permissions to the contents of the bucket.

As an extension of the first case, the JSON datastore will be read in a simliar way since it is also stored within S3. The additional security concerns come from the `/onboard` API. To lock this feature down, only the administrators of the workspace will be allowed to use the `/onboard` API to onboard users. We will use the `admin.users.list` API to determine if a user is an administrator of the workspace and then only allow them to write too the `users.json` datastore. Once the datastore is updated and if and only if the event listener request queue is empty, we will prompt the event listener to rehydrate the user cache. This role will have read and write permissions to the bucket and actions will be stored in AWS CloudTrail for 90 days.

## Scaling considerations

To start, this system is built to be very approachable and cost-effective, the application (the go code) can be run on any Docker host that is connected to the internet and has the right permissions deployed to it and the cloud storage is using S3 which is the most cost effective storage solution available for the cloud.

There are two types of scaling events that the system owners will need to consider - horizontal and upwards.

Scaling upwards is relatively easy to do with this infrastructure because S3 can scale nearly infinitely in the storage that it has given that the system owners are able to pay for it and in the cloud, we can deploy this application to a Docker container running on a host that has access to more resources, up to, frankly, a ridiculous size.

Scaling horizontally usually presents more challenges to systems. This system has specifically been designed in a way to enable system owners to scale it without having to make too many changes. The next few paragraphs consider each part of the system that will need scaling concerns.

The application is already internally modeled as a network of microservices. We can leverage docker networking and split out the singular application into it's components. The input and output queues become websocket connections backed by input and output queues and the singular application becomes 3 separate containers that share a network. This allows each bit of the application to be scaled out independently and also allows each component to manage their own load balancing.

As mentioned above already, S3 can grow to nearly infinite scale so the storage of the user pictures growing is already accounted for. The bottleneck that will impact performance will be the JSON datastore. Luckily, the JSON datastore has been modeled in a way that will ease the migration to a document store like AWS DynamoDB. This will also solve the concurrency problems that we would see with trying to support multiple readers and writers to a singular file-based datastore.

As the application scales out, the metrics and observability story will also need to be scaled to match. To start, the host stores log files and metrics files that are exported to S3 periodically. We can tune the periodicity of the export to support large amounts of data as well as onboard the application to something like AWS CloudWatch Metrics and AWS CloudWatch Logs or Datadog to help with storage, analysis, and review of the metrics and the logs.

In the interest of the complete documentation, as the application scales, each component and quite possibly the entire application will need to be put behind layers of load balancers. This can be a hidden cost that surprises system owners are they are scaling their system.

## Open source considerations

This project and all documentation written for it will be open sourced at a later date. While nascent, Github will be used as the storage and collaboration tool but the project will be marked private until ready for public release.

## Open questions

- How do we detect which properties changed when a `user_change` event comes in?
- Can admins/owners of unpaid Slack communities update another user's profile?

# Appendices

## Appendix A - references

1. [Slack Events API overview](https://api.slack.com/events-api)
1. [Slack Events API - `user_change` event](https://api.slack.com/events/user_change)
1. [Slack API - `users.profile.set`](https://api.slack.com/methods/users.profile.set)
1. [Slack API - `users.setPhoto`](https://api.slack.com/methods/users.setPhoto)
1. [Slack API - `users.profile.get`](https://api.slack.com/methods/users.profile.get)
1. [Slack API - `admin.users.list`](https://api.slack.com/methods/admin.users.list)
1. [AWS Cloud Development Kit getting started](https://docs.aws.amazon.com/cdk/latest/guide/getting_started.html)
