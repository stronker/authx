# Authx
​
The authx component is responsible of managing the credentials of the different elements of the platform that require
those (e.g., users). Thus, this component will be used by either components trying to check a set of credentials,
or by components creating/managing credentials.
​
## Getting Started
​
The authx component is a requirement for having a running platform, therefore any error in the deployment may cause
issue at system level.
​
### Prerequisites
​
The following components are required
​
* [scylla-deploy](https://github.com/nalej/scylladb-deploy)
* A secret named `authx-secret` is required and created by the [installer](https://github.com/nalej/installer). The content
of this secret is used to create the JWT tokens so different installations are expected to use different secrets.
* A certification authoritity created by the [installer](https://github.com/nalej/installer) is required to issue new certificates.
​
### Build and compile
​
In order to build and compile this repository use the provided Makefile:
​
```
make all
```
​
This operation generates the binaries for this repo, download dependencies,
run existing tests and generate ready-to-deploy Kubernetes files.
​
### Run tests
​
Tests are executed using Ginkgo. To run all the available tests:
​
```
make test
```
​
### Update dependencies
​
Dependencies are managed using Godep. For an automatic dependencies download use:
​
```
make dep
```
​
In order to have all dependencies up-to-date run:
​
```
dep ensure -update -v
```
​
## Known Issues

* The interceptors are being migrated to their [own repository](https://github.com/nalej/authx-interceptors) to limit inter-repository dependencies.
* Some authx packages do not follow the handler-manager approach as implemented in other repositories. This issue will be part of a refactor on future
versions.
* Some scylladb providers do not make use of the [scylladb-utils](https://github.com/nalej/scylladb-utils). This will be refactored in future versions.
​
​
## Contributing
​
Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.
​
​
## Versioning
​
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/nalej/authx/tags). 
​
## Authors
​
See also the list of [contributors](https://github.com/nalej/authx/contributors) who participated in this project.
​
## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.

