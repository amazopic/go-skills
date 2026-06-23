<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go सेवाओं के लिए आर्किटेक्चरल पैटर्न और पद्धति — Claude Code skills के रूप में पैकेज की गई।</p>

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

> **go-skills में नए हैं? [Project Assessment](skills/workflow/project-assessment.md) से शुरू करें।** इसे अपने मौजूदा Go प्रोजेक्ट पर लगाएँ और एक मैच्योरिटी स्कोर के साथ एक प्राथमिकता-क्रम वाला रोडमैप पाएँ, जहाँ हर आइटम ठीक उसी go-skills पैटर्न या अध्याय से जुड़ा होता है जिसे लागू करना है — केवल-पठन, यह आपके कोड को कभी नहीं बदलता।

## अंदर क्या है

- **पद्धति (Methodology)** — प्रोडक्शन Go सेवाएँ बनाने के लिए 18-अध्याय की एक प्लेबुक (डायरेक्टरी लेआउट, स्तरित आर्किटेक्चर, मैनुअल DI, कॉन्फ़िगरेशन, रिट्राई, स्टोरेज, ट्रांसपोर्ट, जॉब्स, लॉगिंग, सत्यापन, एरर, टेस्टिंग, बिल्ड, डिप्लॉय)। [प्रामाणिक पूर्ण दस्तावेज़](skills/methodology/00-canonical-full.md) पढ़ें या एक [अध्याय-skill](METHODOLOGY.md) चुनें।
- **पैटर्न (Patterns)** — 10 श्रेणियों में 52 प्रविष्टियाँ: creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns। देखें [PATTERNS.md](PATTERNS.md)।
- **उदाहरण (Examples)** — `_test.go` के साथ 52 चलाने योग्य Go पैटर्न पैकेज। [`examples/`](examples/) के अंतर्गत एक साझा मॉड्यूल।
- **वर्कफ़्लो (Workflows)** — चलाने योग्य टीम-आधारित skills: [Project Assessment](skills/workflow/project-assessment.md) (एक मौजूदा प्रोजेक्ट का आकलन करें), [Feature Development](skills/workflow/feature-development.md) (एक फ़ीचर बनाएँ), [Security Code Review](skills/workflow/security-review.md) (ऑडिट + सुधार)।
- **साइट (Site)** — Linear-शैली वाला GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP सर्वर** — भविष्य के एक संस्करण के लिए प्लेसहोल्डर। देखें [`mcp/`](mcp/)।

## Claude Code skills के रूप में इंस्टॉल करें

सबसे आसान तरीका — इसे Claude Code में पेस्ट करें, बाकी सब हो जाएगा:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

या फिर खुद से सेटअप करें:

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

`~/.claude/settings.json` में जोड़ें:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Claude Code को फिर से शुरू करें।

## उदाहरण चलाएँ

```bash
cd examples
go test ./...
```

## लाइसेंस

MIT — देखें [LICENSE](LICENSE)।

## लेखक

Yevgeniy Achin · <https://github.com/amazopic>
