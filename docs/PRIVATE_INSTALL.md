# Private Provider Installation

Use this process for pre-release testing, private mirrors, or environments that cannot install directly from the Terraform Registry.

Terraform configurations should still use the final provider address:

```hcl
terraform {
  required_providers {
    mailu = {
      source  = "cetic-group/mailu"
      version = "0.1.0-rc.1"
    }
  }
}
```

## Install From GitHub Release Assets

Download the archive that matches your platform from the GitHub release assets.

Use one of these platform directory names:

- `darwin_amd64`
- `darwin_arm64`
- `linux_amd64`
- `linux_arm64`
- `windows_amd64`
- `windows_arm64`

For macOS or Linux:

```shell
VERSION="0.1.0-rc.1"
OS_ARCH="darwin_arm64"

mkdir -p "$HOME/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$VERSION/$OS_ARCH"
unzip "terraform-provider-mailu_${VERSION}_${OS_ARCH}.zip" \
  -d "$HOME/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$VERSION/$OS_ARCH"
chmod +x "$HOME/.terraform.d/plugins/registry.terraform.io/cetic-group/mailu/$VERSION/$OS_ARCH"/terraform-provider-mailu_v*
```

For Windows, use the same directory structure under:

```text
%APPDATA%\terraform.d\plugins\registry.terraform.io\cetic-group\mailu\0.1.0-rc.1\windows_amd64
```

Then run:

```shell
terraform init
terraform providers
```

`terraform providers` should show `registry.terraform.io/cetic-group/mailu`.

## Verify Checksums

Each release includes a `SHA256SUMS` file. Verify the downloaded archive before installing it:

```shell
shasum -a 256 terraform-provider-mailu_0.1.0-rc.1_darwin_arm64.zip
grep terraform-provider-mailu_0.1.0-rc.1_darwin_arm64.zip terraform-provider-mailu_0.1.0-rc.1_SHA256SUMS
```

The checksum values must match.

## Current Limitations

- Terraform cannot download this provider automatically from a private GitHub release using only `source = "cetic-group/mailu"`.
- For private distribution, use local plugin installation or a provider mirror.
- Prefer Terraform Registry installation for stable public releases.
