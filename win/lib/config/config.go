package config

import (
    "fmt"
    "path"
    "os"

    log "github.com/jbrodriguez/mlog"
    "gopkg.in/ini.v1"

    _ "../logger"
)

var (
    Version = "0.0.1"

    ServerEndpoint = "https://audit.zhenai.com"

    // shell and rdp exec path
    XShellExec = ""
    SecureCRTExec = ""
    RDPExec = "mstsc"

    // config path
    confPath = path.Join(path.Join(os.Getenv("APPDATA"), "Jumpserver"), "config.ini")

    // check update URL
    CheckUpdateURL = ServerEndpoint + "/api/extension/v1/check-update"
)

func init() {
    serverEndpoint, err := tryToGetFromConf("Default", "Server")
    if err != nil {
        log.Warning("%s", err)
    } else {
        ServerEndpoint = serverEndpoint
    }

    xshellExec, err := tryToGetFromConf("XShell", "Path")
    if err != nil {
        log.Warning("%s", err)
    } else {
        XShellExec = xshellExec
    }

    secureCRT, err := tryToGetFromConf("SecureCRT", "Path")
    if err != nil {
        log.Warning("%s", err)
    } else {
        SecureCRTExec = secureCRT
    }

    checkUpdateURL, err := tryToGetFromConf("Update", "CheckURL")
    if err != nil {
        log.Warning("%s", err)
    } else {
        CheckUpdateURL = checkUpdateURL
    }
}

func tryToGetFromConf(section, key string) (string, error) {
    if _, err := os.Stat(confPath); os.IsNotExist(err) {
        return "", fmt.Errorf("config %s not exist", err)
    }

    conf, err := ini.LooseLoad(confPath)
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    sec, err := conf.GetSection(section)
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    val, err := sec.GetKey(key)
    if err != nil {
        return "", fmt.Errorf("Error while reading config - %s", err)
    }

    return val.String(), nil
}
