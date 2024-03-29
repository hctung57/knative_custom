# Test

This directory contains tests and testing docs for `Knative Serving`:

- [Unit tests](#running-unit-tests) currently reside in the codebase alongside
  the code they test
- [End-to-end tests](#running-end-to-end-tests), of which there are two types:
  - Conformance tests in [`/test/conformance`](./conformance)
  - Other end-to-end tests in [`/test/e2e`](./e2e)
- [Performance tests](#running-performance-tests) reside in
  [`/test/performance`](./performance)

The conformance tests are a subset of the end to end test with
[more strict requirements](./conformance/README.md#requirements) around what can
be tested.

If you want to add more tests, see [adding_tests.md](./adding_tests.md).

## Presubmit tests

[`presubmit-tests.sh`](./presubmit-tests.sh) is the entry point for both the
[end-to-end tests](./e2e) and the [conformance tests](./conformance)

This script, and consequently, the e2e and conformance tests will be run before
every code submission. You can run these tests manually with:

```shell
test/presubmit-tests.sh
```

_Note that to run `presubmit-tests.sh` or `e2e-tests.sh` scripts, you'll need
kubernetes `kubetest` installed:_

```bash
go get -u k8s.io/test-infra/kubetest
```

## Running unit tests

To run all unit tests:

```bash
go test ./...
```

_By default `go test` will not run [the e2e tests](#running-end-to-end-tests),
which need [`-tags=e2e`](#running-end-to-end-tests) to be enabled._

## Running end to end tests

To run [the e2e tests](./e2e), you need to have a running environment that meets
[the e2e test environment requirements](#environment-requirements), and you need
to specify the build tag `e2e`.

```bash
go test -v -tags=e2e -count=1 ./test/e2e
```

## Running conformance tests

To run [the conformance tests](./conformance), you need to have a running environment that meets
[the e2e test environment requirements](#environment-requirements), and you need
to specify the build tag `e2e`.

```bash
go test -v -tags=e2e -count=1 ./test/conformance/...
```

### Running conformance as a project admin

It is possible to run the conformance tests by a user with reduced privileges, e.g. project admin.
The environment needs to meet similar requirements as in [running conformance tests](#running-conformance-tests)
but the test resources can be limited to minimum and so the user can install them.
Running the conformance tests then consists of these steps:
1. The cluster admin creates the cluster scope resources
   ```bash
   kubectl apply -f test/config/cluster-resources.yaml
   ```
1. The project admin installs minimum test resources:
    ```bash
    kubectl apply -f test/config/test-resources.yaml
    ```
1. The project admin then runs the conformance test suite using the `--disable-logstream` flag:
    ```bash
    go test -v -tags=e2e -count=1 ./test/conformance/... \
      -disable-logstream \
      -kubeconfig=$PROJECT_ADMIN_KUBECONFIG
    ```

The tests can be run in arbitrary test namespaces. If you modify the namespaces of the
resources above you must pass the updated values to the Golang test suite:
```bash
go test -tags=e2e -count=1 ./test/conformance/... \
  -disable-logstream \
  -kubeconfig=$PROJECT_ADMIN_KUBECONFIG \
  -test-namespace=serving-tests \
  -alt-test-namespace=serving-tests-alt \
  -tls-test-namespace=tls
```

## Running performance tests

Performance tests using [kperf](https://github.com/knative-extensions/kperf) can be run in
a pull request using the optional test job `/test pull-knative-serving-performance-tests-kperf`

A report will be generated in the artifacts folder of the test run.

### Running a single test case

To run one e2e test case, e.g. TestAutoscaleUpDownUp, use
[the `-run` flag with `go test`](https://golang.org/cmd/go/#hdr-Testing_flags):

```bash
go test -v -tags=e2e -count=1 ./test/e2e -run ^TestAutoscaleUpDownUp$
```

### Running tests in short mode

Running tests in short mode excludes some large-scale E2E tests and saves
time/resources required for running the test suite. To run the tests in short
mode, use
[the `-short` flag with `go test`](https://golang.org/cmd/go/#hdr-Testing_flags)

```bash
go test -v -tags=e2e -count=1 -short ./test/e2e
```

To get a better idea where the flag is used, search for `testing.Short()`
throughout the test source code.

### Environment requirements

These tests require:

1. [A running `Knative Serving` cluster.](../DEVELOPMENT.md#prerequisites)
1. The `knative-testing` resources:

   ```bash
   ko apply -f test/config
   ```

1. A docker repo containing [the test images](#test-images)

### Common Flags

- By default the e2e tests against the current cluster in `~/.kube/config` using
  the environment specified in
  [your environment variables](../DEVELOPMENT.md#setup-your-environment).
- Since these tests are fairly slow, running them with logging enabled is
  recommended (`-v`).
- Using `-count=1` is
  [the idiomatic way to disable test caching](https://golang.org/doc/go1.10#test)

You can [use test flags](#flags) to control the environment your tests run
against, i.e. override
[your environment variables](../DEVELOPMENT.md#setup-your-environment):

```bash
go test -v -tags=e2e -count=1 ./test/conformance/... --kubeconfig ~/special/kubeconfig --cluster myspecialcluster --dockerrepo myspecialdockerrepo
go test -v -tags=e2e -count=1 ./test/e2e --kubeconfig ~/special/kubeconfig --cluster myspecialcluster --dockerrepo myspecialdockerrepo
```

## Test images

### Building the test images

Note: this is only required when you run conformance/e2e tests locally with
`go test` commands.

The [`upload-test-images.sh`](./upload-test-images.sh) script can be used to
build and push the test images used by the conformance and e2e tests. The script
expects your environment to be setup as described in
[DEVELOPMENT.md](../DEVELOPMENT.md#install-requirements).

To run the script for all end to end test images:

```bash
./test/upload-test-images.sh
```

A docker tag may be passed as an optional parameter. This can be useful on
Minikube in tandem with the `--tag` [flag](#using-a-docker-tag):

`PLATFORM` environment variable is optional. If it is specified, test images
will be built for specific hardware architecture, according to its value (for
instance,`linux/arm64`).

```bash
eval $(minikube docker-env)
./test/upload-test-images.sh any-old-tag
```

### Adding new test images

New test images should be placed in `./test/test_images`.

## Flags

These flags are useful for running against an existing cluster, making use of
your existing [environment setup](../DEVELOPMENT.md#setup-your-environment).

Tests importing [`knative.dev/serving/test`](#test-library) recognize these
flags:

- [All flags added by `knative/pkg/test`](https://github.com/knative/pkg/tree/main/test#flags)
  such as:
  - [`--dockerrepo`](#overriding-docker-repo)
  - [`--tag`](#using-a-docker-tag)
  - [`--ingressendpoint`](#using-a-custom-ingress-endpoint)
- [`--resolvabledomain`](#using-a-resolvable-domain)
- [`--https`](#using-https)

### Overriding docker repo

The `--dockerrepo` argument lets you specify the docker repo from which images
used by your tests should be pulled. This will default to the value of your
[`KO_DOCKER_REPO` environment variable](../DEVELOPMENT.md#setup-your-environment)
if not specified.

```bash
go test -v -tags=e2e -count=1 ./test/conformance/... --dockerrepo gcr.myhappyproject
go test -v -tags=e2e -count=1 ./test/e2e --dockerrepo gcr.myhappyproject
```

### Using a docker tag

The default docker tag used for the test images is `latest`, which can be
problematic on Minikube. To avoid having to configure a remote container
registry to support the `Always` pull policy for `latest` tags, you can have the
tests use a specific tag:

```bash
go test -v -tags=e2e -count=1 ./test/conformance/... --tag any-old-tag
go test -v -tags=e2e -count=1 ./test/e2e --tag any-old-tag
```

Of course, this implies that you tagged the images when you
[uploaded them](#building-the-test-images).

### Using a custom ingress endpoint

Some environments (like minikube) do not support a Loadbalancer to make Knative
services externally available. These environments usually rely on rewriting the
Loadbalancer to a NodePort. The external address of such a NodePort is usually
not easily obtained within the cluster automatically, but can be provided from
the outside through the `--ingressendpoint` flag. For a minikube setup for
example, you'd want to run tests against the default `ingressgateway` (port
number 31380) running on the minikube node:

```
go test -v -tags=e2e -count=1 ./test/conformance/... --ingressendpoint "$(minikube ip):31380"
go test -v -tags=e2e -count=1 ./test/e2e --ingressendpoint "$(minikube ip):31380"
```

### Using a resolvable domain

If you set up your cluster using
[the getting started docs](../DEVELOPMENT.md#prerequisites), Routes created in
the test will use the domain `example.com`, unless the route has label
`app=prod` in which case they will use the domain `prod-domain.com`. Since these
domains will not be resolvable to deployments in your test cluster, in order to
make a request against the endpoint, the test use the IP assigned to the service
`istio-ingressgateway` in the namespace `istio-system` and spoof the `Host` in
the header.

If you have configured your cluster to use a resolvable domain, you can use the
`--resolvabledomain` flag to indicate that the test should make requests
directly against `Route.Status.Domain` and does not need to spoof the `Host`.

### Overriding the gateway used for spoofing

If you are using an ingress provider other than Istio, and have not set up a
resolvable domain (above), you will also need to set the `GATEWAY_OVERRIDE` and
`GATEWAY_NAMESPACE_OVERRIDE` environment variables to allow the gateway to be
discovered for spoofing. For example, for kourier you would do:

```
export GATEWAY_OVERRIDE=kourier
export GATEWAY_NAMESPACE_OVERRIDE=kourier-system
```

### Using https

You can use the `--https` flag to have all tests run with https.
