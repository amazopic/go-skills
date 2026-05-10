<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go サービスのアーキテクチャパターンと方法論 — Claude Code skills としてパッケージ化。</p>

<p align="center">
  <a href="README.md">English</a> ·
  <a href="README.ru.md">Русский</a> ·
  <a href="README.fr.md">Français</a> ·
  <a href="README.de.md">Deutsch</a> ·
  <a href="README.uk.md">Українська</a> ·
  <a href="README.sl.md">Slovenščina</a> ·
  <a href="README.it.md">Italiano</a> ·
  <a href="README.es.md">Español</a> ·
  <a href="README.zh.md">中文</a> ·
  <a href="README.ja.md">日本語</a> ·
  <a href="README.ko.md">한국어</a> ·
  <a href="README.ar.md">العربية</a>
</p>

---

## 内容

- **方法論** — 本番 Go サービス構築のための 18 章の実践ガイド（ディレクトリ構成、レイヤードアーキテクチャ、手動 DI、設定、リトライ、ストレージ、transport、ジョブ、ログ、バリデーション、エラー、テスト、ビルド、デプロイ）。[完全な標準ドキュメント](skills/methodology/00-canonical-full.md)を読むか、[chapter-skill](METHODOLOGY.md) を個別に選択してください。
- **パターン** — 9 カテゴリにわたる 50 以上のエントリ：生成型、構造型、振る舞い型、並行性、同期、メッセージング、安定性、プロファイリング、イディオム、アンチパターン。[PATTERNS.md](PATTERNS.md) を参照。
- **サンプル** — `_test.go` 付きの 23 個の実行可能な Go パターンパッケージ。[`examples/`](examples/) 以下に共有モジュールが 1 つ。
- **サイト** — Linear スタイルの GitHub Pages：<https://amazopic.github.io/go-skills/>
- **MCP server** — 将来のイテレーション向けのプレースホルダー。[`mcp/`](mcp/) を参照。

## Claude Code skills としてインストール

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

`~/.claude/settings.json` に追加：

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Claude Code を再起動してください。

## サンプルの実行

```bash
cd examples
go test ./...
```

## ライセンス

MIT — [LICENSE](LICENSE) を参照。

## 作者

Yevgeniy Achin · <https://github.com/amazopic>
