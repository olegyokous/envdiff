# envdiff

> Compare `.env` files across environments and surface missing or mismatched keys with structured output for CI pipelines.

---

## Installation

```bash
go install github.com/yourname/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envdiff.git && cd envdiff && go build -o envdiff .
```

---

## Usage

```bash
envdiff --base .env.example --compare .env.production
```

**Example output:**

```
MISSING KEYS (in production, not in example):
  - DATABASE_POOL_SIZE

EXTRA KEYS (in example, not in production):
  + REDIS_URL
  + FEATURE_FLAG_NEW_UI

MISMATCHED (key present in both, value type differs):
  ~ LOG_LEVEL  [expected: "info"]  [got: ""]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Base `.env` file to compare against | `.env.example` |
| `--compare` | Target `.env` file to validate | `.env` |
| `--format` | Output format: `text`, `json` | `text` |
| `--strict` | Exit with code 1 if any diff is found | `false` |

### CI Integration

```yaml
- name: Validate env
  run: envdiff --base .env.example --compare .env.production --strict --format json
```

---

## License

MIT © [yourname](https://github.com/yourname)