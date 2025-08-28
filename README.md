# KVMage: Virt Image Creator

KVMage is an image creation software similar to something like "HashiCorp Packer" that is used to assist with and even fully automate the creation of qcow2 image for use with KVM. What makes KVMage uniique is that it is designed to leverage tools you may already have installed on a KVM hypervisor like virt-install and virt-customize. This is huge when you dont want to install new packages or programs and just want to work with what you have and since its written in Go, its super fast and runs and a single compiled binary on your system. 

## Requirements

Before you just download the binary or compile KVMage on your system, there are a few things you want to make sure that you are able to do on your system. Since KVMage is dependent on existing software on the system, you need to make sure that you are able to execute the following commands. These commands can be installed a number of ways depending on the distro you are using but as long as you are able to execute them, that is all that matters.

```
virt-install
virt-customize
qwemu-img
```

If you are able to successfully execute these commands without any issue, then you are good to get started.

## Installation

Here are a few ways that you are able to install

- Manually compile the code on your system by downloading the repo and using the compile instructions (requires Go)
- Download a precompiled binary [available soon]
- Use our Docker image with everything ready to go [roadmap]


### Manual Installation

Use `git clone` to download the repo locally:

``` Bash
git clone https://gitlab.ctos.io/code/kvmage-virt-image-creator.git
```

``` Bash
cd kvmage-virt-image-creator
mkdir -p dist
bash build.sh
bash install.sh
cd ..
rm -rf kvmage-virt-image-creator
```

Autoinstall script
```
bash <(curl -s https://gitlab.ctos.io/code/kvmage-virt-image-creator/-/raw/main/autoinstall.sh)
```

## How to Use KVMage

KVMage provides a streamlined method for creating qcow images that is designed to feel like a natural extension to KVM by using existing commands and features already available to users with a deployed KVM hypervisor.

### KVMage Build Methods

KVMage has two operating modes:

- `install`: creates a brand-new image using an installation media (an ISO or URL) and a startup script to perform the automated unattended installation of the system.

> **NOTE**
> The only supported methods for install currently are using a Kickstart file with RHEL-based distros such as Fedora, Alma, Rocky, etc...

- `customize`: creates an image using an existing qcow2 image as a source and modifies it with an identified script file (such as bash)

### KVMage Operating Modes

KVMage supports two different methods for operating:

- `run`: Use the `-r, --run` option with the `kvmage` command to perform setup using command line arguments and options.
- `config`: Use the `-f, --config` option with the `kvmage` command to perform setup using a config file (YAML) where the options are defined. Config mode is particularly useful for managing multiple image builds as you can stack builds.

### KVMage Install

RUN Example:
```bash
kvmage \
    --run \
    --install \
    --image-name almalinux01 \
    --os-var almalinux9 \
    --image-size 100G \
    --ks-file ks.cfg \
    --install-media almalinux9.5-minimal.iso \
    --image-dest $PWD
```

CONFIG Example:
```yaml
---
kvmage:
  almalinux9:
    image_name: almalinux01
    virt_mode: install
    os_var: almalinux9
    image_size: 100G
    ks_file: ks.cfg
    install_media: almalinux9.5-minimal.iso
    image_dest: $PWD
```

### KVMage Customize

RUN Example:

```bash
kvmage \
    --run \
    --customize \
    --image-name almalinux02 \
    --os-var almalinux9 \
    --image-src almalinux01.qcow2 \
    --image-dest $PWD \
    --custom-script script.sh
```

```yaml
---
kvmage:
  almalinux02:
    image_name: almalinux02
    virt_mode: customize
    os_var: almalinux9
    image_src: almalinux01.qcow2
    img_dest: $PWD
    custom_script: script.sh
```
### KVMage Options

> **NOTE**
> The following options are not functional (currently): verbose, debug

```
kvamge [command]

-h, --help (help menu)
-v, --verbose (enabled verbose for kvmage)
-d, --debug (enables --debug flag for virt-install and virt-customize)
-V, --version (prints version information for kvamge, kvm, )

(kvmage requires run OR config but not both)
-r, --run [command args]
-f, --config <config_file>

(kvmage requires install or config but not both)
-i, --install | virt_mode: install (requires image_name, os_var, image_size, ks_file, install_media, and image_dest)
-c, --customize | virt_mode: customize (requires image_name, os_var, image_src, image_dest)

-n, --image-name | image_name <name>
-o, --os-var <os> | os_var: <os> (this comes from osinfo-query os)
-s, --image-size <size> | image_size: <size> (eg. 100G)
-k, --ks-file | ks_file: <ks_file> 
-l, --install-media | install_media: <media> (virt-install --location)


-H, --hostname <hostname> | hostname: <hostname> (optional)
-C, --custom-script | custom_script: <script_file> (bash file, optional but should be included with config option)
-S, --image-src <src_qcow> | image_src: <src_qcow>
-D, --image-dest <dest_qcow> | image_dest: <dest_qcow>
-W, --network <iface> | network: <iface> (optional)
```