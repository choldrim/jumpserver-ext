package updater

import (
    "fmt"

    "github.com/mozillazg/request"
    "github.com/inconshreveable/go-update"
    Version "github.com/hashicorp/go-version"
    log "github.com/jbrodriguez/mlog"

    "../config"
)


// fetch the binary and replace itself
func doUpdate(URL string) error {
    req := request.NewRequest(nil)
    resp, err := req.Get(URL)
    if err != nil {
        return fmt.Errorf("error while getting update binary - %s", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("error while getting update binary, status code is not correct - status code: %d", resp.StatusCode)
    }

    err = update.Apply(resp.Body, update.Options{})
    if err != nil {
        return fmt.Errorf("error while updating - %s", err)
    }

    return nil
}


// return latest-version, download-url
func checkUpdate() (string, string, error) {
    req := request.NewRequest(nil)
    resp, err := req.Get(config.CheckUpdateURL)
    if err != nil {
        return "", "", fmt.Errorf("error while making a request to check update - %s", err)
    }

    defer resp.Body.Close()

    if !resp.Ok() {
        return "", "", fmt.Errorf("error while checking update with status code: %d", resp.StatusCode)
    }

    data, err := resp.Json()
    if err != nil {
        return "", "", fmt.Errorf("error while reading check update body - %s", err)
    }

    latestVersionData := data.Get("version")
    latestVersion, err := latestVersionData.String()
    if err != nil {
        return "", "", fmt.Errorf("error while parsing check update data - %s", err)
    }

    binaryURLData := data.Get("binary_URL")
    binaryURL, err := binaryURLData.String()
    if err != nil {
        return "", "", fmt.Errorf("error while parsing check update data - %s", err)
    }

    return latestVersion, binaryURL, nil
}

// return true if v1 > v2, else return false
func checkVersion(v1, v2 string) (bool, error) {
    _v1, err := Version.NewVersion(v1)
    if err !=nil {
        return false, err
    }

    _v2, err := Version.NewVersion(v2)
    if err !=nil {
        return false, err
    }

    if _v2.LessThan(_v1) {
        return true, nil
    } else {
        return false, nil
    }
}


func Check() {
    log.Info("checking update...")
    latestVersion, URL, err := checkUpdate()
    if err != nil {
        log.Error(err)
        return
    }

    ret, err := checkVersion(latestVersion, config.Version)
    if err != nil {
        log.Error(fmt.Errorf("error while comparing version - %s", err))
        return
    }

    if ret {
        log.Info("updating app...")
        err = doUpdate(URL)
        if err != nil {
            log.Error(err)
            return
        }

        log.Info("finish app update")
    } else {
        log.Info("new version not found, skip update")
    }
}
