# This project is a WIP and is not in a working state yet.

# S3 Policy tester

S3 Policy tester is a web app (with accompanying docker file) used to test S3 Policies and is intended to be used with automated tests for testing IAM roles.

Running the application uses the default credentials on the system to test if the current setup allows access to:

- List all buckets in the current application
- Put an object within the bucket
- Get an object from the bucket

The application is intended to run from ECS as well as EC2 if needed, because of this configuration is done through environment variables.

## Environment variables

**These may change when the application is finalised**

* `AWS_ACCESS_KEY_ID` (Optional if supplied by attached IAM Role) AWS access key - used to authenticate access to the bucket, supplied by either instance/container configuration or the IAM role
* `AWS_SECRET_ACCESS_KEY` (Optional if supplied by attached IAM Role) AWS secret key - See above
* `CHECK_POLICIES` (Optional, defaults to `get, put`) A comma delimited list of policies to check from the following:
  * `listall` List all buckets in the account
  * `put` Putting an object in the specified bucket
  * `get` Getting an object from the specified bucket
* `S3_BUCKET_NAME` Bucket to test the policies against, should be created prior to S3 Policy Tester running
* `S3_BUCKET_OBJECT` (Optional if `put` policy check is enabled) A relative path to an object to check for checking get object policies, this can be omitted if the put object policies are enabled
* `SERVER_PORT` (Optional, defaults to `80`) what port the server listens on