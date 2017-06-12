package session

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "path"
    "strconv"
    "time"
    "strings"

    log "github.com/jbrodriguez/mlog"

    "../parameter"
    "../config"
    "../encrypter"
)

type SessionJson struct {
    Result string           `json:"result"`
    AssetName string        `json:"asset_name"`
    Password string         `json:"password"`  // for rdp
    Err bool                `json:"err"`
}

func ReadySessionFile(params *parameter.InputParams) (string, error) {
    sessionDir, err := ioutil.TempDir("", "jms-")
    if err != nil {
        return "", err
    }

    //defer os.RemoveAll(sessionDir)

    sessionFile, err := downloadSessionFile(params.Token, params.AssetID, params.ShellType, sessionDir)

    return sessionFile, err
}

func downloadSessionFile(token string, assetID int, shellType string, sessionDir string) (string, error) {
    respJson := SessionJson{}
    err := getJson(token, assetID, shellType, &respJson)
    if err != nil {
        return "", fmt.Errorf("Error while get asset(id: %d) session file - %s", assetID, err)
    }

    if respJson.Err {
        return "", fmt.Errorf(respJson.Result)
    }

    assetName := respJson.AssetName
    if len(assetName) == 0 {
        assetName = path.Base(sessionDir)
    }

    suffix := getSuffix(shellType)
    sessionFile := path.Join(sessionDir, assetName) + suffix
    if shellType == "RDP" {
        enc, err := encrypter.EncryptPWD(respJson.Password)
        if err != nil {
            log.Error(err)
        }

        respJson.Result = strings.Replace(respJson.Result, "__PASSWORD__", enc, -1)
    }

    err = ioutil.WriteFile(sessionFile, []byte(respJson.Result), 0644)
    if err != nil {
        return "", fmt.Errorf("Error while writing %s - %s", sessionFile, err)
    }

    return sessionFile, nil
}

func getSuffix(shellType string) string {
    switch shellType {
        case "XShell":
            return ".xsh"
        case "RDP":
            return ".rdp"
        default:
            return ""
    }
}

func getJson(token string, assetID int, shellType string, target *SessionJson) error {
    url := config.ServerEndpoint + "/api/shell/v1/session-file/" + strconv.Itoa(assetID) + "/"
    client := &http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return fmt.Errorf("Error while requesting %s - %s", url, err)
    }

    req.Header.Set("Authorization", "Bearer " + token)

    q := req.URL.Query()
    q.Add("shell_type", shellType)
    req.URL.RawQuery = q.Encode()

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Error while downloading %s - %s", url, err)
    }
    defer resp.Body.Close()

    if ! strings.HasPrefix(resp.Status, "2") {
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return err
        }
        return fmt.Errorf(string(b))
    }

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    err = json.Unmarshal(b, &target)
    return err
}
