A Collection of deepflow wasm plugins

| plugin | desc |
| --- | --- |
| header-extract | extract all headers transport by http |

## usage

- compile
```bash
# compile
cd header-extract
tinygo build -o header.wasm -target wasi -gc=custom -panic=trap -tags=custommalloc -scheduler=none -no-debug .
# upload plugin
deepflow-ctl plugin create --type wasm --image header.wasm --name header-extract-plugin
```

- update deepflow-agent config
```yaml
wasm_plugins:
- header-extract-plugin
```