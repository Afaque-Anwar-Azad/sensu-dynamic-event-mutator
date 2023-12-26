[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/Afaque-Anwar-Azad/sensu-dynamic-event-mutator)
![Go Test](https://github.com/Afaque-Anwar-Azad/sensu-dynamic-event-mutator/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/Afaque-Anwar-Azad/sensu-dynamic-event-mutator/workflows/goreleaser/badge.svg)


# sensu-dynamic-event-mutator

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Mutator definition](#mutator-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview


The Sensu dynamic event mutator enriches the event payload by seamlessly incorporating 
details of both labels and annotations. These specific labels and annotations information
are sourced from the associated entity(event.entity.name), enriching the overall
contextual information of the event payload.

Key enhancements to the event fields:
event.entity.labels
event.entity.annotations


## Usage examples

Usage:
  sensu-dynamic-event-mutator [command]

example commands:
1. sensu-dynamic-event-mutator --add-labels
2. sensu-dynamic-event-mutator --add-annotations
3. sensu-dynamic-event-mutator --add-all

### Environment variables

| Argument          | Environment Variable  |
|-------------------|-----------------------|
| --api-url         | SENSU_API_URL         |
| --api-key         | SENSU_API_KEY         |

**Security Note:** Care should be taken to not expose the API key or access token for this handler by explicitly specifying either on the command line or by directly setting the environment variable(s) in the handler definition.
It is suggested to make use of [secrets management][3] to provide the API key or access token as environment variables.



## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add Afaque-Anwar-Azad/sensu-dynamic-event-mutator
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/Afaque-Anwar-Azad/sensu-dynamic-event-mutator].

### Mutator definition

```yml
---
type: Mutator
api_version: core/v2
metadata:
  name: sensu-dynamic-event-mutator
  namespace: default
spec:
  command: sensu-dynamic-event-mutator --add-labels
  runtime_assets:
  - Afaque-Anwar-Azad/sensu-dynamic-event-mutator
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-dynamic-event-mutator repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[6]: https://docs.sensu.io/sensu-go/latest/reference/mutators/
[8]: https://bonsai.sensu.io/
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
