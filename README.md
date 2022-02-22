# json-validator (Go)

## Summary

This is a simple REST http application written in Go. It stores json
schemas on the filesystem and validates user provided documents against
these schemas.

**NOTE:** It is only a prototype, don't use it in production!

## Install and run

First you need to clone the repository to get the application code.

```
$ git clone https://github.com/janospapp/json-validator.git
$ cd json-validator
```

Once you have the code you can either run the application directly
on your machine, if you have go 1.17+ installed, or you can run it
inside a Docker container.

### Running it from a local Go environment

To try out the application you can simply run it from the repository
without installing it.

```
$ go run .
```

It will start the application on **localhost:8000**.

If you want to use the application on a regular basis and have it
installed then run

```
$ go install
```

Which installs it to your **GOBIN** folder. If its on your $PATH
you can simply run the application from any directory.

```
$ json-validator
```

**Note:** It's important which directory are you in when you start
the application, because it saves the uploaded schemas under
**./saved_schemas/** directory. You need to ensure that the current
directory is writable by the application.

### Running it from Docker

You can use the provided Dockerfile to build and run the application.

```
$ docker build -t json-validator .
$ docker run --rm -p 8000:8000 --name sample-validator json-validator
```

The application will be started in a new container named **sample-validator**
(or whatever name you provided).

To use a different local port, change the **-p** option accordingly. E.g
to access the application from **localhost:3001** use

```
-p 3001:8000
```

**Note:** The schemas will be saved inside the container, so when you
delete the container, all the uploaded schemas will be lost. To preserve
them, attach a volume to the container by setting the **-v** parameter
to docker run

```
-v /local/path/to/save/schemas:/usr/src/app/saved_schemas
```

## API endpoints

```
POST /schema/:id
```

The **id** can't be empty, and must be unique. The request fails if there
is already a schema with the given id.

The application checks if the provided data is a valid JSON document.

If everything is correct, then it saves the schema under
**./saved_schemas/id.json**

```
GET /schema/:id
```

It returns the saved schema if the id is valid and the schema exists. Otherwise
it returns a corresponding error.

```
POST /validate/:id
```

It validates the request data against the given schema. The request fails if
the **id** is invalid, there is no such schema or the request data is not a
valid JSON document.

Otherwise it succeeds and the response data will tell if the document complies
to the schema. If not, the first error message is also included in the
response.

## Tests

The main endpoint handler functions are tested by the application. You can run
the tests from the repository root by

```
$ go test ./...
```

or run

```
$ go test
```

from either **schema** or **validator** folders to test only that package.

## Improvement bits

The application is far from perfect or ready. A few ideas what to improve:

* Make the schema folder configurable, or use a proper database backend.
* Validate the schema id, if it results a proper filename (e.g. it does not
  contain '/'). If it's not a proper filename, then the user gets a different
  error than for an empty id at the moment.
* Improve logging. Instead of printing everything, separate INFO / ERROR / DEBUG
  logs, and make the log level configurable.
* Make the request body processing more secure.
* Write additional unit tests for internal functions and stores.
* Use a proper web framework or router module.
