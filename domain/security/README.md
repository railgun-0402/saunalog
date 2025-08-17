# Saunalog Auth & Security README

Go/Echo ベースの Saunalog における認証・認可とセキュリティ設計。
Web（SPA/SSR）・モバイル/WebView どちらも視野に入れて設計する。

## アーキテクチャ方針
### 短命 Access トークン + 長命 Refresh トークン を採用
- Access: JWT（5–15 分）… API 認可用。サーバー側で stateless 検証
- Refresh: ランダム ID（JTI）を DB で管理（7–30 日）… ローテーション＆失効可能

### Web: Cookie（HttpOnly/Secure/SameSite）で配布
### モバイル/ネイティブ: ```Authorization: Bearer``` ヘッダーで配布。

> Cookie セッション方式も選べるが、モバイル/サービス横断を考慮し JWT 方式を第一候補

## 脅威モデル（ざっくりベース）

- 資格情報詰め合わせ（クレデンシャルスタッフィング）
- XSS によるトークン窃取
- CSRF による不正リクエスト
- token再利用（リプレイ攻撃）
- DB 流出（パスワード・Refresh JTI）
- 権限昇格/認可不備
- ブルートフォース/DoS

## パスワード取り扱い

- ハッシュ: bcrypt（Cost 10–12）または Argon2id（推奨）
- 検証: CompareHashAndPassword
- ポリシー: 最低 8–12 文字、一般的な漏えいパスワード拒否（Pwned Passwords 照合は任意）
- Pepper（任意）: アプリ固有の秘密文字列を連結してハッシュ（KMS/環境変数で管理）
- 将来更新: コスト引き上げ検知（NeedsRehash）を用意する予定

## JWT 設計
### 署名方式: 
HS256（共有鍵）で開始、ローテーションを見据えて kid 利用。将来的に RS256/ECDSA へ移行も

### Access Claims（例）
- iss（発行者）: "saunalog"
- sub（ユーザーID）: 文字列（例: "12345")
- aud（対象）: ["api"]（任意）
- iat, nbf, exp: 発行/有効化/失効
  - （nbf = now - 30s で時計ズレ吸収）
- jti: 一意 ID（newJTI()）
- カスタム: kind: "access", role（必要最小限。PII は入れない）

### Refresh Claims（例）

- kind: "refresh"
- ```jti``` は DB に保存し、デバイス/UA/IP/期限を紐付け
- ローテーション必須（使うたびに新発行・旧無効化）

### JTI 生成（暗号学的乱数）

```go
func NewJTI() (string, error) {
    var b [16]byte
    if _, err := rand.Read(b[:]); err != nil { return "", err }
    return hex.EncodeToString(b[:]), nil // 32文字hex
}
```

## トークン保管と配送

- Web: Cookie
  - HttpOnly, Secure, SameSite=Lax(or Strict), Path=/, 適切な Domain
  - CSRF 対策: Cookie 配布時は SameSite + CSRF トークン（状態変更系のみ）
- ネイティブ/SPA: 
  - メモリ or セキュアストレージに保持、Authorization: Bearer <access>
  - LocalStorage は XSS に弱かった気がするので使用しない予定

## 代表エンドポイント (予定)

### ```POST /auth/login```
- 入力: email, password
- 処理: ユーザー検索 → bcrypt で照合 → Access/Refresh 発行 → Cookie/ヘッダーに設定
  
### ```POST /auth/refresh```
- 入力: Refresh（Cookie/ヘッダー）
- 処理: Refresh 検証 → JTI を DB で照合 → ローテーション（新 Refresh/Access 発行、旧失効）
- 再利用検知: 旧 Refresh を再度使われたら全セッション失効

### ```POST /auth/logout```
- 入力: Refresh
- 処理: JTI を失効（DB 更新）、Cookie 削除

### ```GET /me```
- 入力: Access
- 処理: JWT 検証 → ユーザー情報返却

## DB スキーマ（例）
```sql
CREATE TABLE users (
    id           BIGINT PRIMARY KEY AUTO_INCREMENT,
    email        VARCHAR(255) NOT NULL UNIQUE,
    password     VARCHAR(72)  NOT NULL, -- bcryptは60文字。余裕を持たせる
    name         VARCHAR(100) NOT NULL,
    gender       VARCHAR(16)  NOT NULL,
    age          INT          NOT NULL,
    prefecture   VARCHAR(64)  NOT NULL,
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    jti          CHAR(32)     PRIMARY KEY,      -- NewJTI()
    user_id      BIGINT       NOT NULL,
    issued_at    DATETIME     NOT NULL,
    expires_at   DATETIME     NOT NULL,
    revoked_at   DATETIME     NULL,
    user_agent   VARCHAR(255) NULL,
    ip_address   VARBINARY(16) NULL,            -- IPv4/6 を保存したければ
    CONSTRAINT fk_refresh_user FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_refresh_user (user_id),
    INDEX idx_refresh_expires (expires_at)
);
```
---
### Echo 実装ひな形
### 生成（Access/Refresh）
```go
type TokenKind string
const (
    Issuer = "saunalog"
    KindAccess  TokenKind = "access"
    KindRefresh TokenKind = "refresh"
)
type Claims struct {
    Kind TokenKind `json:"kind"`
    Role string    `json:"role,omitempty"`
    jwt.RegisteredClaims
}
var secret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(userID uint64, role string, ttl time.Duration, kind TokenKind) (string, time.Time, error) {
    now := time.Now()
    exp := now.Add(ttl)
    jti, _ := NewJTI()
    sub := strconv.FormatUint(userID, 10)
    
    c := &Claims{
        Kind: kind,
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{
        Issuer:    Issuer,
        Subject:   sub,
        IssuedAt:  jwt.NewNumericDate(now),
        NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)),
        ExpiresAt: jwt.NewNumericDate(exp),
        ID:        jti,
        },
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
    s, err := t.SignedString(secret)
    return s, exp, err
}
```

### ミドルウェア（Access 検証）
```go
func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
    tokenStr := extractFromCookieOrHeader(c) // 実装は環境に応じて
    if tokenStr == "" {
        return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
    }
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        return secret, nil
    }, jwt.WithIssuer(Issuer), jwt.WithValidMethods([]string{"HS256"}))
    if err != nil || !token.Valid {
        return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
    }
    claims := token.Claims.(*Claims)
    if claims.Kind != KindAccess {
        return echo.NewHTTPError(http.StatusUnauthorized, "invalid token kind")
    }
    c.Set("userID", claims.Subject)
    c.Set("role", claims.Role)
    return next(c)
    }
}

```

### Refresh ローテーション（概略）
```go
// 1) 受け取ったRefreshをJWT検証 → jti取得
// 2) DBのrefresh_tokensで jti が未失効 && 期限内かチェック
// 3) 新しいRefresh(新jti)を発行・保存、旧jtiをrevoked_atで失効
// 4) 新しいAccessも発行し直す
// 5) 再利用検知：既に失効済みjtiが使われたらそのユーザーの全jtiを失効
```

## CSRF / XSS / ヘッダー

- CSRF: Cookie 利用時は SameSite=Lax(or Strict) + CSRF トークン（状態変更リクエストのみ必須）
- XSS: CSP（default-src 'self' を基本）、dangerouslySetInnerHTML 禁止、入力値エスケープ
- セキュリティヘッダー: Strict-Transport-Security, X-Content-Type-Options, X-Frame-Options/frame-ancestors, Referrer-Policy

---

## レート制限・ロックアウト

- ログイン API に IP/ユーザーID 単位のレート制限（例: 5/min）
- 連続失敗で一時ロック（指数バックオフ）

---

## 権限管理（認可）

- MVP: RBAC（role: admin|user）から開始
- 余力があればスコープ（read:jobs, write:jobs など）を Access に付与
- サーバー側でも最終チェック（トークンだけに依存しない）
---
## ログ / 監査
- ログイン成功/失敗、Refresh ローテーション、ログアウト、失効イベントを構造化ログで記録
- 可能なら監査テーブルに保存（だれが・いつ・どこから）

---

## 秘密情報と鍵管理

- JWT_SECRET_KEY は .env ではなく Secret Manager/KMS を推奨（少なくとも環境変数）
- 鍵ローテーション: JWT Header の kid を使い、検証側は複数鍵に対応
- すべて TLS（Secure Cookie 必須）

### 設定項目（ENV 例）
```ini
JWT_SECRET_KEY=...
ACCESS_TTL_MINUTES=15
REFRESH_TTL_DAYS=14
BCRYPT_COST=12
COOKIE_DOMAIN=.example.com
COOKIE_SECURE=true
COOKIE_SAMESITE=Lax
```
---

## テスト指針

- 単体: パスワードハッシュ/比較、JWT 生成/検証、JTI 生成の衝突なし
- 結合: /auth/login→/auth/refresh→/me→/auth/logout の一連
- 悪性ケース: 期限切れ、未来発行（nbf）、改ざん（署名不一致）、再利用（旧 Refresh）

## クイックチェックリスト (AI作成)

- パスワードは平文保存せず、bcrypt/Argon2id でハッシュ
- Access は短命、Refresh は DB 保持＆ローテーション
- JWT に PII を入れない、iss/sub/exp/iat/nbf/jti 設定
- Cookie は HttpOnly/Secure/SameSite、CSRF 対策あり
- ミドルウェアで JWT 検証（発行者/アルゴリズム/種別/期限）
- 鍵は Secret Manager/KMS、kid でローテーション可能
- ログ・失効・再利用検知の仕組みがある
- レート制限・ロックアウト・MFA（任意）
