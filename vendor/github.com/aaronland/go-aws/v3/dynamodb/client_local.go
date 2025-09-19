package dynamodb

// LOCAL_CLIENT_URI is a suitable URI for passing to the `NewClient` and `NewClientV1` methods for use with a local instance of DynamoDB.
const LOCAL_CLIENT_URI string = "aws://?region=localhost&credentials=anon:&local=true"
