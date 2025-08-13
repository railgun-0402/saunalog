# ✅ Saunalog のコンセプト

## 🎯 サービス概要
- ユーザー自身が訪れたサウナの「体験」をログに残す
- 共有・検索できる投稿型Webサービス

## 要件
| 機能            | 概要             | 優先度 |
| ------------- | -------------- | --- |
| 🔍 サウナ施設一覧・詳細 | 名前・住所・温度・水風呂など | ◎   |
| ➕ サウナ施設登録     | ユーザーが新規施設を投稿可能 | ◎   |
| 📝 体験ログ投稿     | 混雑度・体調・「整い度」など | ◎   |
| 🗓️ マイページ     | 自分の体験ログ一覧表示    | ◯   |
| 👍 ログへのいいね    | 共感できる体験にリアクション | ◯   |
| 🔐 認証（仮）      | ログイン機能（匿名でも可）  | △   |
| 📊 統計（仮）      | 最も「整った」施設など表示  | △   |

## 🔧 技術スタック（予定）
| レイヤー   | 技術候補                                 |
| ------ | ------------------------------------ |
| フロント   | Next.js（React） |
| バックエンド | **Go + Echo（Clean Architecture）**    |
| DB     | DynamoDB               |
| 認証     | Cognito                   |
| インフラ   | AWS + Terraform                      |
| 監視     | NewRelic（APM、ログ、外形監視）                |

## Labmdaへデプロイ(最終的にはTerraformでIaC化)
```
cd lambda/compress_image

GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip function.zip bootstrap
```

### ・experience logs

```Request```

```
curl -X POST http://localhost:8080/logs \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u123",
    "facility_id": "f456",
    "date": "2025-08-02",
    "congestion_level": 3,
    "totonoi_level": 5,
    "cost_performance": 4,
    "service_quality": 5,
    "comment": "キンキン水風呂で昇天しました"
  }'
```

```Response```

```
{
    "id": 1
}
```
