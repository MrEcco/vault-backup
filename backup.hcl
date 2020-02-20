# Special policy for backing up. Use root token for restore from backups.
path "mykv/*" {
  capabilities = ["read", "list"]
}
path "mykv" {
  capabilities = ["read", "list"]
}
path "sys/policies" {
  capabilities = ["read", "list"]
}
path "sys/policies/*" {
  capabilities = ["read", "list"]
}
