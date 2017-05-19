# SQS Push

[![Build Status](https://travis-ci.org/SSENSE/sqspush.svg?branch=master)](https://travis-ci.org/SSENSE/sqspush)

simple cli command that push STDIN to your AWS SQS queue

---
## Install
[Binary packages](https://github.com/SSENSE/sqspush/releases) are available for Mac, Linux and Windows.

To build from source you can:

1. Clone this repository into `$GOPATH/src/github.com/SSENSE/sqspush` and
   change directory into it
2. Run `make build`

This will leave you with `./sqspush`, which you can put in your `$PATH` if
you'd like. (You can also take a look at `make install` to install for
you.)


## Usage

```bash
# login to AWS
sqspush login --key=aws_access_key_id --secret=aws_secret_access_key

# Send payload.json to your queue
cat ./mock/payload.json | sqspush --queue=https://sqs.us-west-2.amazonaws.com/XXYYYZZZZXXX/queue-name.fifo --region=us-west-2 --group=mygroup

# or use echo to send data
echo '{ "foo": "bar" }' | sqspush --queue=https://sqs.us-west-2.amazonaws.com/XXYYYZZZZXXX/queue-name --region=us-west-2 --group=mygroup

# Getting help
sqspush --help
# NAME:
#    SQS PUSH - push stdin into sqs queue
# USAGE:
#    cat ./payload.json | sqspush [global options] command [command options] [arguments...]
#
# COMMANDS:
#    login, l  login aws credentials
#    help, h   Shows a list of commands or help for one command
#
# GLOBAL OPTIONS:
#    --cred-path value, -c value  AWS credentials path (default: "~/.aws/credentials")
#    --queue value, -q value      SQS queue URL
#    --region value, -r value     SQS queue REGION
#    --profile value, -p value    AWS login profile
#    --retries value              SQS queue retries (default: 2)
#    --group value, -g value      SQS queue MessageGroupId (default: "5ePlmiO61JA5iqU1jBkjLYv2Xp4=")
#    --help, -h                   show help
#    --version, -v                print the version
```

### Configuring Credentials
Before using the SQS PUSH, ensure that you've configured credentials. The best way to configure credentials on a development machine is to use the `~/.aws/credentials` file, which might look like:

```bash
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

or use SQS PUSH login command to create your `~/.aws/credentials`

```bash
# login to AWS will create or overwrite ~/.aws/credentials file
sqspush login --key=aws_access_key_id --secret=aws_secret_access_key
```

### LICENSE

This package is made available under an MIT-style license. See LICENSE.txt.
