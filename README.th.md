<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">แพตเทิร์นสถาปัตยกรรมและระเบียบวิธีสำหรับเซอร์วิส Go — แพ็กเป็น Claude Code skills</p>

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

> **เพิ่งรู้จัก go-skills? เริ่มที่ [Project Assessment](skills/workflow/project-assessment.md)** ชี้มันไปที่โปรเจกต์ Go ที่มีอยู่ของคุณ แล้วรับคะแนนความสมบูรณ์พร้อมโรดแมปที่จัดลำดับความสำคัญ โดยทุกรายการเชื่อมโยงกับแพตเทิร์นหรือบท go-skills ที่ตรงเป๊ะให้นำไปใช้ — อ่านอย่างเดียว ไม่แตะโค้ดของคุณ

## มีอะไรอยู่ข้างใน

- **ระเบียบวิธี** — คู่มือปฏิบัติ 18 บทสำหรับสร้างเซอร์วิส Go ระดับโปรดักชัน (การจัดวางไดเรกทอรี, สถาปัตยกรรมแบบเลเยอร์, DI แบบทำมือ, การตั้งค่า, การลองใหม่, ที่จัดเก็บข้อมูล, การขนส่ง, งาน, การบันทึกล็อก, การตรวจสอบความถูกต้อง, ข้อผิดพลาด, การทดสอบ, การคอมไพล์, การดีพลอย) อ่าน[เอกสารฉบับสมบูรณ์เต็ม](skills/methodology/00-canonical-full.md) หรือเลือก[chapter-skill](METHODOLOGY.md)
- **แพตเทิร์น** — 52 รายการใน 10 หมวด: creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns ดูที่ [PATTERNS.md](PATTERNS.md)
- **ตัวอย่าง** — 52 แพ็กเกจแพตเทิร์น Go ที่รันได้พร้อม `_test.go` มีหนึ่งโมดูลที่ใช้ร่วมกันภายใต้ [`examples/`](examples/)
- **Workflows** — สกิลที่รันได้และอิงทีม: [Project Assessment](skills/workflow/project-assessment.md) (ประเมินโปรเจกต์ที่มีอยู่), [Feature Development](skills/workflow/feature-development.md) (สร้างฟีเจอร์), [Security Code Review](skills/workflow/security-review.md) (ตรวจสอบ + แก้ไข)
- **เว็บไซต์** — GitHub Pages สไตล์ Linear: <https://amazopic.github.io/go-skills/>
- **MCP server** — ส่วนที่เว้นว่างไว้สำหรับรุ่นในอนาคต ดูที่ [`mcp/`](mcp/)

## ติดตั้งเป็น Claude Code skills

ง่ายสุด — วางอันนี้ลงใน Claude Code แล้วมันจะตั้งค่าให้ทั้งหมด:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

หรือจะตั้งค่าเองก็ได้:

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

เพิ่มลงใน `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

รีสตาร์ท Claude Code

## รันตัวอย่าง

```bash
cd examples
go test ./...
```

## ใบอนุญาต

MIT — ดูที่ [LICENSE](LICENSE)

## ผู้เขียน

Yevgeniy Achin · <https://github.com/amazopic>
