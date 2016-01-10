package envar

import (
  "os"
)

type decideStr func(string,string) string
type decideInt func(string, int) int

func String(varName string, defaultVal string) string {
  v := os.Getenv(varName)
  if v != "" {
    return v
  } else {
    return defaultVal
  }
}

func ByteSlice(varName string, defaultVal string) []byte {
  v := os.Getenv(varName)
  if v != "" {
    return []byte(v)
  } else {
    return []byte(defaultVal)
  }
}

func IntFunc(varName string, f decideInt, defaultVal int) int {
  v := os.Getenv(varName)
  return f(v, defaultVal)
}

func StringFunc(varName string, f decideStr, defaultVal string) string {
  v := os.Getenv(varName)
  return f(v, defaultVal)
}
