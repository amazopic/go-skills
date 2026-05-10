// themes.js — toggle between dark (default, Linear-aligned) and a light flip.
// Stored in localStorage as `go-skills.theme`.

export function getTheme() {
  return localStorage.getItem('go-skills.theme') || 'dark';
}

export function setTheme(name) {
  localStorage.setItem('go-skills.theme', name);
  document.documentElement.setAttribute('data-theme', name);
}

export function init() {
  const initial = getTheme();
  document.documentElement.setAttribute('data-theme', initial);
}
