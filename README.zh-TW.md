<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go 服務的架構模式與方法論 — 封裝成 Claude Code skills。</p>

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
  <a href="README.pl.md">Polski</a>
</p>

---

## 內容概覽

- **方法論** — 一套用於建構正式 Go 服務的 18 章實戰手冊（目錄結構、分層架構、手動 DI、設定、重試、儲存、傳輸、工作排程、日誌、驗證、錯誤處理、測試、建置、部署）。閱讀[權威完整文件](skills/methodology/00-canonical-full.md)，或挑選一個[章節 skill](METHODOLOGY.md)。
- **模式** — 橫跨 10 大類別的 52 多條條目：建立型、結構型、行為型、並行、同步、訊息傳遞、穩定性、效能分析、慣用法、反模式。詳見 [PATTERNS.md](PATTERNS.md)。
- **範例** — 52 個可執行的 Go 模式套件，附 `_test.go`。位於 [`examples/`](examples/) 之下的一個共用模組。
- **網站** — Linear 風格的 GitHub Pages：<https://amazopic.github.io/go-skills/>
- **MCP server** — 為未來迭代預留的佔位。詳見 [`mcp/`](mcp/)。

## 安裝為 Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

加入 `~/.claude/settings.json`：

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

重新啟動 Claude Code。

## 執行範例

```bash
cd examples
go test ./...
```

## 授權條款

MIT — 詳見 [LICENSE](LICENSE)。

## 作者

Yevgeniy Achin · <https://github.com/amazopic>
