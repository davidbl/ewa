package commands

func CheckErr(err error, msg string) {
  if err != nil {
    config.Log.Printf("%s: %v\n", msg, err)
  }
}
func CheckErrFatal(err error, msg string) {
  if err != nil {
    config.Log.Fatalf("%s: %v\n", msg, err)
  }
}
