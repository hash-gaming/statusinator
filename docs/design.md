# Statusinator design

## Purpose

The purpose of this project is to automate the changing of profile pictures and names in a Slack community once a user changes their status emoji. This can be used to more visibly alert other users of a status change; for example, switching from working to vacation. Additionally, users gain the ability to have preset name/picture pairs that they can easily toggle by just switching their emoji, giving them easier personalization options.

## User stories

Given the purpose statement above, there are some user stories that we can derive.

- As a user, I want my picture and name changed automatically so that I can alert other members that my availability has changed
- As a user, I want to configure name and picture pairs that go with certain status emoji so that it is easier for me to customize my profile

## High level design

The high level design for this project comprises of 3 parts - the event listener server that checks if a onboarded user changed their status, the utilities to handle storing and retrieval of images and names, and the Slack API client that handles changes to users' pictures and names. Technically, the code will be written in Golang mostly because Golang has a well-built standard library and running it under Docker doesn't require a runtime. The assets (emoji and the corresponding names and pictures) will be stored in AWS Simple Storage Service (S3).
