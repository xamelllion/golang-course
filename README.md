TestApp. CLI tool that provide you information about github repo

```
Usage: TestApp owner/repo
```

## Steps for bulding

### 1. Build tool
```bash
go build .
```

### 2. Run it
```bash
./TestApp Xamelllion/TestApp
```
output:
```
TestApp:
        description: Homework for GoLang course 2026
        start count: 0
        forks count: 0
        date creation: 2026-03-03T16:45:04Z
```

## Install it to the system

### 0. Add go/bin folder to the PATH
```bash
  echo "export PATH=\$PATH:~/go/bin" >> ~/.bashrc
```

### 1. Install package
```bash
go install .
```

### 2. Run it
```bash
TestApp Xamelllion/TestApp
````