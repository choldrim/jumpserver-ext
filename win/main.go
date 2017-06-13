package main

import (
    "time"
    "sync"

    log "github.com/jbrodriguez/mlog"

    _ "./lib/logger"
    "./lib/config"
    "./lib/launcher"
    "./lib/parameter"
    "./lib/session"
    "./lib/updater"
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
    config.Version = Version
    log.Info("Start up, current version is %s", Version)

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        updater.Check()
        defer wg.Done()
    }()

    err := work()
    if err != nil {
        log.Error(err)
        time.Sleep(time.Second * 2)
    }

    wg.Wait()
}
