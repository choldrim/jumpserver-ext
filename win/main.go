package main

import (
    "time"

    _ "./lib/logger"
    "./lib/parameter"
    "./lib/session"
    "./lib/launcher"

    log "github.com/jbrodriguez/mlog"
)


var (
    Version = "0.0.1"
)


func work() (error) {
    log.Info("decoding params")
    params, err := parameter.CheckParams()
    if err != nil {
        return err
    }

    log.Info("Downloading session file")
    sessionFile, err := session.ReadySessionFile(params)
    if err != nil {
        return err
    }

    log.Info("Launching shell")
    err = launcher.LaunchShell(params.ShellType, sessionFile)
    if err != nil {
        return err
    }

    return nil
}


func main() {
    err := work()
    if err != nil {
        log.Error(err)
        time.Sleep(time.Second * 2)
    }
}
