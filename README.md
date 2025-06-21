# Go Serialization Benchmarks

Go 言語における各種シリアライゼーション形式のパフォーマンス比較ツールです。

## 概要

このプロジェクトは、異なるシリアライゼーション形式のパフォーマンスを測定・比較するためのベンチマークツールです。特に以下の用途に適しています：

- アプリケーションに最適なシリアライゼーション形式の選択
- データ永続化や API 通信における性能最適化
- Redis キャッシュなどでのシリアライゼーション性能評価

## 対応シリアライザー

- **JSON** - Go 標準ライブラリ ([`encoding/json`](https://pkg.go.dev/encoding/json))
- **CBOR** - [`github.com/fxamacker/cbor/v2`](https://github.com/fxamacker/cbor)
- **EasyJSON** - 高性能 JSON with コード生成 ([`github.com/mailru/easyjson`](https://github.com/mailru/easyjson)) - コード生成による高性能の JSON シリアライザー
- **Gob** - Go 標準ライブラリ ([`encoding/gob`](https://pkg.go.dev/encoding/gob))
- **GoJSON** - 高性能 JSON ([`github.com/goccy/go-json`](https://github.com/goccy/go-json)) - 標準ライブラリの 100%互換高性能版
- **JSONiter** - 高性能 JSON ([`github.com/json-iterator/go`](https://github.com/json-iterator/go)) - 標準ライブラリの 100%互換高性能版
- **Msgp** - 高性能 MessagePack with コード生成 ([`github.com/tinylib/msgp`](https://github.com/tinylib/msgp)) - コード生成による高性能の MessagePack シリアライザー
- **MsgPack** - [`github.com/vmihailenco/msgpack/v5`](https://github.com/vmihailenco/msgpack)
- **Protobuf** - Google Protocol Buffers ([`google.golang.org/protobuf`](https://pkg.go.dev/google.golang.org/protobuf#section-readme)) - 効率的で言語に依存しないシリアライゼーション形式

## 測定項目

### 1. シリアライゼーション性能

- Marshal/Unmarshal 速度（平均値・中央値）
- シリアライズ後のデータサイズ

### 2. Marshal/Unmarshal の対称性テスト

- 空スライス/マップの Marshal→Unmarshal 対称性
- nil スライス/マップの Marshal→Unmarshal 対称性

### 3. Redis 性能測定（オプション）

- Redis SET/GET 操作の性能測定
- 実際のキャッシュ使用シナリオでの評価

## プロジェクト構造

```
go-serialization-benchmarks/
├── cmd/
│   └── benchmark/
│       └── main.go                 # 実行エントリーポイント
├── internal/
│   ├── benchmark/
│   │   └── runner.go              # ベンチマーク実行ロジック
│   ├── models/
│   │   └── test_data.go           # テストデータ構造体
│   ├── redis/
│   │   └── client.go              # Redis性能測定
│   ├── reporter/
│   │   └── reporter.go            # 結果出力・保存
│   └── serializers/
│       ├── serializer.go          # 共通インターフェース
│       ├── json.go                # JSON実装
│       ├── cbor.go                # CBOR実装
│       ├── easyjson.go            # EasyJSON実装
│       ├── gob.go                 # Gob実装
│       ├── gojson.go              # GoJSON実装
│       ├── jsoniter.go            # JSONiter実装
│       ├── msgp.go                # Msgp実装
│       ├── msgpack.go             # MsgPack実装
│       └── protobuf.go            # Protobuf実装
├── proto/
│   ├── user.proto                  # Protocol Buffersスキーマ定義
│   └── user.pb.go                  # 生成されたProtocol Buffersコード
├── results/                        # 結果出力先
├── go.mod                          # Go モジュール設定
└── README.md                       # このファイル
```

## インストール

```bash
git clone https://github.com/tomotakashimizu/go-serialization-benchmarks.git
cd go-serialization-benchmarks
go mod tidy
```

## 使用方法

### 基本実行

```bash
# デフォルト設定で実行（10万件データ、5回測定）
go run cmd/benchmark/main.go

# レコード数を指定して実行
go run cmd/benchmark/main.go -count=10000

# Redis測定をスキップ
go run cmd/benchmark/main.go -skip-redis

# ヘルプ表示
go run cmd/benchmark/main.go -help
```

### コマンドライン引数

| 引数              | デフォルト     | 説明                     |
| ----------------- | -------------- | ------------------------ |
| `-count`          | 100000         | 生成するテストレコード数 |
| `-iterations`     | 5              | ベンチマーク測定回数     |
| `-redis-addr`     | localhost:6379 | Redis サーバーアドレス   |
| `-redis-password` | ""             | Redis パスワード         |
| `-redis-db`       | 0              | Redis データベース番号   |
| `-output`         | ./results      | 結果出力ディレクトリ     |
| `-skip-redis`     | false          | Redis 測定をスキップ     |
| `-help`           | false          | ヘルプ表示               |

### 実行例

```bash
# 小規模テスト（1万件、Redis無し）
go run cmd/benchmark/main.go -count=10000 -skip-redis

# カスタムRedis設定での実行
go run cmd/benchmark/main.go -redis-addr=192.168.1.100:6379 -redis-password=secret

# 10回測定
go run cmd/benchmark/main.go -iterations=10
```

## テストデータ

4 層ネスト構造を持つ User モデルを使用：

```go
type User struct {
    ID        int64
    Name      string
    Email     string
    Age       int
    IsActive  bool
    Profile   Profile                // 2層目
    Settings  Settings               // 2層目
    Tags      []string
    Metadata  map[string]interface{}
    CreatedAt time.Time
}

type Profile struct {
    FirstName   string
    LastName    string
    Bio         string
    Avatar      string
    SocialLinks []Link             // 3層目
    Preferences Preferences        // 3層目
}

type Preferences struct {
    Theme         string
    Language      string
    Notifications map[string]bool
    Privacy       PrivacySettings  // 4層目
}
```

複雑なネスト構造により、実際のアプリケーションデータに近い条件でのベンチマークが可能です。

## 結果出力

### コンソール出力

実行時に以下の情報がテーブル形式で表示されます：

1. **シリアライゼーション性能結果**

   - データサイズ（バイト）
   - Marshal/Unmarshal 速度

2. **Marshal/Unmarshal の対称性テスト結果**

   - 空/nil スライス・マップの処理結果

3. **Redis 性能結果**（Redis 測定を行った場合）
   - SET/GET 操作速度

### ファイル出力

`results/` ディレクトリに以下の CSV ファイルが保存されます：

- `serialization_results_YYYYMMDD_HHMMSS.csv` - シリアライゼーション性能
- `symmetry_results_YYYYMMDD_HHMMSS.csv` - Marshal/Unmarshal の対称性テスト結果
- `redis_results_YYYYMMDD_HHMMSS.csv` - Redis 性能（実行した場合）
