<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go 服务的架构模式与方法论 — 打包为 Claude Code skills。</p>

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

> **初次使用 go-skills？请从 [Project Assessment](skills/workflow/project-assessment.md) 开始。** 让它分析你现有的 Go 项目，即可获得一份成熟度评分以及一份按优先级排序的路线图，其中每一项都链接到要应用的具体 go-skills 模式或章节——只读，它绝不会改动你的代码。

## 内容

- **方法论** — 一份 18 章的实用指南，用于构建生产级 Go 服务（目录结构、分层架构、手动 DI、配置、重试、存储、transport、任务、日志、验证、错误、测试、构建、部署）。阅读[完整规范文档](skills/methodology/00-canonical-full.md)或选择单独的 [chapter-skill](METHODOLOGY.md)。
- **模式** — 10 个类别中的 52 条目：创建型、结构型、行为型、并发、同步、消息传递、稳定性、性能分析、惯用写法、反模式。参见 [PATTERNS.md](PATTERNS.md)。
- **示例** — 52 个可运行的 Go 模式包，包含 `_test.go`。在 [`examples/`](examples/) 下有一个共享模块。
- **工作流** — 可运行的基于团队的 skills：[Project Assessment](skills/workflow/project-assessment.md)（评估现有项目）、[Feature Development](skills/workflow/feature-development.md)（构建功能）、[Security Code Review](skills/workflow/security-review.md)（审计 + 修复）。
- **网站** — Linear 风格的 GitHub Pages：<https://amazopic.github.io/go-skills/>
- **MCP server** — 未来迭代的占位实现。参见 [`mcp/`](mcp/)。

## 安装为 Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

添加到 `~/.claude/settings.json`：

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

重启 Claude Code。

## 运行示例

```bash
cd examples
go test ./...
```

## 许可证

MIT — 参见 [LICENSE](LICENSE)。

## 作者

Yevgeniy Achin · <https://github.com/amazopic>
