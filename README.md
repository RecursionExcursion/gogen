## ⚠️ Requirement: Precompiled Standard Library for Cross-Platform Targets

This tool generates code/scripts intended to run on multiple platforms (Windows, macOS, Linux).  
To enable this, the Go toolchain must have the standard library (`std`) installed for each target.

Before using this tool, you **must run**:

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go install std
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go install std
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go install std
