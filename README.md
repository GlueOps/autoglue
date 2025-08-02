# AutoGlue

## Generate Docs
```shell
make swagger
```

## Generate Binary
```shell
make build
```
## Routes
See the swagger interface at http://localhost:8080/swagger

## Config Priorities
| Source            | Priority |
|-------------------|----------|
| CLI Flags         | Highest  |
| ENV vars / `.env` | Higher   |
| `config.yaml`     | Medium   |
| Defaults          | Lowest   |


## Config Loading
| Source      | Viper Handles It?                         | Example                   |
| ----------- | ----------------------------------------- | ------------------------- |
| Defaults    | ✅ `viper.SetDefault()`                    | `"8080"`                  |
| config.yaml | ✅ auto-loaded                             | `bind_port: 8080`         |
| .env        | ✅ via `godotenv` + `viper.AutomaticEnv()` | `AUTOGLUE_BIND_PORT=8081` |
| CLI Flags   | ✅ with `viper.BindPFlag()`                | `--bind-port 8082`        |

## Configs
### config.yaml
```yaml
auth:
  secret: 5lCzHfAyGbwXfdDI59IoPaIyN_q0cmOjglE8xh0XXpo=
bind_address: 0.0.0.0
bind_port: "8080"
database:
    dsn: postgres://user:pass@localhost:5432/autoglue?sslmode=disable

```

### .env
```shell
# Used for docker-compose
DB_CLIENT=postgresql
DB_USER=autoglue
DB_PASSWORD=autoglue
DB_HOST=localhost
DB_PORT=5432
DB_NAME=autoglue

# Used by Autoglue
AUTOGLUE_DATABASE_DSN="postgres://autoglue:autoglue@localhost:5432/autoglue"
AUTOGLUE_BIND_ADDRESS=127.0.0.1
AUTOGLUE_BIND_PORT=9090
AUTOGLUE_AUTH_SECRET=kNAoX/QKZ14ewmYJ5MHRLDPVK7uuPNmHqP5vJ/WJ8NM=
```