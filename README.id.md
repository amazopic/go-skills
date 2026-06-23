<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Pola arsitektur dan metodologi untuk layanan Go — dikemas sebagai Claude Code skills.</p>

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

> **Baru mengenal go-skills? Mulai dengan [Project Assessment](skills/workflow/project-assessment.md).** Arahkan ke proyek Go Anda yang sudah ada dan dapatkan skor kematangan beserta peta jalan yang diprioritaskan, di mana setiap item ditautkan ke pola atau bab go-skills tepat yang harus diterapkan — hanya-baca, tidak pernah mengubah kode Anda.

## Apa isinya

- **Metodologi** — panduan 18 bab untuk membangun layanan Go di produksi (tata letak direktori, arsitektur berlapis, DI manual, konfigurasi, retry, penyimpanan, transport, job, logging, validasi, error, pengujian, build, deploy). Baca [dokumen lengkap kanonik](skills/methodology/00-canonical-full.md) atau pilih sebuah [bab-skill](METHODOLOGY.md).
- **Pola** — 52 entri di 10 kategori: kreasional, struktural, perilaku, konkurensi, sinkronisasi, pesan, stabilitas, profiling, idiom, anti-pola. Lihat [PATTERNS.md](PATTERNS.md).
- **Contoh** — 52 paket pola Go yang dapat dijalankan dengan `_test.go`. Satu modul bersama di bawah [`examples/`](examples/).
- **Alur Kerja** — skill berbasis tim yang dapat dijalankan: [Project Assessment](skills/workflow/project-assessment.md) (menilai proyek yang sudah ada), [Feature Development](skills/workflow/feature-development.md) (membangun sebuah fitur), [Security Code Review](skills/workflow/security-review.md) (audit + perbaikan).
- **Situs** — GitHub Pages bergaya Linear: <https://amazopic.github.io/go-skills/>
- **Server MCP** — placeholder untuk iterasi mendatang. Lihat [`mcp/`](mcp/).

## Pasang sebagai Claude Code skills

Paling gampang — tempel ini ke Claude Code dan semua langsung beres:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

Atau atur manual:

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Tambahkan ke `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Mulai ulang Claude Code.

## Menjalankan contoh

```bash
cd examples
go test ./...
```

## Lisensi

MIT — lihat [LICENSE](LICENSE).

## Penulis

Yevgeniy Achin · <https://github.com/amazopic>
