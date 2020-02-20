# Vault Backup utility

Just little tool which can recursively copy content in HashiCorp Vault Storage and put it in YAML manifest.
Now it can only backup and restore KVs and policies only.

## Configuration

This container use command line arguments for configuration (See golang flag lib). Nothing difficult :)

### Usage

```bash
  -addr       string
        Address of vault service.
  -authtype   string
        Method of token asquition. "token" read VAULT_TOKEN env var. "kubejwt" use Kubernetes vault plugin and default way.
  -token      string
        Token for work with vault service.
  -jwtpath    string
        Custom path to Kubernetes ServiceAccount JWT file. Useless for any non-"kubejwt" auth methods. Default is "/var/run/secrets/kubernetes.io/serviceaccount/token".
  -jwtrole    string
        Custom vault role to assume via Kubernetes ServiceAccount JWT. Useless for any non-"kubejwt" auth methods.
  [ command ] string
        What to do: only "backup" and "restore" available.
        "backup"  - copy all policyes and specified KVs. Receive zero or more KV names to backup. See examples.
        "restore" - restore from backup file. Receive only one filename via args. See examples.
  [ kv ]      string
        Name of KV to backup. See examples.
  [ file ]    string
        Name of backup file to restore. See examples.
```

### Examples

```bash
  vault-backup -addr="http://127.0.0.1:8200" -token=s.111111111111111111111111 -authtype=token backup KV_1 KV_2 # backup policies and KVs "KV_1" and "KV_2"
  vault-backup -addr="http://127.0.0.1:8200" -token=s.111111111111111111111111 -authtype=token restore my/backup/file.yml # restore from file
  vault-backup -addr="http://127.0.0.1:8200" backup # auth via Kubernetes SA JWT and backup policies only
```
