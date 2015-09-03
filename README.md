# Feature flags API in Go
Documentation on its way.

## API Endpoints
- [`GET` /features](#get-features)
- [`POST` /features](#post-features)
- [`GET` /features/:featureKey](#get-featuresfeaturekey)
- [`DELETE` /features/:featureKey](#delete-featuresfeaturekey)
- [`PATCH` /features/:featureKey](#patch-featuresfeaturekey)
- [`GET` /features/:featureKey/access](#get-featuresfeaturekeyaccess)

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

#### `DELETE` `/features`
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
        "foo",
        "bar"
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
        "foo",
        "bar"
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
    Common reasons:
    - the feature key already exists. The `message` will be `Feature already exists`
    - the percentage must be between `0` and `100`
    - the feature key must be between `3` and `50` characters
    - the feature key must only contain digits, lowercase letters and underscores

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
