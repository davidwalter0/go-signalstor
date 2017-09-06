---
This is test software experimenting with:

- xml - json conversion using an xml library
- ftp transfer tests

Todo 

- database i/o using persist library
- encryption

The goal is to demonstrate transformation and secure backup of some
source data

---

```
go get github.com/davidwalter0/xml2json

cd ${GOPATH}/src/github.com/davidwalter0/xml2json/examples
```

- fetch a version of a signal backup xml file 

- edit an environment file, export the corresponding environment
  variables or set the corresponding command line flags

```
. ftp.environment
export FTP_FILENAME=/path/to/SignalPlaintextBackup.xml
go run signalxmlparse.go
```
