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

function setupCopyButtons() {
  document.querySelectorAll('[data-copy-target]').forEach(btn => {
    btn.addEventListener('click', async () => {
      const target = document.querySelector(btn.getAttribute('data-copy-target'));
      if (!target) return;
      try {
        await navigator.clipboard.writeText(target.textContent.trim());
        const label = btn.querySelector('span') || btn;
        const orig = label.textContent;
        label.textContent = btn.getAttribute('data-copy-done') || 'Copied';
        btn.classList.add('copied');
        setTimeout(() => {
          label.textContent = orig;
          btn.classList.remove('copied');
        }, 1500);
      } catch (e) {
        console.error('Copy failed', e);
      }
    });
  });
}

function init() {
  themes.init();
  apply(getLocale());
  buildLangSwitcher();
  setupCopyButtons();
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
