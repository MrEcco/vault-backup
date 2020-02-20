package main

// ItemHandler struct
type ItemHandler struct {
	Path    string                 `yaml:"path"`
	Content map[string]interface{} `yaml:"content"`
}

// PolicyHandler struct
type PolicyHandler struct {
	Name string `yaml:"name"`
	ACL  string `yaml:"acl"`
}

// BackupHandler struct
type BackupHandler struct {
	KVBuckets []struct {
		Name    string        `yaml:"name"`
		Entries []ItemHandler `yaml:"entries"`
	} `yaml:"kvbuckets"`
	Policies []PolicyHandler `yaml:"policies"`
}
