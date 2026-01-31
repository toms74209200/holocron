# 書庫アプリ 要件・Tech Stack

## 要件

### 機能要件
1. **ユーザー登録**
   - Firebase Anonymous認証後、ユーザー名を設定
   - ユーザー名は変更可能

2. **書籍登録**
   - ISBN入力で書籍情報自動取得
   - バーコードスキャン対応

3. **書籍一覧・検索**
   - 貸出可能/貸出中のステータス表示
   - 貸出中の場合は利用者名を表示
   - タイトル・著者で検索

4. **貸出・返却**
   - バーコードスキャンで貸出（貸出者名を記録）
   - バーコードスキャンで返却

### 非機能要件
- イベントソーシング
- CQRS（Command/Query分離）
- 結果整合性（トランザクションロック回避）
- アプリ内goroutineでイベント処理→読み取りモデル更新

## Tech Stack

### フロントエンド
- Next.js (App Router)
- Firebase Auth (Anonymous Login)
- バーコードスキャンライブラリ

### バックエンド
- Go (標準 `net/http`)
- OpenAPI Generator (`oapi-codegen`)
- sqlc
- Firebase Admin SDK（UID検証用）

### データベース
- SQLite (`:memory:` or 一時ファイル)
- 本番移行時はPostgreSQL + pg_cron

### テスト
- Small test: Property-Based Testing (`gopter`)
- Medium test: Testcontainers
- Large test: 外部システムに接続する場合のテスト
- Web API test: Testcontainers

### アーキテクチャ
- Presentation層: OpenAPI自動生成
- Orchestration層: Command/Query分離
- Logic層: 純粋関数
- Data層: EventStore + ReadModel

## 画面仕様

### 画面遷移

```
[アプリ起動]
    │
    ↓ (自動匿名認証 + ユーザー自動生成「ゲストXXXX」)
    │
    └─→ [ホーム画面] ←→ [書籍登録画面]
```

### 画面一覧

| 画面 | 用途 |
|------|------|
| ホーム（書籍一覧） | 書籍の一覧表示、検索 |
| 書籍登録 | ISBN入力/バーコードスキャンで書籍追加 |

### ユーザー名の自動生成

- サーバー側でユーザー初回アクセス時に `ゲスト{ランダム4桁}` を自動設定
- または Firebase UIDの末尾4桁を使う
