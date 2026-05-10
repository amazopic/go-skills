// main.js — entrypoint. Initializes i18n + theme.

import { apply, getLocale, setLocale, supportedLocales } from './i18n.js';
import * as themes from './themes.js';

function buildLangSwitcher() {
  const sel = document.createElement('select');
  sel.className = 'lang-switch';
  sel.setAttribute('aria-label', 'Language');
  for (const loc of supportedLocales) {
    const opt = document.createElement('option');
    opt.value = loc.code;
    opt.textContent = `${loc.code.toUpperCase()} · ${loc.native}`;
    sel.appendChild(opt);
  }
  sel.value = getLocale();
  sel.addEventListener('change', () => setLocale(sel.value));
  document.body.appendChild(sel);
}

function init() {
  themes.init();
  apply(getLocale());
  buildLangSwitcher();
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
