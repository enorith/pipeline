## actions define

```js
[
  {
    id: "00",
    action: "http-input",
    inputs: [],
    options: [],
    outputs: [
      {
        type: "httpInput",
        label: "http请求",
      },
    ],
  },
  {
    id: "01",
    action: "data-source",
    name: "database数据源",
    inputs: [
      {
        type: "string",
        label: "数据模型",
      },
      {
        type: "string",
        label: "数据连接",
      },
    ],
    options: [
      {
        type: "string",
        label: "数据模型",
        value: "user",
      },
      {
        type: "string",
        label: "数据连接",
        value: "db DSN",
      },
    ],
    outputs: [
      {
        type: "dataSource",
        label: "数据模型",
      },
    ],
  },
  {
    id: "02",
    action: "query-list",
    inputs: [
      {
        type: "dataSource",
        label: "数据模型",
        // 连接关系
        from: {
          "01": 0,
        },
      },
      {
        type: "httpInput",
        label: "http请求",
        // 连接关系
        from: {
          "00": 0,
        },
      },
    ],
    options: [],
    outputs: [
      {
        type: "list",
        label: "列表",
      },
    ],
  },
];
```

1. action 预排序，尾节点>>>中间节点>>>头节点
2. call(第一节点)

   - 查找 inputs 依赖，nodeId，定位上游节点
   - call(上游节点), [重复步骤 2，递归]
   - outputs，输出，传递至下游节点 inputs[output=outputs:index]

3. 递归结束，返回值

## action compose

## libraries

[Golang JS Engine](https://github.com/dop251/goja)
`github.com/dop251/goja`
