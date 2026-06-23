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
  <a href="README.pl.md">Polski</a> ·
  <a href="README.th.md">ไทย</a> ·
  <a href="README.he.md">עברית</a> ·
  <a href="README.bn.md">বাংলা</a> ·
  <a href="README.ur.md">اردو</a>
</p>

---

> **جديد على go-skills؟ ابدأ بـ [Project Assessment](skills/workflow/project-assessment.md).** وجّهه إلى مشروع Go الحالي الخاص بك واحصل على درجة نضج إضافةً إلى خارطة طريق مُرتّبة حسب الأولوية حيث يرتبط كل عنصر بنمط go-skills أو الفصل الدقيق المطلوب تطبيقه — للقراءة فقط، لا يُغيّر شيفرتك أبداً.

## المحتوى

<div dir="rtl">

- **المنهجية** — دليل عملي من 18 فصلاً لبناء خدمات Go في بيئة الإنتاج (هيكل الدلائل، البنية متعددة الطبقات، حقن التبعيات اليدوي، الإعداد، إعادة المحاولة، التخزين، transport، المهام، التسجيل، التحقق، الأخطاء، الاختبار، البناء، النشر). اقرأ [الوثيقة الكاملة](skills/methodology/00-canonical-full.md) أو اختر [chapter-skill](METHODOLOGY.md) منفرداً.
- **الأنماط** — أكثر من 52 مدخلاً في 10 فئات: إنشائية، هيكلية، سلوكية، تزامن، مزامنة، رسائل، استقرار، تحليل الأداء، تعابير اصطلاحية، أنماط مضادة. انظر [PATTERNS.md](PATTERNS.md).
- **الأمثلة** — 52 حزمة Go قابلة للتشغيل مع `_test.go`. وحدة مشتركة واحدة ضمن [`examples/`](examples/).
- **سير العمل** — مهارات قابلة للتشغيل قائمة على الفريق: [Project Assessment](skills/workflow/project-assessment.md) (تقييم مشروع قائم)، [Feature Development](skills/workflow/feature-development.md) (بناء ميزة)، [Security Code Review](skills/workflow/security-review.md) (تدقيق + معالجة).
- **الموقع** — GitHub Pages بتصميم Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — عنصر احتياطي لتكرار مستقبلي. انظر [`mcp/`](mcp/).

</div>

## التثبيت كـ Claude Code skills

<div dir="rtl">

أسهل طريقة — الصق هذا في Claude Code وهو يجهّز كل شيء:

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

أو جهّزه يدويًا:

</div>

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
