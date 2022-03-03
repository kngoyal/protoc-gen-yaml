## Overview
Create a Protocol Buffers plugin, which is executed with the `protoc` compiler.

## What is a plugin, and how do I write one?

A [protoc plugin][3] allows users to generate custom code according
to the *structured representation* of Protobuf files. In simpler terms,
plugins are programs that take Protobuf files as input, and produce code
as output.

As you'll notice in some of the [additional plugin documentation][4],
the plugin *must* either be discoverable by your machine's `PATH`, or
be manually targeted with the `--plugin` flag. In this case, we recommend
placing the plugin on your `PATH` so that you can quickly `go install`
between test invocations, but this is up to you.

For references, there are many examples of plugins you can find on the web.

  [3]: https://developers.google.com/protocol-buffers/docs/reference/other#plugins
  [4]: https://developers.google.com/protocol-buffers/docs/reference/cpp/google.protobuf.compiler.plugin

## What does the plugin need to do?

The plugin *must* be implemented in Go.

The plugin *must* output one YAML file *for each* of the input files specified
in the `protoc` invocation.

The YAML file *must* conform to the following constraints:
* The plugin is named `protoc-gen-yaml`.
* The plugin outputs one YAML file for every input file.
* Each of the files should be named by adding the `.yaml` suffix to their filename.
  * For example, an input file named `proto/echo.proto` will create `proto/echo.proto.yaml`.
* The YAML output will represent a subset of the Protobuf file's content. Specifically,
  the plugin will create two `maps` in each file: one for `messages` and another for `services`.
  Each map key have a `list` of `messages` and `services`, respectively.
* For each message, a `list` of fields will be represented under the `fields` key.
* For each method, a `list` of methods will be represented under the `methods` key.

A template representation of the expected output is shown below. *Pay close attention to
the names of the map keys, as well as the required values associated with them.*

> If your are confused at all by any of these definitions, such as a `field` number,
  please refer to the [Language Guide][8]. The [Testing](#testing) section below gives
  an example that will make it easier to understand.

```yaml
messages:
# For each message in this file, sorted by name...
- name: <$message_name>
  fields:
  # For each field in this message, sorted by number...
  - name: <$field_name>
    number: <$field_number>
services:
# For each service in this file, sorted by name...
- name: <$service_name>
  # For each method in this service, sorted by name...
  methods:
  - name: <$method_name>
    input_type: <$input_type>
    output_type: <$output_type>
```

  [8]: https://developers.google.com/protocol-buffers/docs/overview

> Hint: Use `Marshal` from gopkg.in/yaml.v2 instead of using gopkg.in/yaml.v3 to get the correct indentation.

## Testing

The following file located in `test/input/proto/echo.proto`.

```
syntax = "proto3";

package echo.v1;

message EchoRequest {
  string value = 1;
}

message EchoResponse {
  Foo foo = 1;
}

message Foo {
  message Bar {
    string value = 1;
  }
  string two = 2;
  string one = 1;
}

service EchoService {
  rpc Echo(EchoRequest) returns (EchoResponse);
}
```

Assuming that your plugin is installed in `PATH`, we should be able to execute
the plugin like so:

```shell
$ protoc -I test/input test/input/proto/echo.proto --yaml_out=./output
```

The expected output for the `test/input/proto/echo.proto` file is shown below:

```yaml
messages:
- name: echo.v1.EchoRequest
  fields:
  - name: value
    number: 1
- name: echo.v1.EchoResponse
  fields:
  - name: foo
    number: 1
- name: echo.v1.Foo
  fields:
  - name: one
    number: 1
  - name: two
    number: 2
- name: echo.v1.Foo.Bar
  fields:
  - name: value
    number: 1
services:
- name: echo.v1.EchoService
  methods:
  - name: Echo
    input_type: echo.v1.EchoRequest
    output_type: echo.v1.EchoResponse
```

You can verify if you've implemented the correct solution with `diff`:

```shell
$ protoc -I test/input test/input/proto/echo.proto --yaml_out=./output
$ diff ./output/proto/echo.proto.yaml ./test/output/proto/echo.proto.yaml
```

> Note that the ouptut is *sanitized*, such that any `.` prefixes
  are trimmed in the relevant places. A proper solution will need
  to handle these cases, so pay special attention!

## Submitting

Send a tar ball containing your solution and, optionally, any additional test cases you've created.
We should be able to build the plugin with `go build`.
