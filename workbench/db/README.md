# Database Setup for CQRS Pattern

このディレクトリは、CQRSパターンを実装するためのMySQLレプリケーション環境を提供します。

## ディレクトリ構成

```text
workbench/db/
├── README.md              # このファイル
├── Makefile              # タスク実行用のMakefile
├── compose.yaml          # Docker Compose設定
├── command/              # Command（書き込み）用データベース設定
│   ├── my.cnf           # MySQL設定ファイル
│   ├── ddl/             # DDLファイル
│   │   ├── create_object.sql  # オブジェクト作成SQL
│   │   ├── create_record.sql  # レコード作成SQL
│   │   └── master.db          # ダンプファイル
│   ├── scripts/         # 実行スクリプト
│   │   └── dump.sh      # データベースダンプスクリプト
│   └── sql/             # 初期化SQL
│       └── create_repl_user.sql  # レプリケーションユーザー作成
├── query/               # Query（読み込み）用データベース設定
│   ├── my.cnf          # MySQL設定ファイル
│   ├── ddl/            # DDLファイル
│   │   ├── master.db   # レストア用ダンプファイル
│   │   └── replication.sql  # レプリケーション設定SQL
│   └── sql/            # 初期化SQL
└── logs/               # MySQLログファイル
    ├── command/        # Command DBのログ
    └── query/          # Query DBのログ
```

## CQRSレプリケーション構成

### システム全体アーキテクチャ

```mermaid
graph TB
    subgraph "Host Environment"
        subgraph "Application"
            gRPC[gRPC Microservice]
            WriteService[Write Service<br/>書き込み処理]
            ReadService[Read Service<br/>読み込み処理]
        end
        
        subgraph "Docker Compose Environment"
            subgraph "Database Containers"
                CommandContainer[command_db<br/>Container]
                QueryContainer[query_db<br/>Container]
                AdminContainer[db_admin<br/>Container]
            end
            
            subgraph "Persistent Volumes"
                CommandVol[(db-command<br/>Volume)]
                QueryVol[(db-query<br/>Volume)]
            end
            
            subgraph "Host Mounts"
                LogsDir[./logs/<br/>Log Files]
                ConfigDir[./command/<br/>./query/<br/>Config Files]
            end
        end
    end
    
    subgraph "Network: net"
        CommandDB[(Command DB<br/>MySQL 8.0<br/>Port: 3306)]
        QueryDB[(Query DB<br/>MySQL 8.0<br/>Port: 3307)]
        Admin[phpMyAdmin<br/>Port: 3100]
    end
    
    %% アプリケーション層の接続
    gRPC --> WriteService
    gRPC --> ReadService
    
    %% データベース接続
    WriteService -->|Write Operations| CommandDB
    ReadService -->|Read Operations| QueryDB
    
    %% コンテナとデータベースの関係
    CommandContainer --> CommandDB
    QueryContainer --> QueryDB
    AdminContainer --> Admin
    
    %% レプリケーション
    CommandDB -.->|MySQL Replication<br/>Binlog Sync| QueryDB
    
    %% ボリュームマウント
    CommandDB --> CommandVol
    QueryDB --> QueryVol
    
    %% ログとコンフィグマウント
    CommandContainer -.-> LogsDir
    QueryContainer -.-> LogsDir
    CommandContainer -.-> ConfigDir
    QueryContainer -.-> ConfigDir
    
    %% 管理画面接続
    Admin -.-> CommandDB
    Admin -.-> QueryDB
    
    %% スタイル設定
    classDef app fill:#e3f2fd
    classDef write fill:#ffebee
    classDef read fill:#e8f5e8
    classDef db fill:#f3e5f5
    classDef infra fill:#fff3e0
    classDef admin fill:#fafafa
    
    class gRPC,WriteService app
    class WriteService,CommandDB,CommandContainer write
    class ReadService,QueryDB,QueryContainer read
    class CommandVol,QueryVol,LogsDir,ConfigDir infra
    class Admin,AdminContainer admin
```

### アーキテクチャ概要

- **Command DB** (Port: 3306): 書き込み専用のマスターデータベース
- **Query DB** (Port: 3307): 読み込み専用のスレーブデータベース
- **DB Admin** (Port: 3100): phpMyAdmin管理画面

### データフロー

CQRSパターンにおけるデータの流れを以下のMermaid図で示します：

```mermaid
graph TB
    subgraph "Application Layer"
        App[アプリケーション]
        WriteAPI[Write API<br/>書き込み処理]
        ReadAPI[Read API<br/>読み込み処理]
    end
    
    subgraph "Database Layer"
        CommandDB[(Command DB<br/>Port: 3306<br/>書き込み専用)]
        QueryDB[(Query DB<br/>Port: 3307<br/>読み込み専用)]
    end
    
    subgraph "Management"
        Admin[phpMyAdmin<br/>Port: 3100]
    end
    
    %% データフロー
    App --> WriteAPI
    App --> ReadAPI
    
    WriteAPI -->|INSERT/UPDATE/DELETE| CommandDB
    ReadAPI -->|SELECT| QueryDB
    
    CommandDB -.->|MySQL Replication<br/>Binlog同期| QueryDB
    
    Admin -.->|管理画面| CommandDB
    Admin -.->|管理画面| QueryDB
    
    %% スタイル設定
    classDef writeFlow fill:#ffcccc,stroke:#ff0000
    classDef readFlow fill:#ccffcc,stroke:#00ff00
    classDef replication fill:#ccccff,stroke:#0000ff
    classDef management fill:#ffffcc,stroke:#ffaa00
    
    class WriteAPI,CommandDB writeFlow
    class ReadAPI,QueryDB readFlow
    class Admin management
```

#### データフローの詳細

1. **書き込みフロー（赤色）**
   - アプリケーション → Write API → Command DB
   - INSERT、UPDATE、DELETE操作

2. **読み込みフロー（緑色）**
   - アプリケーション → Read API → Query DB  
   - SELECT操作

3. **レプリケーションフロー（青色）**
   - Command DB → Query DB
   - MySQLバイナリログによる自動同期

4. **管理フロー（黄色）**
   - phpMyAdmin経由での両DB管理

### レプリケーション詳細フロー

MySQLレプリケーションの内部動作を詳しく示します：

```mermaid
sequenceDiagram
    participant App as アプリケーション
    participant CDB as Command DB<br/>(Master)
    participant Binlog as Binary Log
    participant QDB as Query DB<br/>(Slave)
    
    Note over App,QDB: 1. 書き込み処理
    App->>CDB: INSERT/UPDATE/DELETE
    CDB->>CDB: データ更新
    CDB->>Binlog: バイナリログに記録
    
    Note over App,QDB: 2. レプリケーション
    QDB->>Binlog: ログイベント要求
    Binlog->>QDB: バイナリログ送信
    QDB->>QDB: ログイベント適用
    
    Note over App,QDB: 3. 読み込み処理
    App->>QDB: SELECT
    QDB->>App: 結果返却
```

## セットアップ手順

### 1. 環境起動

データベース環境を起動します。

```bash
make up
```

### 2. 初期データ作成

Command DBにテーブルとサンプルデータを作成します。

```bash
make create-data
```

### 3. データダンプ

Command DBのデータをダンプファイルに出力します。

```bash
make dump
```

### 4. データリストア

Command DBのダンプをQuery DBにリストアします。

```bash
make restore
```

### 5. レプリケーション開始

Query DBでレプリケーションを開始します。
DumpからGTIDを取得し、`query/ddl/replication.sql`に設定した後、レプリケーションを設定します。

```bash
make start-replication
```

### セットアップフロー図

CQRSレプリケーション環境のセットアップ手順を図解します：

```mermaid
flowchart TD
    Start([開始]) --> Step1[make up<br/>環境起動]
    Step1 --> Check1{Command DB<br/>Query DB<br/>起動確認}
    Check1 -->|OK| Step2[make create-data<br/>初期データ作成]
    Check1 -->|NG| Error1[エラー:<br/>Docker環境確認]
    
    Step2 --> Step3[make dump<br/>データダンプ]
    Step3 --> Step4[make restore<br/>データリストア]
    Step4 --> Step5[make start-replication<br/>レプリケーション開始]
    Step5 --> Check2{レプリケーション<br/>動作確認}
    
    Check2 -->|OK| Success([セットアップ完了])
    Check2 -->|NG| Debug[トラブルシューティング<br/>・ログ確認<br/>・権限確認<br/>・設定確認]
    Debug --> Step4
    
    Error1 --> Start
    
    %% スタイル設定
    classDef startEnd fill:#e1f5fe
    classDef process fill:#f3e5f5
    classDef decision fill:#fff3e0
    classDef error fill:#ffebee
    classDef success fill:#e8f5e8
    
    class Start,Success startEnd
    class Step1,Step2,Step3,Step4,Step5 process
    class Check1,Check2 decision
    class Error1,Debug error
```

## Makefileコマンド一覧

| コマンド | 説明 |
|---------|------|
| `make up` | データベース環境を起動 |
| `make down` | データベース環境を停止 |
| `make dump` | Command DBをダンプ |
| `make restore` | Query DBにダンプをリストア |
| `make start-replication` | レプリケーションを開始 |
| `make create-data` | テストデータを作成 |
| `make help` | 利用可能なコマンドを表示 |

## アクセス情報

### データベース接続

- **Command DB**: `localhost:3306`
- **Query DB**: `localhost:3307`
- **ユーザー**: `root`
- **パスワード**: `password`

### 管理画面

- **phpMyAdmin**: <http://localhost:3100>

## 注意事項

1. **レプリケーション順序**: 必ず上記の手順1-5の順序で実行してください
2. **データ整合性**: Command DBでのデータ変更は自動的にQuery DBに反映されます
3. **ログ確認**: レプリケーション状況は`logs/`ディレクトリで確認できます
4. **初期化**: 環境をリセットする場合は`make down`後に`make up`から再実行してください

## トラブルシューティング

### レプリケーションが動作しない場合

1. Command DBとQuery DBが正常に起動しているか確認
2. レプリケーションユーザーが正しく作成されているか確認
3. ログファイルでエラーメッセージを確認

### 接続エラーの場合

1. ポートが正しく開放されているか確認
2. Docker Composeのhealthcheckが通っているか確認
