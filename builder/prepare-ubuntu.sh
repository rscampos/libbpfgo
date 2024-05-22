#!/bin/bash

#
# This shell script is meant to prepare a building/exec environment for libbpfgo.
#


# variables

[ -z "${GO_VERSION}" ] && GO_VERSION="1.22"
[ -z "${CLANG_VERSION}" ] && CLANG_VERSION="14"
[ -z "${ARCH}" ] && ARCH=$(uname -m)

case "${ARCH}" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        die "unsupported architecture ${ARCH}"
        ;;
esac


# functions

die() {
    echo "ERROR: ${*}"
    exit 1
}

info() {
    echo "INFO: ${*}"
}

check_tooling() {
    local tools="sudo apt-get"
    for tool in ${tools}
    do
        command -v "${tool}" >/dev/null 2>&1 || die "missing required tool ${tool}"
    done
}

install_pkgs() {
    # silence 'dpkg-preconfigure: unable to re-open stdin: No such file or directory'
    export DEBIAN_FRONTEND=noninteractive

    sudo -E apt-get update || die "coud not update package list"
    for pkg in "${@}"
    do
        info "Installing ${pkg}"
        sudo -E apt-get install -y "${pkg}" || die "could not install ${pkg}"
        info "${pkg} installed"
    done
}

install_clang_format_12() {
    info "Installing clang-format-12"

    echo "deb http://cz.archive.ubuntu.com/ubuntu jammy main universe" | sudo -E tee /etc/apt/sources.list.d/jammy.list
    sudo -E apt-get update
    sudo -E apt-get install -y clang-format-12 || die "could not install clang-format-12"
    sudo -E rm /etc/apt/sources.list.d/jammy.list || die "could not remove jammy.list"
    sudo -E apt-get update
    sudo -E update-alternatives --install /usr/bin/clang-format clang-format /usr/bin/clang-format-12 100

    info "clang-format-12 installed"
}

setup_go() {
    info "Setting Go ${GO_VERSION} as default"
    
    local tools="go gofmt"
    for tool in ${tools}
    do
        sudo -E update-alternatives --install "/usr/bin/${tool}" "${tool}" "/usr/lib/go-${GO_VERSION}/bin/${tool}" 100
    done

    info "Go ${GO_VERSION} set as default"
}

setup_clang() {
    info "Setting Clang ${CLANG_VERSION} as default"

    local tools="clang llc llvm-strip"
    for tool in ${tools}
    do
        sudo -E update-alternatives --install "/usr/bin/${tool}" "${tool}" "/usr/bin/${tool}-${CLANG_VERSION}" 100
    done

    info "Clang ${CLANG_VERSION} set as default"
}


# startup

info "Starting preparation"

check_tooling

install_pkgs \
    coreutils bsdutils findutils \
    build-essential pkgconf \
    golang-"${GO_VERSION}"-go \
    llvm-"${CLANG_VERSION}" clang-"${CLANG_VERSION}" \
    linux-headers-generic \
    linux-tools-generic linux-tools-"$(uname -r)" \
    libbpf-dev libelf-dev libzstd-dev zlib1g-dev

install_clang_format_12

setup_go
setup_clang

info "Preparation finished"
