# Chirpy

Twitter server clone for study purposes. None of the features are intended to be used in a production environment.

## Features

- Chirpy is a RESTful API that handles simple CRUD operations on chirps (equivalent of Tweets)
- Uses JWT to authenticate users
- Provides a single Webhook, used to grant "payed" features to users after the subscribe to the platform.

## Limitations

- This project was made to explore the http library provided by Golang. Not meant to be used in a production environment.
- It doesn't use a proper DataBase. It relies on a single json file for storage

## Requirements

- [Go](https://go.dev/doc/install) 1.22

## How to install

```go
go install github.com/RealMotz/chirpy
```
Inside a go module:
```go
go get github.com/RealMotz/chirpy
```
