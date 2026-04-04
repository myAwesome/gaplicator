import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom';
import PostPage from './pages/PostPage';
function ThemeToggle() {
  const [dark, setDark] = useState(() => {
    const stored = localStorage.getItem('theme');
    if (stored === 'dark') return true;
    if (stored === 'light') return false;
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
    localStorage.setItem('theme', dark ? 'dark' : 'light');
  }, [dark]);

  return (
    <button
      className="btn-theme-toggle"
      onClick={() => setDark(d => !d)}
      title={dark ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {dark ? '☀' : '☾'}
    </button>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <nav className="app-nav">
        <NavLink to="/posts" className={({ isActive }) => isActive ? 'active' : ''}>Post</NavLink>
        <span className="nav-spacer" />
        <ThemeToggle />
      </nav>
      <main className="app-main">
        <Routes>
          <Route path="/posts" element={<PostPage />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

