<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">דפוסים ארכיטקטוניים ומתודולוגיה לשירותי Go — ארוזים כ-Claude Code skills.</p>

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

> **חדשים ב-go-skills? התחילו עם [Project Assessment](skills/workflow/project-assessment.md).** כוונו אותו לפרויקט ה-Go הקיים שלכם וקבלו ציון בשלות בתוספת מפת דרכים מתועדפת שבה כל פריט מקושר לדפוס go-skills או לפרק המדויק שיש ליישם — לקריאה בלבד, אינו משנה לעולם את הקוד שלכם.

</div>

## מה בפנים

<div dir="rtl">

- **מתודולוגיה** — ספר משחקים בן 18 פרקים לבניית שירותי Go בפרודקשן (פריסת ספריות, ארכיטקטורה שכבתית, DI ידני, הגדרות, ניסיונות חוזרים, אחסון, transport, משימות, לוגים, ולידציה, שגיאות, בדיקות, בנייה, פריסה). קראו את [המסמך הקנוני המלא](skills/methodology/00-canonical-full.md) או בחרו [chapter-skill](METHODOLOGY.md).
- **דפוסים** — 52 ערכים על פני 10 קטגוריות: יצירתיים, מבניים, התנהגותיים, מקביליות, סנכרון, מסרים, יציבות, פרופיילינג, ניבים, אנטי-דפוסים. ראו [PATTERNS.md](PATTERNS.md).
- **דוגמאות** — 52 חבילות דפוס Go הניתנות להרצה עם `_test.go`. מודול משותף אחד תחת [`examples/`](examples/).
- **Workflows** — skills מבוססי-צוות הניתנים להרצה: [Project Assessment](skills/workflow/project-assessment.md) (הערכת פרויקט קיים), [Feature Development](skills/workflow/feature-development.md) (בניית פיצ'ר), [Security Code Review](skills/workflow/security-review.md) (ביקורת + תיקון).
- **אתר** — GitHub Pages בסגנון Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — מציין מקום לאיטרציה עתידית. ראו [`mcp/`](mcp/).

</div>

## התקנה כ-Claude Code skills

<div dir="rtl">

הדרך הכי קלה — הדבק את זה ב-Claude Code והוא יגדיר הכול:

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

או שתגדיר ידנית:

</div>

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

<div dir="rtl">

הוסיפו אל `~/.claude/settings.json`:

</div>

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

<div dir="rtl">

הפעילו מחדש את Claude Code.

</div>

## הרצת הדוגמאות

```bash
cd examples
go test ./...
```

## רישיון

<div dir="rtl">

MIT — ראו [LICENSE](LICENSE).

</div>

## מחבר

<div dir="rtl">

Yevgeniy Achin · <https://github.com/amazopic>

</div>
