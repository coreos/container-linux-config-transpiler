# Container Linux Config Transpiler

The Config Transpiler ("ct" for short) is the utility responsible for transforming a human-friendly Container Linux Config into a JSON file. This resulting file can be provided to a Container Linux machine when it first boots to provision the machine.

**NOTE: This tool is for Container Linux, not Fedora CoreOS. See [FCCT](https://github.com/coreos/fcct) for the Fedora CoreOS equivalent.**

## Documentation

If you're looking to begin writing configs for your Container Linux machines, check out the [getting started][get-started] documentation.

The [configuration][config] documentation is a comprehensive resource specifying what options can be in a Container Linux Config.

For a more in-depth view of ct and why it exists, take a look at the [Overview][overview] document.

Please use the [bug tracker][issues] to report bugs.

[ignition]: https://github.com/coreos/ignition
[issues]: https://issues.coreos.com
[overview]: doc/overview.md
[get-started]: doc/getting-started.md
[config]: doc/configuration.md

## Examples

There are plenty of small, self-contained examples [in the documentation][examples].

[examples]: doc/examples.md

## Installation

### macOS homebrew

```bash
brew install coreos-ct
```

### Prebuilt binaries

The easiest way to get started using ct is to download one of the binaries from the [releases page on GitHub][releases].

One can use the following script to download and verify the signature of Config Transpiler:

```bash
# Specify Config Transpiler version
CT_VER=v0.6.1

# Specify Architecture
# ARCH=aarch64 # ARM's 64-bit architecture
ARCH=x86_64

# Specify OS
# OS=apple-darwin # MacOS
# OS=pc-windows-gnu.exe # Windows
OS=unknown-linux-gnu # Linux

# Specify download URL
DOWNLOAD_URL=https://github.com/coreos/container-linux-config-transpiler/releases/download

# Remove previous downloads
rm -f /tmp/ct-${CT_VER}-${ARCH}-${OS} /tmp/ct-${CT_VER}-${ARCH}-${OS}.asc /tmp/coreos-app-signing-pubkey.gpg

# Download Config Transpiler binary
curl -L ${DOWNLOAD_URL}/${CT_VER}/ct-${CT_VER}-${ARCH}-${OS} -o /tmp/ct-${CT_VER}-${ARCH}-${OS}
chmod u+x /tmp/ct-${CT_VER}-${ARCH}-${OS}

# Download and import CoreOS application signing GPG key
curl https://coreos.com/dist/pubkeys/app-signing-pubkey.gpg -o /tmp/coreos-app-signing-pubkey.gpg
gpg2 --import --keyid-format LONG /tmp/coreos-app-signing-pubkey.gpg

# Download and import CoreOS application signing GPG key if it has not already been imported
curl -L ${DOWNLOAD_URL}/${CT_VER}/ct-${CT_VER}-${ARCH}-${OS}.asc -o /tmp/ct-${CT_VER}-${ARCH}-${OS}.asc
gpg2 --verify /tmp/ct-${CT_VER}-${ARCH}-${OS}.asc /tmp/ct-${CT_VER}-${ARCH}-${OS}
```

[releases]: https://github.com/coreos/container-linux-config-transpiler/releases

### Building from source

To build from source you'll need to have the go compiler installed on your system.

```shell
git clone --branch v0.8.0 https://github.com/coreos/container-linux-config-transpiler
cd container-linux-config-transpiler
make
```

The `ct` binary will be placed in `./bin/`.

Note: Review releases for new branch versions.

## Related projects

- [https://github.com/coreos/ignition](https://github.com/coreos/ignition)
- [https://github.com/coreos/coreos-metadata/](https://github.com/coreos/coreos-metadata/)
- [https://github.com/coreos/matchbox](https://github.com/coreos/matchbox)
