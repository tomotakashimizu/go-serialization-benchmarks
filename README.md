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
- **FlatBuffers** - ゼロコピーシリアライゼーション ([`github.com/google/flatbuffers`](https://github.com/google/flatbuffers)) - メモリ効率に優れたクロスプラットフォームシリアライゼーション形式
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

```plaintext
go-serialization-benchmarks/
├── cmd/
│   └── benchmark/
│       └── main.go                 # 実行エントリーポイント
├── internal/
│   ├── benchmark/
│   │   └── runner.go              # ベンチマーク実行ロジック
│   ├── models/
│   │   └── test_data.go           # テストデータ構造体
│   ├── flatbuffers/
│   │   ├── user.fbs               # FlatBuffersスキーマ定義
│   │   └── generated/             # FlatBuffers生成コード
│   ├── proto/
│   │   ├── user.proto             # Protocol Buffersスキーマ定義
│   │   └── user.pb.go             # 生成されたProtocol Buffersコード
│   ├── redis/
│   │   └── client.go              # Redis性能測定
│   ├── reporter/
│   │   └── reporter.go            # 結果出力・保存
│   └── serializers/
│       ├── serializer.go          # 共通インターフェース
│       ├── json.go                # JSON実装
│       ├── cbor.go                # CBOR実装
│       ├── easyjson.go            # EasyJSON実装
│       ├── flatbuffers.go         # FlatBuffers実装
│       ├── gob.go                 # Gob実装
│       ├── gojson.go              # GoJSON実装
│       ├── jsoniter.go            # JSONiter実装
│       ├── msgp.go                # Msgp実装
│       ├── msgpack.go             # MsgPack実装
│       └── protobuf.go            # Protobuf実装
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

#### コマンド実行例

```bash
go run cmd/benchmark/main.go -count=10000 -skip-redis
```

```bash
Serializer Performance Benchmark
=================================
Test data count: 10000
Benchmark iterations: 5
Output directory: ./results
Redis: localhost:6379 (skip: true)

========================================================================================================================
SERIALIZATION BENCHMARK RESULTS
========================================================================================================================
Serializer   | Data Size    | Marshal Avg  | Marshal Med  | Unmarshal Avg | Unmarshal Med
             | (MB)         | (ms)         | (ms)         | (ms)         | (ms)        
------------------------------------------------------------------------------------------------------------------------
JSON         | 9.06         | 139.28       | 137.27       | 218.02       | 215.54      
CBOR         | 6.87         | 56.14        | 56.50        | 131.92       | 132.44      
EasyJSON     | 9.06         | 37.16        | 35.80        | 92.34        | 87.53       
FlatBuffers  | 9.27         | 94.89        | 93.25        | 49.40        | 47.49       
Gob          | 4.72         | 62.01        | 61.74        | 88.61        | 88.28       
GoJSON       | 9.06         | 116.68       | 121.31       | 122.82       | 130.12      
JSONiter     | 9.06         | 42.17        | 36.15        | 147.17       | 139.42      
Msgp         | 6.94         | 25.53        | 25.08        | 38.18        | 38.39       
MsgPack      | 6.94         | 77.75        | 74.31        | 165.28       | 147.20      
Protobuf     | 4.97         | 114.01       | 111.78       | 134.50       | 136.57      
========================================================================================================================

====================================================================================================
STRICT TYPE PRESERVATION TEST RESULTS
====================================================================================================
Serializer   | Empty→Empty  | Empty{}→{}   | Nil→Nil      | Nil→Nil     
             | (Slices)     | (Maps)       | (Slices)     | (Maps)      
----------------------------------------------------------------------------------------------------
JSON         | ✓            | ✓            | ✓            | ✓           
CBOR         | ✓            | ✓            | ✓            | ✓           
EasyJSON     | ✓            | ✓            | ✓            | ✓           
FlatBuffers  | ✓            | ✓            | ✗            | ✗           
Gob          | ✗            | ✓            | ✓            | ✓           
GoJSON       | ✓            | ✓            | ✓            | ✓           
JSONiter     | ✓            | ✓            | ✓            | ✓           
Msgp         | ✗            | ✓            | ✓            | ✗           
MsgPack      | ✓            | ✓            | ✓            | ✓           
Protobuf     | ✗            | ✗            | ✓            | ✗           
====================================================================================================

Benchmark completed successfully!
Results saved to: ./results
```

#### 出力内容

1. **シリアライゼーション性能結果**
   - データサイズ（MB）
   - Marshal/Unmarshal 速度（平均・中央値）

2. **Marshal/Unmarshal の対称性テスト結果**
   - 空/nil スライス・マップの型保持確認
   - ✓: 厳密な型保持、✗: 型変換あり

3. **Redis 性能結果**（Redis 測定を行った場合）
   - SET/GET 操作速度

### ファイル出力

`results/` ディレクトリに以下の CSV ファイルが保存されます：

- `serialization_results_YYYYMMDD_HHMMSS.csv` - シリアライゼーション性能
- `symmetry_results_YYYYMMDD_HHMMSS.csv` - Marshal/Unmarshal の対称性テスト結果
- `redis_results_YYYYMMDD_HHMMSS.csv` - Redis 性能（実行した場合）
