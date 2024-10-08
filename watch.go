/*
    watch.go
    Sayed Khader sy@sayedkhader.com
    
    when invoked, external file is read that maps a remote github repository
    to a local active directory. If the remote commit hash mismatches the hash
    in the local repository, perform a pull to the active directory

    external file format is a valid json:
    { 
        "repo_name": "local_directory",
        . . . 
    }
*/


package main

import "os"
import "os/exec"
import "fmt"
import "bytes"
import "encoding/json"
import "io/ioutil"
import "net/http"
import _ "crypto/sha256"

type RepoResponse struct {
    Commit struct {
        Sha string
    }
}

func main() {

    GITHUB_HOST := "https://api.github.com/"

    home := os.Getenv("HOME")
    confdir := home + "/.gowatch/"
    repoFilePath := confdir + "repo_map.json"
    fileBuffer, err := ioutil.ReadFile(repoFilePath)

    branch := os.Getenv("ENV")
    if branch == "prod" {
        branch = "master"
    }

    token_path := confdir + "github_creds"
    gitToken, tokenErr := ioutil.ReadFile(token_path)

    if tokenErr != nil {
        panic(tokenErr) // problem reading github OAuth token
    }
    gitToken = gitToken[:len(gitToken)-1]

    // perform git auth
    if err != nil {
        panic(err) // problem creating auth request
    }
    client := &http.Client{}

    if err != nil {
        panic(err)
    }
    //var githubConfig map[string]interface{}
    var githubConfig map[string]map[string]string
    err = json.Unmarshal(fileBuffer, &githubConfig)

    if err != nil {
        panic(err)
    }

    for gitUser, repoMap := range githubConfig {
        for repoName, activeDir := range repoMap {
            // determine the local commit hash
            masterHashPath := activeDir +"/.git/refs/heads/" + branch
            masterHashBuffer, masterError := ioutil.ReadFile(masterHashPath)
            var localHash string
            if masterError == nil {
                localHash = string(masterHashBuffer)
                localHash = localHash[:len(localHash)-1]
            } else {
                panic(masterError)
            }

            repoPath := "repos/"+gitUser+"/"+repoName+"/branches/" + branch
            repoUrl := GITHUB_HOST + repoPath
            //fmt.Printf("repo url: %s\n", repoUrl)

            branchReq, err := http.NewRequest("GET", repoUrl, nil)
            if err != nil {
                panic(err)
            }
            branchReq.Header.Add("Authorization","token "+string(gitToken))
            resp, err := client.Do(branchReq)
            if err != nil {
                panic(err)
            }
            respBody, err := ioutil.ReadAll(resp.Body)
            // determine remote commit hash, perform auth
            rr := &RepoResponse{}
            err = json.Unmarshal(respBody, &rr)
            if err != nil {
                panic(err)
            }
            if localHash != rr.Commit.Sha {
                // issue a pull inside the local active directory
                err = os.Chdir(activeDir)
                fmt.Printf("chdr(%s)\n", activeDir)
                if err != nil {
                    panic(err)
                }
                var out bytes.Buffer
                pullCmd := exec.Command("ssh-agent", "/bin/sh", "-c", "ssh-add "+ home+"/.ssh/gowatch_id_rsa; git pull origin " + branch)
                fmt.Println(pullCmd)
                pullCmd.Stdout = &out
                err = pullCmd.Run()
                if err != nil {
                    fmt.Println(err)
                    panic(err) // err running git pull
                }
                fmt.Println(out.String())
            }
        }
    }

}
