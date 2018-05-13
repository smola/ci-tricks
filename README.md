# ci-tricks [![Build Status](https://travis-ci.org/smola/ci-tricks.svg?branch=master)](https://travis-ci.org/smola/ci-tricks)  [![Build status](https://ci.appveyor.com/api/projects/status/github/smola/ci-tricks?branch=master&svg=true)](https://ci.appveyor.com/project/smola/ci-tricks)

**ci-tricks** is a single-binary that can be run on multiple continuous integration providers to setup the environment in an efficient way.

**WARNING: This approach is hacky, dirty and ugly. It is meant to provide some last resort hacks for projects that run on multiple CI providers and platforms.**

## Supported platforms

* Appveyor/Windows/amd64
* Travis/Linux/amd64 ([sudo-enabled](https://docs.travis-ci.com/user/reference/overview/#Sudo-enabled))
* Travis/macOS/amd64

## Usage

1. Set the environment variables as required by each trick.
2. Run `ci-tricks`.

## Tricks

### PostgreSQL

PostgreSQL can be set up by setting the `POSTGRESQL_VERSION` environment variable to the desired PostgreSQL version.

| Version       | Travis/Linux  | Travis/macOS | Appveyor/Windows      |
| ------------- |:-------------:|:------------:|:---------------------:|
| `9.2`         | ✅             | ❌           | ❌                     |
| `9.3`         | ✅             | ❌           | ❌                     |
| `9.4`         | ✅             | ✅           | ❌                     |
| `9.5`         | ✅             | ✅           | ✅                     |
| `9.6`         | ✅             | ✅           | ✅                     |
| `10`          | ❌             | ✅           | ✅                     |

### RabbitMQ

RabbitMQ can be set up by setting the `RABBITMQ_VERSION` environment variable to the desired PostgreSQL version.

| Version       | Travis/Linux  | Travis/macOS | Appveyor/Windows      |
| ------------- |:-------------:|:------------:|:---------------------:|
| `any`         | ✅             | ✅           | ✅                     |

## Contributing

[Issues](https://github.com/smola/ci-tricks/issues) and [pull requests](https://github.com/smola/ci-tricks/pulls) are welcome.

**Platforms:** We are **not** accepting issues or pull requests related to running ci-tricks outside a continuous integration provider.

**Bugs:** If you find a bug, please, include a link to an affected CI build.

**Features:** We are open to include new tricks. For new services, we are initially interested only in those supported either [by Travis](https://docs.travis-ci.com/user/database-setup/) or [by Appveyor](https://www.appveyor.com/docs/services-databases/).

## License

This project is released under the terms of the Apache License 2.0.
