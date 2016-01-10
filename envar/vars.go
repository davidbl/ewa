package envar

import (
  "bytes"
  "fmt"
  "os"
)

type EnvVar struct {
  key string
  _default string
  usage string
  value string
}

func (e EnvVar) String() string {
  return fmt.Sprintf("%-25s\t%-25q\t%-15q\t%-15s", e.key, e.value, e._default, e.usage)
}

type decideStr func(string,string) string

var vars map[string]EnvVar

func init(){
  vars = make(map[string]EnvVar)
}

func Add(key string, defaultV string, usage string, fn decideStr) {
  var v string
  if fn != nil {
    v = stringFunc(key, fn, defaultV)
  } else {
    v = _string(key, defaultV)
  }
  _envVar := EnvVar{key: key, _default: defaultV, usage: usage, value: v}
  vars[key] = _envVar
}

func Get(key string) string {
  return vars[key].value
}

func _string(varName string, defaultVal string) string {
  v := os.Getenv(varName)
  if v != "" {
    return v
  } else {
    return defaultVal
  }
}

func stringFunc(varName string, f decideStr, defaultVal string) string {
  v := os.Getenv(varName)
  return f(v, defaultVal)
}

func Help() string {
  var buffer  bytes.Buffer
  buffer.WriteString(fmt.Sprintf("%-25s\t%-25s\t%-15s\t%-15s\n", "Key","Current Value","Default Value","Description"))
  for _, k := range vars {
    buffer.WriteString(fmt.Sprintf("%v\n",k))
  }
  return buffer.String()
}
