# syugo
Collect files from git repositorys.

## Setup
```
$ go get -u github.com/nissy/syugo/cmd/syugo
```

## Config
syugo.sample.toml
```toml
[[collects]]
  repository = "https://github.com/google/fonts"
  requests = [
    "ufl/ubuntu/Ubuntu-Medium.ttf",
    "ufl/ubuntu/Ubuntu-MediumItalic.ttf",
  ]
  version = "master"
  dir = "fonts"

[[collects]]
  repository = "https://github.com/Microsoft/fonts"
  requests = [
    "Symbols/winjs-symbols.ttf",
  ]
  version = "master"
  dir = "fonts"
```

## Usage
```
$ syugo -c syugo.sample.toml
$ tree fonts/
fonts/
├── Symbols
│   └── winjs-symbols.ttf
└── ufl
    └── ubuntu
        ├── Ubuntu-Medium.ttf
        └── Ubuntu-MediumItalic.ttf
```
