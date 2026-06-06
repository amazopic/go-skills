<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Các mẫu thiết kế kiến trúc và phương pháp luận cho dịch vụ Go — đóng gói thành Claude Code skills.</p>

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

## Bên trong có gì

- **Phương pháp luận** — một cẩm nang 18 chương để xây dựng các dịch vụ Go cho production (bố cục thư mục, kiến trúc phân tầng, DI thủ công, cấu hình, thử lại, lưu trữ, vận chuyển, tác vụ nền, ghi log, kiểm tra hợp lệ, xử lý lỗi, kiểm thử, build, triển khai). Đọc [tài liệu đầy đủ chính thống](skills/methodology/00-canonical-full.md) hoặc chọn một [chương-skill](METHODOLOGY.md).
- **Mẫu thiết kế** — hơn 52 mục trải khắp 10 nhóm: khởi tạo, cấu trúc, hành vi, đồng thời, đồng bộ hóa, truyền tin, ổn định, profiling, thành ngữ, phản mẫu. Xem [PATTERNS.md](PATTERNS.md).
- **Ví dụ** — 52 package mẫu thiết kế Go chạy được kèm `_test.go`. Chung một module dưới [`examples/`](examples/).
- **Trang web** — GitHub Pages theo phong cách Linear: <https://amazopic.github.io/go-skills/>
- **Máy chủ MCP** — chỗ giữ chỗ cho một bản lặp tương lai. Xem [`mcp/`](mcp/).

## Cài đặt làm Claude Code skills

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

Thêm vào `~/.claude/settings.json`:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Khởi động lại Claude Code.

## Chạy các ví dụ

```bash
cd examples
go test ./...
```

## Giấy phép

MIT — xem [LICENSE](LICENSE).

## Tác giả

Yevgeniy Achin · <https://github.com/amazopic>
