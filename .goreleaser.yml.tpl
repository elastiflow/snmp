[[/*
  Following values are substituted by https://github.com/hairyhenderson/gomplate:
  [[ .Env.BINARY ]] - binary name, used in repo paths
  [[ .Env.PKG_NAME ]] - used to set package name in DEB/RPM manifest. This name will appear in the apt/yum repositories
  [[ .Env.PKG_DESCRIPTION ]] - package description
*/]]

[[ $etcDirs := coll.Slice
    "autodiscover"
    "defaults"
    "device_groups"
    "devices"
    "enums"
    "object_groups"
    "objects"
    "traps"
-]]

version: 2
project_name: [[ .Env.BINARY ]]

before:
  hooks:
    - gpg --pinentry-mode loopback --import --passphrase {{ .Env.NFPM_PASSPHRASE }} [[ .Env.GPG_KEY_PATH ]]

builds:
  - id: "[[ .Env.BINARY ]]"
    builder: go
    main: ./cmd/validate
    [[/* "whitelabel", rename binary to be flowcoll */]]
    binary: "[[ .Env.BINARY ]]"
    goos:
      - linux
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{ .Version }}
      - -X main.date={{ .Date }}
    env:
      - CGO_ENABLED=0
      - >-
        {{- if eq .Os "linux" }}
          {{- if eq .Arch "arm64"}}CC=aarch64-linux-gnu-gcc{{- end }}
        {{- end }}

checksum:
  algorithm: sha256
  split: true

source:
  prefix_template: "{{ .ProjectName }}-{{ .Version }}/"

nfpms:
  - id: "[[ .Env.PKG_NAME ]]"
    package_name: "[[ .Env.PKG_NAME ]]"
    maintainer: ElastiFlow Inc. <support@elastiflow.com>
    description: "[[ .Env.PKG_DESCRIPTION ]]"
    license: ElastiFlow EULA
    formats:
      - deb
      - rpm
    section: custom
    priority: optional
    bindir: /usr/share/elastiflow/bin
    release: 1
    file_name_template: "[[ .Env.PKG_NAME ]]_{{ .Version }}_{{ .Arch }}"

    contents:
      [[ range $p := ("./" | file.Walk) -]]
      [[ if has $etcDirs (index ($p | strings.Split "/") 0) -]]
      [[ if not ($p | file.IsDir) -]]
      - src: [[ $p ]]
        dst: [[ printf "/etc/elastiflow/snmp/%s" $p ]]
        type: "config|noreplace"
      [[ end -]]
      [[ end -]]
      [[ end ]]

    deb:
     signature:
       key_file: "[[ .Env.GPG_KEY_PATH ]]"

    rpm:
     signature:
       key_file: "[[ .Env.GPG_KEY_PATH ]]"

    overrides:
      deb:
        file_name_template: "[[ .Env.PKG_NAME ]]_{{ .Version }}_{{ .Arch }}"
      rpm:
        file_name_template: "[[ .Env.PKG_NAME ]]-{{ .Version }}-{{ .Release }}.{{ .Arch }}"

signs:
  - id: "[[ .Env.PKG_NAME ]]"
    artifacts: package
    args:
      - --pinentry-mode
      - loopback
      - --output
      - ${artifact}.sig
      - --detach-sig
      - --sign
      - --default-key
      - 6A2E26EFDE24AA7A634A442ED5C0572E5D212F6B
      - --passphrase
      - "{{ .Env.NFPM_PASSPHRASE }}"
      - ${artifact}
