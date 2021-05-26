# OpenShift Dynamic Target Handler

## What is it?

This is an OpenShift target handler for a [Dynamic Target Registration Server](https://github.com/jacobsee/dynamic-target-registration-server). It is intended to run alongside a Prometheus instance, and will keep a directory of target files up-to-date with the most current OpenShift beacons.

## How is it configured?

This server takes the following environment variables:

| Variable Name | Description | Example |
| --- | --- | --- |
| `SERVER_URL` | A fully-formed URL (with protocol prefix) at which the registration service can be reached (defaults to `http://localhost:8081` if unset) | `https://1.2.3.4` |
| `TARGET_PATH` | The location on disk to which dynamic OpenShift targets should be written as Prometheus-parseable target files | `/targets/openshift/` |
| `AUTH_TOKEN` | The token to pass in an `Authorization` header to the registration service | `12345abcde` |
