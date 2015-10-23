# gopivnet
A go library and cli executable to download a product from network.pivotal.io

# Get It

`go get github.com/cfmobile/gopivnet`

# Usage

```
gopivnet -help
Usage of gopivnet:
  -file="": filename where to save the pivotal product
  -product="": product to download
  -token="": pivnet token
  -version="": version of the product. If missing download the latest version
```

Example: `gopivnet -product p-redis -token <token> -version "1.4.7" -file p-redis.pivotal`

# Fetching a pivnet token

https://network.pivotal.io/docs/api

Relevane section: 
> Some products and releases may require authentication to access or modify. Your Pivotal Network API Token can be found on your [Edit Profile](https://network.pivotal.io/users/dashboard/edit-profile) page. This API token should be used in the Authorization header of the request. The Authentication API can be used to test that you are using your authorization token correctly.

# Library

The api package is meant to make it simple to fetch a pivotal product of a specific version and download it.
