// i18n.js — UI string switching across 23 locales, lazy-loaded per locale.
// Default: en (bundled inline below as the default locale + universal fallback).
// Every other locale lives in ./locales/<code>.js and is dynamic-imported on demand,
// so a visitor downloads only the dictionary they actually use.
// Activate with ?lang=xx (also written to localStorage for stickiness).

export const supportedLocales = [
  { code: 'en', label: 'English', native: 'English' },
  { code: 'ru', label: 'Russian', native: 'Русский' },
  { code: 'fr', label: 'French', native: 'Français' },
  { code: 'de', label: 'German', native: 'Deutsch' },
  { code: 'uk', label: 'Ukrainian', native: 'Українська' },
  { code: 'sl', label: 'Slovenian', native: 'Slovenščina' },
  { code: 'it', label: 'Italian', native: 'Italiano' },
  { code: 'es', label: 'Spanish', native: 'Español' },
  { code: 'zh', label: 'Chinese', native: '中文' },
  { code: 'ja', label: 'Japanese', native: '日本語' },
  { code: 'ko', label: 'Korean', native: '한국어' },
  { code: 'ar', label: 'Arabic', native: 'العربية', rtl: true },
  { code: 'pt-BR', label: 'Portuguese (BR)', native: 'Português (BR)' },
  { code: 'tr', label: 'Turkish', native: 'Türkçe' },
  { code: 'id', label: 'Indonesian', native: 'Bahasa Indonesia' },
  { code: 'vi', label: 'Vietnamese', native: 'Tiếng Việt' },
  { code: 'hi', label: 'Hindi', native: 'हिन्दी' },
  { code: 'zh-TW', label: 'Chinese (Traditional)', native: '繁體中文' },
  { code: 'pl', label: 'Polish', native: 'Polski' },
  { code: 'th', label: 'Thai', native: 'ไทย' },
  { code: 'he', label: 'Hebrew', native: 'עברית', rtl: true },
  { code: 'bn', label: 'Bengali', native: 'বাংলা' },
  { code: 'ur', label: 'Urdu', native: 'اردو', rtl: true },
];

export const defaultLocale = 'en';

// en — bundled inline: the default locale and the fallback for any key missing elsewhere.
const EN = {
  "meta.title": "go-skills — Senior Go developer. Rocket engine included.",
  "meta.description": "Add one skill — code Go like a senior with 10+ years across hundreds of production projects. Claude Code + you + go-skills = idiomatic Go, no rookie mistakes.",
  "hero.brand": "go-skills",
  "hero.title": "Senior Go developer. Rocket engine included.",
  "hero.lede": "Add one skill — code Go like a senior with 10+ years across hundreds of production projects. Claude Code + you + go-skills = no rookie mistakes, idiomatic Go from prompt one.",
  "hero.cta.catalog": "Browse catalog",
  "hero.cta.methodology": "Read methodology",
  "hero.cta.github": "View on GitHub →",
  "hero.cta.assess": "Assess your project",
  "quickstart.eyebrow": "Quick start",
  "quickstart.title": "Done with one Claude Code prompt.",
  "quickstart.lede": "Paste this into Claude Code. It clones go-skills, registers it as a skill source, reloads, and confirms — about 30 seconds, no manual edits.",
  "quickstart.copy": "Copy prompt",
  "quickstart.prompt": "Install go-skills as a Claude Code skill source.\n\n1. Clone https://github.com/amazopic/go-skills to ~/.claude/plugins/go-skills\n2. Add \"~/.claude/plugins/go-skills/skills\" to skillSources in ~/.claude/settings.json\n3. Reload settings\n4. List the skills you now have access to\n\nAfter install I want to code Go with senior-level patterns and the 18-chapter service-architecture methodology.",
  "assess.eyebrow": "Start here",
  "assess.title": "Already have a Go project? Point go-skills at it.",
  "assess.lede": "The fastest way in: a five-role team reads your codebase and returns a maturity score plus a prioritized roadmap — every item linked to the exact go-skills pattern or chapter to apply. Read-only; it never touches your code.",
  "assess.copy": "Copy prompt",
  "assess.prompt": "Using the go-skills \"Project Assessment\" workflow skill, assess my existing Go project.\n\nRead the codebase under <path> and produce a maturity assessment plus a prioritized roadmap where every finding cites the specific go-skills methodology chapter or pattern to apply. Read-only — do not change any code. Then propose the first sprint.",
  "vibe.eyebrow": "Three pillars",
  "vibe.title": "Methodology · Patterns · Skills.",
  "vibe.lede": "A reference for senior Go engineers — opinionated, source-derived, no decoration.",
  "vibe.p1.title": "Methodology",
  "vibe.p1.body": "An 18-chapter playbook for building production Go services.",
  "vibe.p2.title": "Patterns",
  "vibe.p2.body": "52 senior-grade patterns across 10 categories — every entry written by a senior Go consultant, race-safe runnable examples where applicable.",
  "vibe.p3.title": "Skills",
  "vibe.p3.body": "Each pattern and chapter is a Claude Code skill.",
  "numbers.eyebrow": "Inventory",
  "numbers.title": "By the numbers.",
  "numbers.s1": "Categories",
  "numbers.s2": "Patterns",
  "numbers.s3": "Methodology chapters",
  "numbers.s4": "Runnable examples",
  "numbers.s5": "Languages",
  "meth.eyebrow": "Service Methodology",
  "meth.title": "A reference for building services in Go.",
  "meth.lede": "An anonymized, project-agnostic playbook. Read end-to-end as the canonical 930-line document, or pick a chapter-skill below.",
  "meth.cta": "Read the full canonical document →",
  "cat.eyebrow": "Catalog",
  "cat.title": "All patterns.",
  "cat.lede": "All 52 entries are filled with senior-grade content. Most link to a runnable, race-safe Go example.",
  "cat.workflow": "Workflow",
  "cat.concurrency": "Concurrency",
  "cat.stability": "Stability",
  "cat.messaging": "Messaging",
  "cat.synchronization": "Synchronization",
  "cat.idiom": "Idioms",
  "cat.creational": "Creational",
  "cat.structural": "Structural",
  "cat.behavioral": "Behavioral",
  "cat.profiling": "Profiling",
  "cat.antipatterns": "Anti-Patterns",
  "spec.eyebrow": "Specimens",
  "spec.title": "Three patterns, in full.",
  "inst.eyebrow": "Install",
  "inst.title": "Use as Claude Code skills.",
  "inst.lede": "Drop the repo into your skills path. Claude Code will surface skills from this repo when their description matches your task.",
  "inst.settings": "Then add to ~/.claude/settings.json:",
  "inst.restart": "Restart Claude Code.",
  "inst.mcp": "An MCP server that exposes the same content programmatically is on the roadmap.",
  "road.eyebrow": "Roadmap",
  "road.title": "What's next.",
  "road.l1": "Translate all 52 skill bodies into the 11 secondary languages (currently English-only).",
  "road.l2": "MCP server: tools.list_patterns, tools.get_pattern, methodology resources.",
  "road.l3": "Translate methodology and pattern bodies into the 12 secondary languages.",
  "road.l4": "CI: link checks, go tests, Lighthouse, JSON-LD validation.",
  "road.l5": "Slash command /go-skill <pattern> to inject a pattern on demand.",
  "faq.eyebrow": "FAQ",
  "faq.title": "Frequently asked.",
  "faq.q1": "Why build this when go-patterns already exists on GitHub?",
  "faq.a1": "go-skills consolidates two upstream archives, adds a service methodology, and packages everything as Claude Code skills with rich frontmatter.",
  "faq.q2": "Are the patterns idiomatic Go or GoF translations?",
  "faq.a2": "Both. Concurrency, messaging, stability, profiling, idiom are Go-native. GoF patterns are translated with idiomatic adjustments.",
  "faq.q3": "How were the patterns vetted?",
  "faq.a3": "Every one of the 52 entries was written by a senior Go consultant — idiomatic Go, opinionated guidance. All 52 example packages pass go test -race.",
  "faq.q4": "Will the examples run on my machine?",
  "faq.a4": "Yes — one shared go.mod at examples/. Run: cd examples && go test ./...",
  "faq.q5": "When will the MCP server land?",
  "faq.a5": "A future iteration. mcp/ holds a roadmap-only README in v1.",
  "faq.q6": "Why Linear's design system?",
  "faq.a6": "Linear's near-black + single-accent aesthetic matches the content — dense, technical, no decoration.",
  "faq.q7": "License?",
  "faq.a7": "MIT.",
  "faq.q8": "Does this support the methodology in non-English?",
  "faq.a8": "README × 23 languages. Pattern and methodology bodies are English-only in v1."
};

// Runtime cache of loaded locale dictionaries, seeded with the inline en fallback.
const cache = { en: EN };

function isSupported(code) {
  return supportedLocales.some(l => l.code === code);
}

// Dynamic-import a locale chunk on demand, caching it. en (and any failure) resolves to EN.
async function loadLocale(code) {
  if (cache[code]) return cache[code];
  if (!isSupported(code) || code === defaultLocale) return EN;
  try {
    const mod = await import(`./locales/${code}.js`);
    cache[code] = mod.default;
    return mod.default;
  } catch (e) {
    console.error(`i18n: failed to load locale "${code}", falling back to en`, e);
    return EN;
  }
}

// Default to English. Honor explicit ?lang=xx and prior user choice (localStorage),
// but never auto-pick from navigator.language.
export function getLocale() {
  const params = new URLSearchParams(window.location.search);
  const fromQuery = params.get('lang');
  if (fromQuery && isSupported(fromQuery)) return fromQuery;
  const fromStorage = localStorage.getItem('go-skills.lang');
  if (fromStorage && isSupported(fromStorage)) return fromStorage;
  return defaultLocale;
}

export async function setLocale(code) {
  if (!isSupported(code)) return;
  localStorage.setItem('go-skills.lang', code);
  await apply(code);
}

// Synchronous lookup against already-loaded dictionaries (en fallback, then the key itself).
export function t(code, key) {
  return (cache[code] && cache[code][key]) || EN[key] || key;
}

export async function apply(code) {
  const locale = supportedLocales.find(l => l.code === code) || supportedLocales[0];
  const dict = await loadLocale(code);
  const tr = (key) => dict[key] || EN[key] || key;
  document.documentElement.setAttribute('lang', code);
  document.documentElement.setAttribute('dir', locale.rtl ? 'rtl' : 'ltr');
  document.querySelectorAll('[data-i18n]').forEach(el => {
    el.textContent = tr(el.getAttribute('data-i18n'));
  });
  document.querySelectorAll('[data-i18n-html]').forEach(el => {
    el.innerHTML = tr(el.getAttribute('data-i18n-html'));
  });
  document.title = tr('meta.title');
  // Point GitHub README CTAs to the locale-specific README file.
  const repo = 'https://github.com/amazopic/go-skills';
  const readmeUrl = code === defaultLocale ? repo : `${repo}/blob/main/README.${code}.md`;
  document.querySelectorAll('[data-github-readme]').forEach(el => {
    el.setAttribute('href', readmeUrl);
  });
}
