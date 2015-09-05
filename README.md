[![Travis CI](http://img.shields.io/travis/AntoineAugusti/feature-flags/master.svg?style=flat)](https://travis-ci.org/AntoineAugusti/feature-flags)
# Feature flags API in Go
This package uses [boltdb/bolt](https://github.com/boltdb/bolt), a key-value store for storage. You do not need to connect another database! The HTTP routing is done by [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux).

## What are feature flags?
Feature flags let you enable or disable some features of your application, for example when you're under unexpected traffic or when you want to let some users try a new feature you've been working on. They decouple feature release and code deployment, so that you can release features whenever you want, instead of whenever the code happens to ship.

With this package, you can enable the access of a feature for:
- specific user IDs
- specific groups
- a percentage of your user base
- everyone
- no one

And you can combine things! You can give access to a feature for users in the group `dev` or `admin` and for users `1337` and `42` if you want to.

## Getting started
You can grab this package with the following command:
```
go get gopkg.in/antoineaugusti/feature-flags.v0
```

And then build it:
```
cd ${GOPATH%/}/src/github.com/antoineaugusti/feature-flags
go build
```

## Usage
From the `-h` flag:
```
Usage of ./feature-flags:
  -a string
        address to listen (default ":8080")
  -d string
        location of the database file (default "bolt.db")
```

## Authentication
This API does not ship with an authentication layer. You **should not** expose the API to the Internet. This API should be deployed behind a firewall, only your application servers should be allowed to send requests to the API.

## API Endpoints
- [`GET` /features](#get-features) - Get a list of feature flags
- [`POST` /features](#post-features) - Create a feature flag
- [`GET` /features/:featureKey](#get-featuresfeaturekey) - Get a single feature flag
- [`DELETE` /features/:featureKey](#delete-featuresfeaturekey) - Delete a feature flag
- [`PATCH` /features/:featureKey](#patch-featuresfeaturekey) - Update a feature flag
- [`GET` /features/:featureKey/access](#get-featuresfeaturekeyaccess) - Check if someone has access to a feature

### API Documentation
#### `GET` `/features`
Get a list of available feature flags.
- Method: `GET`
- Endpoint: `/features`
- Responses:
    * **200** on success
    ```json
    [
       {
          "key":"homepage_v2",
          "enabled":false,
          "users":[],
          "groups":[
             "dev",
             "admin"
          ],
          "percentage":0
       },
       {
          "key":"portfolio",
          "enabled":false,
          "users":[
             1337,
             42
          ],
          "groups":[
             "dev",
             "admin"
          ],
          "percentage":50
       }
    ]
    ```
    - `key` is the name of the feature flag
    - `enabled`: tell if the feature flag is enabled. If `true`, everybody has access to the feature flag. Otherwise, the access rule depends on the value of the other attributes.
    - `users`: an array of user IDs who can have access to the feature even if it's disabled.
    - `groups`: an array of group names which can have access to the feature even if it's disabled.
    - `percentage`: a number between 0 and 100. If the percentage is `50`, 50% of the user base is going to have access to the feature.

#### `POST` `/features`
Create a new feature flag.
- Method: `POST`
- Endpoint: `/features`
- Input:
    The `Content-Type` HTTP header should be set to `application/json`

    ```json
   {
      "key":"homepage_v2",
      "enabled":false,
      "users":[],
      "groups":[
         "dev",
         "admin"
      ],
      "percentage":0
   }
    ```
- Responses:
    * 201 Created
    ```json
   {
      "key":"homepage_v2",
      "enabled":false,
      "users":[],
      "groups":[
         "dev",
         "admin"
      ],
      "percentage":0
   }
    ```
    * 422 Unprocessable entity:
    ```json
    {
      "status":"invalid_json",
      "message":"Cannot decode the given JSON payload"
    }
    ```
    * 400 Bad Request
    ```json
    {
      "status":"invalid_feature",
      "message":"<reason>"
    }
    ```
    Common reasons:
    - the feature key already exists. The `message` will be `Feature already exists`
    - the percentage must be between `0` and `100`
    - the feature key must be between `3` and `50` characters
    - the feature key must only contain digits, lowercase letters and underscores

#### `GET` `/features/:featureKey`
Get a specific feature flag.
- Method: `GET`
- Endpoint: `/features/:featureKey`
- Responses:
    * 200 OK
    ```json
   {
      "key":"homepage_v2",
      "enabled":false,
      "users":[],
      "groups":[
         "dev",
         "admin"
      ],
      "percentage":0
   }
    ```
    * 404 Not Found
    ```json
    {
      "status":"feature_not_found",
      "message":"The feature was not found"
    }
    ```

#### `DELETE` `/features/:featureKey`
Remove a feature flag.
- Method: `DELETE`
- Endpoint: `/features/:featureKey`
- Responses:
    * 200 OK
    ```json
    {
      "status":"feature_deleted",
      "message":"The feature was successfully deleted"
    }
    ```
    * 404 Not Found
    ```json
    {
      "status":"feature_not_found",
      "message":"The feature was not found"
    }
    ```

#### `PATCH` `/features/:featureKey`
Update a feature flag.
- Method: `PATCH`
- Endpoint: `/features/:featureKey`
- Input:
    The `Content-Type` HTTP header should be set to `application/json`

    ```json
   {
      "enabled":true,
      "users":[
        13,
        37
      ],
      "groups":[
         "dev"
      ],
      "percentage":42
   }
    ```
- Responses:
    * 200 OK
    ```json
   {
      "key":"homepage_v2",
      "users":[
        13,
        37
      ],
      "groups":[
         "dev"
      ],
      "percentage":42
   }
    ```
    * 404 Not Found
    ```json
    {
      "status":"feature_not_found",
      "message":"The feature was not found"
    }
    ```
    * 422 Unprocessable entity:
    ```json
    {
      "status":"invalid_json",
      "message":"Cannot decode the given JSON payload"
    }
    ```
    * 400 Bad Request
    ```json
    {
      "status":"invalid_feature",
      "message":"<reason>"
    }
    ```
    Common reason:
    - the percentage must be between `0` and `100`

#### `GET` `/features/:featureKey/access`
Check if a feature flag is enabled for a user or a list of groups.
- Method: `GET`
- Endpoint: `/features/:featureKey/access`
- Input:
    The `Content-Type` HTTP header should be set to `application/json`

    ```json
   {
      "groups":[
         "dev",
         "test"
      ],
      "user":42
   }
    ```
- Responses:
    * 200 OK
    ```json
   {
      "status":"has_access",
      "message":"The user has access to the feature"
   }
    ```
    ```json
   {
      "status":"not_access",
      "message":"The user does not have access to the feature"
   }
    ```
    * 404 Not Found
    ```json
    {
      "status":"feature_not_found",
      "message":"The feature was not found"
    }
    ```
    * 422 Unprocessable entity:
    ```json
    {
      "status":"invalid_json",
      "message":"Cannot decode the given JSON payload"
    }
    ```
