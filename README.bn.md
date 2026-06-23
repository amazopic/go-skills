<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go সার্ভিসের জন্য আর্কিটেকচারাল প্যাটার্ন ও মেথডোলজি — Claude Code skills হিসেবে প্যাকেজকৃত।</p>

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

> **go-skills-এ নতুন? [Project Assessment](skills/workflow/project-assessment.md) দিয়ে শুরু করুন।** এটি আপনার বিদ্যমান Go প্রজেক্টের দিকে তাক করুন এবং একটি ম্যাচিউরিটি স্কোর সহ একটি অগ্রাধিকারভিত্তিক রোডম্যাপ পান যেখানে প্রতিটি আইটেম প্রয়োগের জন্য সঠিক go-skills প্যাটার্ন বা অধ্যায়ের সাথে লিঙ্ক করা — রিড-অনলি, এটি কখনো আপনার কোড পরিবর্তন করে না।

## ভেতরে কী আছে

- **মেথডোলজি** — প্রোডাকশন Go সার্ভিস তৈরির জন্য একটি ১৮-অধ্যায়ের প্লেবুক (ডিরেক্টরি লেআউট, লেয়ার্ড আর্কিটেকচার, ম্যানুয়াল DI, কনফিগারেশন, রিট্রাই, স্টোরেজ, ট্রান্সপোর্ট, জব, লগিং, ভ্যালিডেশন, এরর, টেস্টিং, বিল্ড, ডিপ্লয়)। [ক্যানোনিক্যাল সম্পূর্ণ ডকুমেন্ট](skills/methodology/00-canonical-full.md) পড়ুন অথবা একটি [chapter-skill](METHODOLOGY.md) বেছে নিন।
- **প্যাটার্ন** — ১০টি ক্যাটাগরি জুড়ে ৫২টি এন্ট্রি: ক্রিয়েশনাল, স্ট্রাকচারাল, বিহেভিওরাল, কনকারেন্সি, সিনক্রোনাইজেশন, মেসেজিং, স্ট্যাবিলিটি, প্রোফাইলিং, ইডিয়ম, অ্যান্টি-প্যাটার্ন। দেখুন [PATTERNS.md](PATTERNS.md)।
- **উদাহরণ** — `_test.go` সহ ৫২টি রানযোগ্য Go প্যাটার্ন প্যাকেজ। [`examples/`](examples/)-এর অধীনে একটি শেয়ার্ড মডিউল।
- **ওয়ার্কফ্লো** — রানযোগ্য টিম-ভিত্তিক স্কিল: [Project Assessment](skills/workflow/project-assessment.md) (একটি বিদ্যমান প্রজেক্ট মূল্যায়ন), [Feature Development](skills/workflow/feature-development.md) (একটি ফিচার তৈরি), [Security Code Review](skills/workflow/security-review.md) (অডিট + প্রতিকার)।
- **সাইট** — Linear-স্টাইলড GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP server** — একটি ভবিষ্যৎ ইটারেশনের জন্য প্লেসহোল্ডার। দেখুন [`mcp/`](mcp/)।

## Claude Code skills হিসেবে ইনস্টল করুন

সবচেয়ে সহজ উপায় — এটি Claude Code-এ পেস্ট করুন, বাকি সব হয়ে যাবে:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

অথবা নিজে হাতে সেটআপ করুন:

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

`~/.claude/settings.json`-এ যোগ করুন:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Claude Code রিস্টার্ট করুন।

## উদাহরণ চালান

```bash
cd examples
go test ./...
```

## লাইসেন্স

MIT — দেখুন [LICENSE](LICENSE)।

## লেখক

Yevgeniy Achin · <https://github.com/amazopic>
