# aivoke

KISS CLI tool for using LLMs and managing prompts

```bash
# do not pipe to Bash at home
./aivoke bash "find all Golang files in the current directory and show me the individual and total line counts"  | bash
15 ./internal/util/id.go
10 ./internal/util/assert.go
36 ./internal/util/fs.go
14 ./internal/util/env.go
47 ./internal/openai/openai.go
39 ./internal/cli/cli.go
8 ./internal/data/data.go
9 ./cmd/main/aivoke/main.go
Total lines: 178
```
