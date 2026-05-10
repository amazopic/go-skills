// i18n.js — UI string switching across 13 locales
// Default: en. Other locales fall back to en for missing keys.
// Activate with ?lang=xx (also writes to localStorage for stickiness).

export const supportedLocales = [
  { code: 'en', label: 'English',     native: 'English'    },
  { code: 'ru', label: 'Russian',     native: 'Русский'    },
  { code: 'fr', label: 'French',      native: 'Français'   },
  { code: 'de', label: 'German',      native: 'Deutsch'    },
  { code: 'uk', label: 'Ukrainian',   native: 'Українська' },
  { code: 'sl', label: 'Slovenian',   native: 'Slovenščina'},
  { code: 'it', label: 'Italian',     native: 'Italiano'   },
  { code: 'es', label: 'Spanish',     native: 'Español'    },
  { code: 'zh', label: 'Chinese',     native: '中文'        },
  { code: 'ja', label: 'Japanese',    native: '日本語'      },
  { code: 'ko', label: 'Korean',      native: '한국어'      },
  { code: 'ar', label: 'Arabic',      native: 'العربية',     rtl: true },
];

export const defaultLocale = 'en';

export const messages = {
  en: {
    'meta.title':   'go-skills — Patterns & methodology for Go services',
    'meta.description': 'A Go knowledge base — 50+ patterns + 18-chapter service architecture methodology — packaged as Claude Code skills.',
    'hero.brand':   'go-skills',
    'hero.title':   'Architectural patterns and methodology for Go services.',
    'hero.lede':    'A curated knowledge base — 50+ design and concurrency patterns plus an 18-chapter service-architecture playbook — packaged as Claude Code skills, exposed via this site, and (soon) an MCP server.',
    'hero.cta.catalog':     'Browse catalog',
    'hero.cta.methodology': 'Read methodology',
    'hero.cta.github':      'View on GitHub →',

    'vibe.eyebrow': 'Three pillars',
    'vibe.title':   'Methodology · Patterns · Skills.',
    'vibe.lede':    'A reference for senior Go engineers — opinionated, source-derived, no decoration.',
    'vibe.p1.title':'Methodology',
    'vibe.p1.body': 'An 18-chapter playbook for building production Go services.',
    'vibe.p2.title':'Patterns',
    'vibe.p2.body': '50+ classic and Go-flavored patterns across 9 categories.',
    'vibe.p3.title':'Skills',
    'vibe.p3.body': 'Each pattern and chapter is a Claude Code skill.',

    'numbers.eyebrow':'Inventory',
    'numbers.title':  'By the numbers.',
    'numbers.s1':'Categories',
    'numbers.s2':'Patterns',
    'numbers.s3':'Methodology chapters',
    'numbers.s4':'Runnable examples',
    'numbers.s5':'Languages',

    'meth.eyebrow': 'Service Methodology',
    'meth.title':   'A reference for building services in Go.',
    'meth.lede':    'An anonymized, project-agnostic playbook. Read end-to-end as the canonical 930-line document, or pick a chapter-skill below.',
    'meth.cta':     'Read the full canonical document →',

    'cat.eyebrow':  'Catalog',
    'cat.title':    'All patterns.',
    'cat.lede':     'Filled cards have a body and (where applicable) a runnable example. Stubs are listed for completeness.',
    'cat.concurrency': 'Concurrency',
    'cat.stability':   'Stability',
    'cat.messaging':   'Messaging',
    'cat.synchronization':'Synchronization',
    'cat.idiom':       'Idioms',
    'cat.creational':  'Creational',
    'cat.structural':  'Structural',
    'cat.behavioral':  'Behavioral',
    'cat.profiling':   'Profiling',
    'cat.antipatterns':'Anti-Patterns',

    'spec.eyebrow':'Specimens',
    'spec.title':  'Three patterns, in full.',

    'inst.eyebrow':'Install',
    'inst.title':  'Use as Claude Code skills.',
    'inst.lede':   'Drop the repo into your skills path. Claude Code will surface skills from this repo when their description matches your task.',
    'inst.settings':'Then add to ~/.claude/settings.json:',
    'inst.restart': 'Restart Claude Code.',
    'inst.mcp':     'An MCP server that exposes the same content programmatically is on the roadmap.',

    'road.eyebrow':'Roadmap',
    'road.title':  "What's next.",
    'road.l1':'Fill the remaining stubbed pattern bodies — port from upstream Go pattern archives.',
    'road.l2':'MCP server: tools.list_patterns, tools.get_pattern, methodology resources.',
    'road.l3':'Translate methodology and pattern bodies into the 12 secondary languages.',
    'road.l4':'CI: link checks, go tests, Lighthouse, JSON-LD validation.',
    'road.l5':'Slash command /go-skill <pattern> to inject a pattern on demand.',

    'faq.eyebrow':'FAQ',
    'faq.title':  'Frequently asked.',
    'faq.q1':'Why build this when go-patterns already exists on GitHub?',
    'faq.a1':'go-skills consolidates two upstream archives, adds a service methodology, and packages everything as Claude Code skills with rich frontmatter.',
    'faq.q2':'Are the patterns idiomatic Go or GoF translations?',
    'faq.a2':'Both. Concurrency, messaging, stability, profiling, idiom are Go-native. GoF patterns are translated with idiomatic adjustments.',
    'faq.q3':'Why are some entries marked as stubs?',
    'faq.a3':'The catalog mirrors a public taxonomy; not every node has a body yet.',
    'faq.q4':'Will the examples run on my machine?',
    'faq.a4':'Yes — one shared go.mod at examples/. Run: cd examples && go test ./...',
    'faq.q5':'When will the MCP server land?',
    'faq.a5':'A future iteration. mcp/ holds a roadmap-only README in v1.',
    'faq.q6':"Why Linear's design system?",
    'faq.a6':"Linear's near-black + single-accent aesthetic matches the content — dense, technical, no decoration.",
    'faq.q7':'License?',
    'faq.a7':'MIT.',
    'faq.q8':'Does this support the methodology in non-English?',
    'faq.a8':'README × 13 languages. Pattern and methodology bodies are English-only in v1.',
  },
  ru: {
    'meta.title': 'go-skills — Архитектурные паттерны и методология Go-сервисов',
    'meta.description': 'База знаний по Go — 50+ паттернов и 18-главная методология построения сервисов, упакованная как Claude Code skills.',
    'hero.brand': 'go-skills',
    'hero.title': 'Архитектурные паттерны и методология Go-сервисов.',
    'hero.lede':  'Курированная база знаний — 50+ паттернов проектирования и параллелизма + 18 глав методологии построения сервисов. Упакована как Claude Code skills, опубликована на этом сайте, в будущем — MCP-сервер.',
    'hero.cta.catalog': 'Каталог',
    'hero.cta.methodology': 'Методология',
    'hero.cta.github': 'GitHub →',
    'vibe.eyebrow': 'Три столпа',
    'vibe.title': 'Методология · Паттерны · Skills.',
    'numbers.title': 'В цифрах.',
    'numbers.s1':'Категорий', 'numbers.s2':'Паттернов', 'numbers.s3':'Глав методологии', 'numbers.s4':'Примеров', 'numbers.s5':'Языков',
    'cat.title':'Все паттерны.',
    'inst.title':'Установка как Claude Code skill.',
    'road.title':'Дальше.',
    'faq.title':'Часто спрашивают.',
  },
  fr: { 'hero.cta.catalog':'Catalogue', 'hero.cta.methodology':'Méthodologie', 'cat.title':'Tous les patterns.' },
  de: { 'hero.cta.catalog':'Katalog',   'hero.cta.methodology':'Methodologie', 'cat.title':'Alle Patterns.' },
  uk: { 'hero.cta.catalog':'Каталог',   'hero.cta.methodology':'Методологія',  'cat.title':'Всі патерни.' },
  sl: { 'hero.cta.catalog':'Katalog',   'hero.cta.methodology':'Metodologija', 'cat.title':'Vsi vzorci.' },
  it: { 'hero.cta.catalog':'Catalogo',  'hero.cta.methodology':'Metodologia',  'cat.title':'Tutti i pattern.' },
  es: { 'hero.cta.catalog':'Catálogo',  'hero.cta.methodology':'Metodología',  'cat.title':'Todos los patrones.' },
  zh: { 'hero.cta.catalog':'目录',       'hero.cta.methodology':'方法论',         'cat.title':'所有模式。' },
  ja: { 'hero.cta.catalog':'カタログ',    'hero.cta.methodology':'方法論',        'cat.title':'全パターン。' },
  ko: { 'hero.cta.catalog':'카탈로그',    'hero.cta.methodology':'방법론',        'cat.title':'모든 패턴.' },
  ar: { 'hero.cta.catalog':'الفهرس',     'hero.cta.methodology':'المنهجية',      'cat.title':'كل الأنماط.' },
};

// Default to English. Honor explicit ?lang=xx and prior user choice
// (localStorage), but never auto-pick from navigator.language.
export function getLocale() {
  const params = new URLSearchParams(window.location.search);
  const fromQuery = params.get('lang');
  if (fromQuery && messages[fromQuery]) return fromQuery;
  const fromStorage = localStorage.getItem('go-skills.lang');
  if (fromStorage && messages[fromStorage]) return fromStorage;
  return defaultLocale;
}

export function setLocale(code) {
  if (!messages[code]) return;
  localStorage.setItem('go-skills.lang', code);
  apply(code);
}

export function t(code, key) {
  return (messages[code] && messages[code][key]) || messages[defaultLocale][key] || key;
}

export function apply(code) {
  const locale = supportedLocales.find(l => l.code === code) || supportedLocales[0];
  document.documentElement.setAttribute('lang', code);
  document.documentElement.setAttribute('dir', locale.rtl ? 'rtl' : 'ltr');
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    el.textContent = t(code, key);
  });
  document.querySelectorAll('[data-i18n-html]').forEach(el => {
    const key = el.getAttribute('data-i18n-html');
    el.innerHTML = t(code, key);
  });
  const titleKey = 'meta.title';
  if (messages[code] && messages[code][titleKey]) {
    document.title = messages[code][titleKey];
  }
  // Point GitHub README CTAs to the locale-specific README file.
  const repo = 'https://github.com/amazopic/go-skills';
  const readmeUrl = code === defaultLocale ? repo : `${repo}/blob/main/README.${code}.md`;
  document.querySelectorAll('[data-github-readme]').forEach(el => {
    el.setAttribute('href', readmeUrl);
  });
}
