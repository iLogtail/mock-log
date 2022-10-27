# Mock Log

iLogtail mock log tool based on [ilogtail-comparison](https://github.com/EvanLjp/ilogtail-comparison)   
 
## Build

Make binary

```bash
make clean
make
```

Make docker image

```bash
make docker
```

## Usage

All available options
```bash
./bin/mock_log --help

Usage of ./bin/mock_log:
  --item-length int
        value length in json log and nginx log (default 100)
  --key-count int
        key count in json log (default 10)
  --log-err-type string
        nginx java random json (default "random")
  --log-file-count int
        max rotated files (default 10)
  --log-file-size int
        max log size (default 20971520)
  --log-type string
        nginx java random json (default "java")
  --logs-per-sec int
        logs per second upper limit (default 1)
  --path string
        output to file path
  --stderr
        output to stderr
  --stdout
        output to stdout (default true)
  --total-count int
        total log count, set -1 for infinity (default 100)
```

Example usage

```bash
./bin/mock_log --log-type=nginx --path="mock_nginx.log" --stdout=false --total-count=300000 --log-file-size=8850000 --log-file-count=5 --logs-per-sec=5000
```
