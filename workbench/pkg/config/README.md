# Config パッケージ

このパッケージは、プロジェクト全体で使用される設定管理の実装を提供します。

## 概要

[Viper](https://github.com/spf13/viper)を使用した設定管理機能を提供し、TOMLファイルと環境変数の両方をサポートします。

## 機能

- **TOMLファイルサポート**: 構造化された設定ファイル
- **環境変数オーバーライド**: 本番環境での柔軟な設定変更
- **型安全なアクセス**: Viperの型安全なゲッターメソッド
- **依存性注入統合**: Uber Fxとの完全な統合
- **設定の階層化**: ネストされた設定構造のサポート

## 使用方法

### 初期化

```go
import (
    "yourproject/pkg/config"
)

// 設定を読み込む
cfg, err := config.NewConfig(".", "config")
if err != nil {
    log.Fatal(err)
}

// 設定値の取得
logLevel := cfg.GetString("log.level")
dbHost := cfg.GetString("mysql.host")
dbPort := cfg.GetInt("mysql.port")
```

### パラメータ

- **configPath**: 設定ファイルの検索パス（例: ".", "./config"）
- **configName**: 設定ファイル名（拡張子なし、例: "config"）

## 設定ファイル

### 基本構造

デフォルトの設定ファイル `config.toml`:

```toml
[log]
level = "info"
format = "text"

[mysql]
dbname = "mydb"
host = "localhost"
port = 3306
user = "dbuser"
password = "dbpass"
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = "1h"
```

### セクション

#### ログ設定 (`[log]`)

| キー | 型 | デフォルト | 説明 |
|------|------|-----------|------|
| level | string | "info" | ログレベル (debug/info/warn/error) |
| format | string | "text" | ログフォーマット (text/json) |

#### データベース設定 (`[mysql]`)

| キー | 型 | デフォルト | 説明 |
|------|------|-----------|------|
| dbname | string | - | データベース名 |
| host | string | "localhost" | データベースホスト |
| port | int | 3306 | データベースポート |
| user | string | - | ユーザー名 |
| password | string | - | パスワード |
| max_idle_conns | int | 10 | 最大アイドル接続数 |
| max_open_conns | int | 100 | 最大オープン接続数 |
| conn_max_lifetime | string | "1h" | 接続の最大ライフタイム |

## 環境変数

環境変数を使用して設定を上書きできます。

### 命名規則

環境変数名は、TOMLの階層構造を`_`（アンダースコア）で区切り、大文字に変換します：

```text
TOML: log.level
環境変数: LOG_LEVEL

TOML: mysql.host
環境変数: DB_MYSQL_HOST  # データベース設定には DB_ プレフィックス
```

### ログ設定の環境変数

```bash
export LOG_LEVEL=debug
export LOG_FORMAT=json
```

### データベース設定の環境変数

データベース設定には `DB_` プレフィックスを使用：

```bash
export DB_MYSQL_DBNAME=production_db
export DB_MYSQL_HOST=db.example.com
export DB_MYSQL_PORT=3306
export DB_MYSQL_USER=prod_user
export DB_MYSQL_PASSWORD=secure_password
export DB_MYSQL_MAX_IDLE_CONNS=20
export DB_MYSQL_MAX_OPEN_CONNS=200
export DB_MYSQL_CONN_MAX_LIFETIME=2h
```

### 環境変数の優先順位

設定値は以下の優先順位で決定されます（上が優先）：

1. 環境変数
2. 設定ファイル（config.toml）
3. デフォルト値（コード内）

## 設定値の取得

### 基本的な取得

```go
// 文字列
logLevel := cfg.GetString("log.level")

// 整数
dbPort := cfg.GetInt("mysql.port")

// 真偽値
debug := cfg.GetBool("app.debug")

// 期間
timeout := cfg.GetDuration("server.timeout")
```

### デフォルト値の設定

存在しない設定にデフォルト値を設定：

```go
cfg.SetDefault("server.port", 8080)
cfg.SetDefault("app.debug", false)
```

### ネストされた設定の取得

```go
// サブツリーの取得
mysqlConfig := cfg.Sub("mysql")
if mysqlConfig != nil {
    host := mysqlConfig.GetString("host")
    port := mysqlConfig.GetInt("port")
}
```

### 構造体へのマッピング

```go
type Config struct {
    Log struct {
        Level  string `mapstructure:"level"`
        Format string `mapstructure:"format"`
    } `mapstructure:"log"`
    
    MySQL struct {
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        DBName   string `mapstructure:"dbname"`
        User     string `mapstructure:"user"`
        Password string `mapstructure:"password"`
    } `mapstructure:"mysql"`
}

var config Config
if err := cfg.Unmarshal(&config); err != nil {
    log.Fatal(err)
}
```

## Uber Fxとの統合

### モジュール定義

```go
// pkg/config/module.go
var Module = fx.Module("config",
    fx.Provide(
        fx.Annotate(
            NewConfig,
            fx.ParamTags(`name:"configPath"`, `name:"configName"`),
        ),
    ),
)
```

### 使用例

```go
func main() {
    fx.New(
        // 設定パラメータを提供
        fx.Supply(
            fx.Annotate(".", fx.ResultTags(`name:"configPath"`)),
            fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
        ),
        
        // 設定モジュールを追加
        config.Module,
        
        // 設定を使用するサービス
        fx.Provide(NewMyService),
        
        fx.Invoke(func(*MyService) {}),
    ).Run()
}

type MyService struct {
    cfg *viper.Viper
}

func NewMyService(cfg *viper.Viper) *MyService {
    return &MyService{cfg: cfg}
}
```

## ベストプラクティス

### 1. 機密情報は環境変数で

パスワードやAPIキーなどの機密情報は、設定ファイルではなく環境変数で管理：

```bash
# Good: 環境変数で機密情報を管理
export DB_MYSQL_PASSWORD=secure_password

# Bad: 設定ファイルに平文で記載
# [mysql]
# password = "secure_password"  # Gitにコミットしない！
```

### 2. 設定ファイルの階層化

関連する設定をグループ化：

```toml
[log]
level = "info"
format = "text"

[server]
host = "0.0.0.0"
port = 8080
timeout = "30s"

[database]
# データベース設定
```

### 3. 型安全なアクセス

適切な型のゲッターメソッドを使用：

```go
// Good
port := cfg.GetInt("server.port")

// Avoid
port := cfg.Get("server.port").(int)  // パニックの可能性
```

### 4. 設定のバリデーション

起動時に必須の設定を検証：

```go
func validateConfig(cfg *viper.Viper) error {
    required := []string{"mysql.host", "mysql.dbname", "mysql.user"}
    for _, key := range required {
        if !cfg.IsSet(key) {
            return fmt.Errorf("required config %s is not set", key)
        }
    }
    return nil
}
```

### 5. 開発環境と本番環境の分離

環境ごとに設定ファイルを用意：

```bash
config/
  ├── config.toml              # 共通設定
  ├── config.development.toml  # 開発環境
  └── config.production.toml   # 本番環境
```

## トラブルシューティング

### 設定ファイルが見つからない

1. ファイルパスを確認：

   ```go
   cfg, err := config.NewConfig(".", "config")  // ./config.toml を探す
   ```

2. ファイルの存在を確認：

   ```bash
   ls -la config.toml
   ```

3. 作業ディレクトリを確認：

   ```go
   pwd, _ := os.Getwd()
   fmt.Println("Current directory:", pwd)
   ```

### 環境変数が反映されない

1. 環境変数名を確認：

   ```bash
   # TOML: mysql.host → 環境変数: DB_MYSQL_HOST
   echo $DB_MYSQL_HOST
   ```

2. プレフィックスを確認：
   - ログ設定: プレフィックスなし (`LOG_LEVEL`)
   - データベース設定: `DB_` プレフィックス (`DB_MYSQL_HOST`)

3. 環境変数の設定を確認：

   ```bash
   env | grep DB_
   env | grep LOG_
   ```

### 型変換エラー

1. 設定値の型を確認：

   ```toml
   port = 3306      # 整数
   host = "localhost"  # 文字列
   debug = true     # 真偽値
   ```

2. 適切なゲッターを使用：

   ```go
   port := cfg.GetInt("mysql.port")      // OK
   port := cfg.GetString("mysql.port")   // "3306" (文字列)
   ```

## Docker環境での使用

### docker-compose.yml

```yaml
version: '3.8'
services:
  app:
    build: .
    environment:
      - LOG_LEVEL=info
      - LOG_FORMAT=json
      - DB_MYSQL_HOST=db
      - DB_MYSQL_DBNAME=mydb
      - DB_MYSQL_USER=dbuser
      - DB_MYSQL_PASSWORD=dbpass
    depends_on:
      - db
  
  db:
    image: mysql:8.0
    environment:
      - MYSQL_DATABASE=mydb
      - MYSQL_USER=dbuser
      - MYSQL_PASSWORD=dbpass
      - MYSQL_ROOT_PASSWORD=rootpass
```

### Dockerfile

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# 設定ファイルをコピー
COPY config.toml .

# アプリケーションをコピー
COPY . .

# ビルド
RUN go build -o server ./cmd/server

# 実行
CMD ["./server"]
```

## 参考リンク

- [Viper公式ドキュメント](https://github.com/spf13/viper)
- [12 Factor App - Config](https://12factor.net/config)
- [Environment Variables Best Practices](https://12factor.net/config)
