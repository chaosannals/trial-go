# wasm

注：go 的安装目录下 misc\wasm 提供 wasm_exec.js 文件。
不同版本的 go 的 wasm_exec.js 不兼容其他版本 go 编译的 wasm。

GO 到 JS 的类型对照

```
| Go                     | JavaScript             |
| ---------------------- | ---------------------- |
| js.Value               | [its value]            |
| js.Func                | function               |
| nil                    | null                   |
| bool                   | boolean                |
| integers and floats    | number                 |
| string                 | string                 |
| []interface{}          | new array              |
| map[string]interface{} | new object             |
```
