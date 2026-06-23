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
  <a href="README.pl.md">Polski</a> ·
  <a href="README.th.md">ไทย</a> ·
  <a href="README.he.md">עברית</a> ·
  <a href="README.bn.md">বাংলা</a> ·
  <a href="README.ur.md">اردو</a>
</p>

---

> **go-skills가 처음이신가요? [Project Assessment](skills/workflow/project-assessment.md)부터 시작하세요.** 기존 Go 프로젝트에 적용하면 성숙도 점수와 우선순위가 매겨진 로드맵을 얻을 수 있으며, 각 항목은 적용할 정확한 go-skills 패턴이나 챕터에 연결됩니다 — 읽기 전용이라 코드를 절대 변경하지 않습니다.

## 내용

- **방법론** — 프로덕션 Go 서비스 구축을 위한 18장 실용 가이드 (디렉터리 구조, 레이어드 아키텍처, 수동 DI, 설정, 재시도, 스토리지, transport, 작업, 로깅, 유효성 검사, 오류, 테스트, 빌드, 배포). [전체 표준 문서](skills/methodology/00-canonical-full.md)를 읽거나 개별 [chapter-skill](METHODOLOGY.md)을 선택하세요.
- **패턴** — 10개 카테고리에 걸친 52개 이상의 항목: 생성형, 구조형, 행동형, 동시성, 동기화, 메시징, 안정성, 프로파일링, 관용구, 안티패턴. [PATTERNS.md](PATTERNS.md) 참조.
- **예제** — `_test.go`가 포함된 52개의 실행 가능한 Go 패턴 패키지. [`examples/`](examples/) 아래에 공유 모듈 하나.
- **워크플로** — 팀 기반의 실행 가능한 skill: [Project Assessment](skills/workflow/project-assessment.md) (기존 프로젝트 진단), [Feature Development](skills/workflow/feature-development.md) (기능 구축), [Security Code Review](skills/workflow/security-review.md) (감사 + 수정).
- **사이트** — Linear 스타일의 GitHub Pages: <https://amazopic.github.io/go-skills/>
- **MCP server** — 미래 이터레이션을 위한 플레이스홀더. [`mcp/`](mcp/) 참조.

## Claude Code skills로 설치

가장 쉬운 방법 — 이걸 Claude Code에 붙여넣으면 전부 설정됩니다:

```text
Set up go-skills for me — handle everything, I'll just vibe.

1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills
2. Add "~/.claude/plugins/go-skills/skills" to skillSources in ~/.claude/settings.json
3. Reload settings and confirm the skills are live
4. Tell me in one line what you can now help with — patterns, methodology, workflows

Then ask what I'm building and start coding Go at senior level — idiomatic, race-safe, no rookie mistakes.
```

또는 직접 설정하기:

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
