A Collection of deepflow wasm plugins

| plugin | desc |
| --- | --- |
| header-extract | extract all headers transport by http |

## usage

- compile
```bash
# compile
cd header-extract
tinygo build -o header.wasm -target wasi -gc=precise -panic=trap -scheduler=none -no-debug .
# upload plugin
deepflow-ctl plugin create --type wasm --image header.wasm --name header-extract-plugin
```

- update deepflow-agent config
```yaml
static_config:
  wasm-plugins:
  - header-extract-plugin
```