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

    "../parameter"
    "../config"
)

type SessionJson struct {
    Result string           `json:"result"`
    AssetName string        `json:"asset_name"`
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

    sessionFile := path.Join(sessionDir, assetName) + ".xsh"
    err = ioutil.WriteFile(sessionFile, []byte(respJson.Result), 0644)
    if err != nil {
        return "", fmt.Errorf("Error while writing %s - %s", sessionFile, err)
    }

    return sessionFile, nil
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
