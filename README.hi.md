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
  <a href="README.pl.md">Polski</a>
</p>

---

## अंदर क्या है

- **पद्धति (Methodology)** — प्रोडक्शन Go सेवाएँ बनाने के लिए 18-अध्याय की एक प्लेबुक (डायरेक्टरी लेआउट, स्तरित आर्किटेक्चर, मैनुअल DI, कॉन्फ़िगरेशन, रिट्राई, स्टोरेज, ट्रांसपोर्ट, जॉब्स, लॉगिंग, सत्यापन, एरर, टेस्टिंग, बिल्ड, डिप्लॉय)। [प्रामाणिक पूर्ण दस्तावेज़](skills/methodology/00-canonical-full.md) पढ़ें या एक [अध्याय-skill](METHODOLOGY.md) चुनें।
- **पैटर्न (Patterns)** — 9 श्रेणियों में 50+ प्रविष्टियाँ: creational, structural, behavioral, concurrency, synchronization, messaging, stability, profiling, idiom, anti-patterns। देखें [PATTERNS.md](PATTERNS.md)।
- **उदाहरण (Examples)** — `_test.go` के साथ 23 चलाने योग्य Go पैटर्न पैकेज। [`examples/`](examples/) के अंतर्गत एक साझा मॉड्यूल।
- **साइट (Site)** — Linear-शैली वाला GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP सर्वर** — भविष्य के एक संस्करण के लिए प्लेसहोल्डर। देखें [`mcp/`](mcp/)।

## Claude Code skills के रूप में इंस्टॉल करें

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
