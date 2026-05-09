// Theme toggle
function toggleTheme() {
  const html = document.documentElement;
  const current = html.getAttribute('data-theme');
  const next = current === 'dark' ? 'light' : 'dark';
  html.setAttribute('data-theme', next);
  localStorage.setItem('theme', next);
}

// Sidebar mobile
function toggleSidebar() {
  document.getElementById('sidebar-mobile').classList.toggle('open');
  document.getElementById('overlay').classList.toggle('show');
}

// Apply saved theme on load
(function() {
  const saved = localStorage.getItem('theme') || 'light';
  document.documentElement.setAttribute('data-theme', saved);
})();
