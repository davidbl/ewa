package main

import (
  "fmt"
  "os"
  "./commands"
)

func main() {
  if err := commands.EwaCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(-1)
  }

}
