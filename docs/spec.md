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
    └─→ [ホーム] ─┬→ [書籍登録]
                  ├→ [書籍詳細] → 貸出/返却
                  └→ [マイページ]
```

### ホーム（書籍一覧）
- 書籍の一覧表示
- 貸出可能/貸出中のステータス表示
- 貸出中の場合は利用者名を表示
- タイトル・著者で検索
- 書籍タップで詳細へ遷移

### 書籍登録
- バーコードスキャンで登録
- ISBN手動入力で登録
- 書籍情報を手動入力で登録（ISBNがない本）

### 書籍詳細
- 書籍情報の表示
- 貸出状況の表示
- 貸出ボタン / 返却ボタン
- （将来）貸出期間の設定

### マイページ
- ユーザー名の変更
- 借りている本の一覧・管理

### ユーザー名の自動生成

- サーバー側でユーザー初回アクセス時に `ゲスト{ランダム4桁}` を自動設定
- または Firebase UIDの末尾4桁を使う
