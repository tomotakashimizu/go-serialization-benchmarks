# Serializer Performance Benchmark

Redis キャッシュのパフォーマンス改善のためのシリアライザー比較ツール

## 概要

このツールは、JSON→ ドメインモデルの Unmarshal が Redis キャッシュ復元時のボトルネックとなっている問題を解決するため、異なるシリアライゼーション形式のパフォーマンスを比較します。

## 比較対象のシリアライザー

- **JSON** - Go 標準ライブラリ (`encoding/json`)
- **MessagePack** - `github.com/vmihailenco/msgpack/v5`
- **CBOR** - `github.com/fxamacker/cbor/v2`
- **Gob** - Go 標準ライブラリ (`encoding/gob`)

## 測定項目

### 1. シリアライゼーション性能

- Marshal/Unmarshal 速度（平均・中央値）
- データサイズ（バイト数）

### 2. 対称性テスト

- 空スライス/マップの Marshal→Unmarshal 対称性
- nil スライス/マップの Marshal→Unmarshal 対称性

### 3. Redis 性能 (オプション)

- Redis SET/GET 操作の性能測定

## プロジェクト構造

```
benchmarks/serializer/
├── go.mod                          # Go モジュール設定
├── cmd/
│   └── benchmark/
│       └── main.go                 # 実行エントリーポイント
├── internal/
│   ├── models/
│   │   └── test_data.go           # テストデータ構造体
│   ├── serializers/
│   │   ├── serializer.go          # 共通インターフェース
│   │   ├── json.go                # JSON実装
│   │   ├── msgpack.go             # MessagePack実装
│   │   ├── cbor.go                # CBOR実装
│   │   └── gob.go                 # Gob実装
│   ├── benchmark/
│   │   └── runner.go              # ベンチマーク実行ロジック
│   ├── redis/
│   │   └── client.go              # Redis性能測定
│   └── reporter/
│       └── reporter.go            # 結果出力・保存
├── results/                        # 結果出力先
└── README.md                       # このファイル
```

## 使用方法

### 前提条件

- Go 1.24.2 以上
- Redis（Redis 性能測定を行う場合）

### 依存関係のインストール

```bash
cd benchmarks/serializer
go mod tidy
```

### 基本的な使用方法

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

### コマンドラインオプション

| オプション        | デフォルト     | 説明                     |
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

# カスタムRedis設定
go run cmd/benchmark/main.go -redis-addr=192.168.1.100:6379 -redis-password=secret

# 高精度測定（10回測定）
go run cmd/benchmark/main.go -iterations=10
```

## 結果出力

### コンソール出力

実行時に以下の表が表示されます：

1. **シリアライゼーション性能結果**

   - データサイズ（バイト）
   - Marshal/Unmarshal 速度（マイクロ秒）

2. **対称性テスト結果**

   - 空/nil スライス・マップの処理結果

3. **Redis 性能結果**（Redis 測定を行った場合）
   - SET/GET 操作速度（マイクロ秒）

### ファイル出力

`results/` ディレクトリに以下の CSV ファイルが保存されます：

- `serialization_results_YYYYMMDD_HHMMSS.csv` - シリアライゼーション性能
- `symmetry_results_YYYYMMDD_HHMMSS.csv` - 対称性テスト結果
- `redis_results_YYYYMMDD_HHMMSS.csv` - Redis 性能（実行した場合）

## テストデータ構造

約 10 フィールドの 3 層ネスト構造を持つ User モデル：

```go
type User struct {
    ID        int64
    Name      string
    Email     string
    Age       int
    IsActive  bool
    Profile   Profile      // 2層目
    Settings  Settings     // 2層目
    Tags      []string
    Metadata  map[string]interface{}
    CreatedAt time.Time
}

type Profile struct {
    // ... 3層目までのネスト構造
}
```

## Redis 測定について

Redis 測定では以下を行います：

1. シリアライズしたデータの Redis SET 操作
2. Redis GET 操作とデシリアライズ
3. データ整合性チェック
4. 測定完了後のキークリーンアップ

Redis 接続に失敗した場合は警告を表示してスキップします。

## トラブルシューティング

### 依存関係エラー

```bash
# 依存関係を再取得
go mod download
go mod tidy
```

### Redis 接続エラー

```bash
# Redisサーバーの起動確認
redis-cli ping

# Redis測定をスキップ
go run cmd/benchmark/main.go -skip-redis
```

### メモリ不足

大量のテストデータ生成時にメモリ不足が発生する場合：

```bash
# レコード数を減らす
go run cmd/benchmark/main.go -count=10000
```

## パフォーマンス期待値

一般的な傾向：

- **速度**: Gob > MessagePack > CBOR > JSON
- **サイズ**: MessagePack ≈ CBOR < Gob < JSON
- **対称性**: Go 標準ライブラリの方が一貫性が高い傾向

※実際の結果は環境により異なります。
