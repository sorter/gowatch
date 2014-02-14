gowatch
=======

A basic command line utility that compares the active commit hash of a local git 
repository with its remote counterpart on GitHub, and issues a pull command if 
the remote counterpart reports a differing commit hash at its HEAD.


Setup
-----

1. create a directory to house configuration `mkdir $HOME/.gowatch`
2. generate a [GitHub application OAuth token](https://github.com/settings/applications)
3. place the token in `.gowatch/github_creds`
4. place the GitHub repository details in `.gowatch/repo_map.json`

repo_map.json
`
{ 
    "repo_owner_username": {
        "repo_name": "/local/repo/copy/path"
    }
}
`
