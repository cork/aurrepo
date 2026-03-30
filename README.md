# aurrepo is a simple manager for building and updating packages in a custom archlinux repo


## Dependencies

- archlinux base
- pacman
- devtools
- git (aur repos and upstream)
- bash (for running pkgver to check for new versions)
- gpg (for signing)

## Building

```bash
go mod tidy
go build
```

# Setup

## GPG

Create a gpg key following https://www.sainnhe.dev/post/create-personal-arch-linux-package-repository/#create-your-gpg-key

Configure makepkg following https://www.sainnhe.dev/post/create-personal-arch-linux-package-repository/#configure-makepkg

## AUR folder

### Build archive

Create the folder used in --aur (~/.cache/aurrepo/aur by default).

git clone the different aur packages you want to build packages for.

### install.json

If a package requires other aur packages to build add them in a install.json file and a glob matching the archive in the repo. (glob means you can use whildcards for version numbers).

#### Example

```json
{
    "tuxedo-control-center-bin": ["tuxedo-drivers-dkms-*-x86_64.pkg.tar.zst"]
}
```

NOTE: You must build the package once before it can be used in another build, aurrepo won't resolve dependencies.

## Repo

The repo given to aurrepo will be created automatically on first package it finds needs to be built.

When the repo folder is populated share it over a http server of your choice and add following to your clients pacman.conf:
```ini
[aurrepo]
Server=https://host-name.org/$arch
```

## Build repo
run ./aurrepo and it will check each package for aur updates and upstream changes following the PKGBUILD.
