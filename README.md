# Container Linux Config Transpiler

The config transpiler ("ct" for short) is the utility responsible for transforming a Container Linux instance's configuration from the human-friendly, YAML form into an [Ignition][ignition] configuration. While it is possible to write Ignition configs directly, CoreOS recommends that this tool be used instead. More details about the reasoning and the overall design of Ignition and this project are provided in the [documentation][overview].

[ignition]: https://github.com/coreos/ignition
[overview]: doc/overview.md

## Building

```shell
git clone --branch v0.1.0 https://github.com/coreos/container-linux-config-transpiler
cd container-linux-config-transpiler
./build
```

The `ct` binary will be placed in `./bin/`.
