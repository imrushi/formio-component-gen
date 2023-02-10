# FormIo Component Generator

This repo will help you to generate the formio components. but you will need to add the entry in `index.js` and `builder.js` inside formio repo `src/component/`.

repo: [https://github.com/formio/formio.js/](https://github.com/formio/formio.js/)

## Flags

- `location` : Location of the component folder or where want the component to be generated.

- `json` : Json file with below formate

```
{
    "Lable": {
        "data":"value"
    },
    "test":{...}
}
```

- `group` : Form Io component gourp in which component you want

## Build

`go build -o component-gen *.go`

To check flags

`./component-gen -h`
