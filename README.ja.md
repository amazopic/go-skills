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
  <a href="README.ar.md">العربية</a> ·
  <a href="README.pt-BR.md">Português (BR)</a> ·
  <a href="README.tr.md">Türkçe</a> ·
  <a href="README.id.md">Bahasa Indonesia</a> ·
  <a href="README.vi.md">Tiếng Việt</a> ·
  <a href="README.hi.md">हिन्दी</a> ·
  <a href="README.zh-TW.md">繁體中文</a> ·
  <a href="README.pl.md">Polski</a> ·
  <a href="README.th.md">ไทย</a> ·
  <a href="README.he.md">עברית</a> ·
  <a href="README.bn.md">বাংলা</a> ·
  <a href="README.ur.md">اردو</a>
</p>

---

> **go-skills は初めてですか？[Project Assessment](skills/workflow/project-assessment.md) から始めましょう。** 既存の Go プロジェクトに向けるだけで、成熟度スコアと優先順位付けされたロードマップが得られます。各項目には、適用すべき go-skills のパターンや章への正確なリンクが付きます — 読み取り専用で、コードを一切変更しません。

## 内容

- **方法論** — 本番 Go サービス構築のための 18 章の実践ガイド（ディレクトリ構成、レイヤードアーキテクチャ、手動 DI、設定、リトライ、ストレージ、transport、ジョブ、ログ、バリデーション、エラー、テスト、ビルド、デプロイ）。[完全な標準ドキュメント](skills/methodology/00-canonical-full.md)を読むか、[chapter-skill](METHODOLOGY.md) を個別に選択してください。
- **パターン** — 10 カテゴリにわたる 52 以上のエントリ：生成型、構造型、振る舞い型、並行性、同期、メッセージング、安定性、プロファイリング、イディオム、アンチパターン。[PATTERNS.md](PATTERNS.md) を参照。
- **サンプル** — `_test.go` 付きの 52 個の実行可能な Go パターンパッケージ。[`examples/`](examples/) 以下に共有モジュールが 1 つ。
- **ワークフロー** — チームベースで実行可能な skills：[Project Assessment](skills/workflow/project-assessment.md)（既存プロジェクトを評価）、[Feature Development](skills/workflow/feature-development.md)（機能を構築）、[Security Code Review](skills/workflow/security-review.md)（監査と修正）。
- **サイト** — Linear スタイルの GitHub Pages：<https://amazopic.github.io/go-skills/>
- **MCP server** — 将来のイテレーション向けのプレースホルダー。[`mcp/`](mcp/) を参照。

## Claude Code skills としてインストール

いちばん簡単 — これを Claude Code に貼るだけで全部セットアップ：

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

または手動でセットアップ：

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
