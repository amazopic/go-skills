<p align="center"><a href="https://amazopic.github.io/go-skills/"><b>go-skills</b></a></p>
<p align="center">Go 서비스를 위한 아키텍처 패턴과 방법론 — Claude Code skills로 패키지화.</p>

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

## 내용

- **방법론** — 프로덕션 Go 서비스 구축을 위한 18장 실용 가이드 (디렉터리 구조, 레이어드 아키텍처, 수동 DI, 설정, 재시도, 스토리지, transport, 작업, 로깅, 유효성 검사, 오류, 테스트, 빌드, 배포). [전체 표준 문서](skills/methodology/00-canonical-full.md)를 읽거나 개별 [chapter-skill](METHODOLOGY.md)을 선택하세요.
- **패턴** — 10개 카테고리에 걸친 52개 이상의 항목: 생성형, 구조형, 행동형, 동시성, 동기화, 메시징, 안정성, 프로파일링, 관용구, 안티패턴. [PATTERNS.md](PATTERNS.md) 참조.
- **예제** — `_test.go`가 포함된 52개의 실행 가능한 Go 패턴 패키지. [`examples/`](examples/) 아래에 공유 모듈 하나.
- **사이트** — Linear 스타일의 GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP server** — 미래 이터레이션을 위한 플레이스홀더. [`mcp/`](mcp/) 참조.

## Claude Code skills로 설치

```bash
git clone https://github.com/amazopic/go-skills.git ~/.claude/plugins/go-skills
```

`~/.claude/settings.json`에 추가:

```json
{
  "skillSources": [
    "~/.claude/plugins/go-skills/skills"
  ]
}
```

Claude Code를 재시작하세요.

## 예제 실행

```bash
cd examples
go test ./...
```

## 라이선스

MIT — [LICENSE](LICENSE) 참조.

## 저자

Yevgeniy Achin · <https://github.com/amazopic>
