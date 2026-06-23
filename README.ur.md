<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go سروسز کے لیے آرکیٹیکچرل patterns اور میتھڈولوجی — بطور Claude Code skills پیک شدہ۔</p>

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

<div dir="rtl">

> **go-skills پر نئے ہیں؟ [Project Assessment](skills/workflow/project-assessment.md) سے شروع کریں۔** اسے اپنے موجودہ Go پراجیکٹ کی طرف اشارہ کریں اور ایک maturity score کے ساتھ ساتھ ترجیحی روڈمیپ حاصل کریں جہاں ہر آئٹم لاگو کرنے کے لیے عین مطابق go-skills pattern یا فصل سے منسلک ہوتا ہے — صرف-پڑھنے کے لیے، یہ آپ کے کوڈ کو کبھی نہیں بدلتا۔

</div>

## اندر کیا ہے

<div dir="rtl">

- **میتھڈولوجی** — پروڈکشن Go سروسز بنانے کے لیے 18-فصلی پلے بک (directory layout، layered architecture، manual DI، configuration، retries، storage، transport، jobs، logging، validation، errors، testing، build، deploy)۔ [مستند مکمل دستاویز](skills/methodology/00-canonical-full.md) پڑھیں یا کوئی [chapter-skill](METHODOLOGY.md) منتخب کریں۔
- **Patterns** — 10 زمروں میں 52 اندراجات: creational، structural، behavioral، concurrency، synchronization، messaging، stability، profiling، idiom، anti-patterns۔ [PATTERNS.md](PATTERNS.md) دیکھیں۔
- **Examples** — `_test.go` کے ساتھ 52 چلنے والے Go pattern packages۔ [`examples/`](examples/) کے تحت ایک مشترکہ ماڈیول۔
- **Workflows** — چلنے والی ٹیم-بنیاد skills: [Project Assessment](skills/workflow/project-assessment.md) (موجودہ پراجیکٹ کا جائزہ)، [Feature Development](skills/workflow/feature-development.md) (ایک feature بنائیں)، [Security Code Review](skills/workflow/security-review.md) (آڈٹ + اصلاح)۔
- **Site** — Linear-اسٹائل GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP server** — مستقبل کی تکرار کے لیے placeholder۔ [`mcp/`](mcp/) دیکھیں۔

</div>

## بطور Claude Code skills انسٹال کریں

<div dir="rtl">

سب سے آسان طریقہ — اسے Claude Code میں پیسٹ کریں، باقی سب ہو جائے گا:

</div>

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

<div dir="rtl">

یا پھر خود سے سیٹ اپ کریں:

</div>

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

<div dir="rtl">

`~/.claude/settings.json` میں شامل کریں:

</div>

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

<div dir="rtl">

Claude Code دوبارہ شروع کریں۔

</div>

## مثالیں چلائیں

```bash
cd examples
go test ./...
```

## لائسنس

<div dir="rtl">

MIT — [LICENSE](LICENSE) دیکھیں۔

</div>

## مصنف

<div dir="rtl">

Yevgeniy Achin · <https://github.com/amazopic>

</div>
