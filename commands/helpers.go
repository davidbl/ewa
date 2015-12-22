package commands
import "fmt"

func CheckErr(err error, msg string) {
  if err != nil {
    panic(fmt.Sprintf("%s: %v", msg, err))
  }
}
