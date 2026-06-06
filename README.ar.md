<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">أنماط معمارية ومنهجية لخدمات Go — مغلفة كـ Claude Code skills.</p>

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

## المحتوى

<div dir="rtl">

- **المنهجية** — دليل عملي من 18 فصلاً لبناء خدمات Go في بيئة الإنتاج (هيكل الدلائل، البنية متعددة الطبقات، حقن التبعيات اليدوي، الإعداد، إعادة المحاولة، التخزين، transport، المهام، التسجيل، التحقق، الأخطاء، الاختبار، البناء، النشر). اقرأ [الوثيقة الكاملة](skills/methodology/00-canonical-full.md) أو اختر [chapter-skill](METHODOLOGY.md) منفرداً.
- **الأنماط** — أكثر من 52 مدخلاً في 10 فئات: إنشائية، هيكلية، سلوكية، تزامن، مزامنة، رسائل، استقرار، تحليل الأداء، تعابير اصطلاحية، أنماط مضادة. انظر [PATTERNS.md](PATTERNS.md).
- **الأمثلة** — 52 حزمة Go قابلة للتشغيل مع `_test.go`. وحدة مشتركة واحدة ضمن [`examples/`](examples/).
- **الموقع** — GitHub Pages بتصميم Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — عنصر احتياطي لتكرار مستقبلي. انظر [`mcp/`](mcp/).

</div>

## التثبيت كـ Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

<div dir="rtl">

أضف إلى `~/.claude/settings.json`:

</div>

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

<div dir="rtl">

أعد تشغيل Claude Code.

</div>

## تشغيل الأمثلة

```bash
cd examples
go test ./...
```

## الترخيص

<div dir="rtl">

MIT — انظر [LICENSE](LICENSE).

</div>

## المؤلف

<div dir="rtl">

Yevgeniy Achin · <https://github.com/amazopic>

</div>
