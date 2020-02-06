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
