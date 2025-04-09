package app

type Logger interface {
    Printf(format string, v ...interface{})
    Println(v ...interface{})
    Fatalf(format string, v ...interface{})
}