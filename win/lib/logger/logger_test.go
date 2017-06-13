package logger

import (
    "fmt"
    . "testing"

    "github.com/jbrodriguez/mlog"
)

func TestCheckUpdate(t *T) {
    mlog.Info("hello logger: info")
    mlog.Warning("hello logger: info")
    mlog.Error(fmt.Errorf("hello logger: error"))
}
