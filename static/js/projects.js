// Глобальные переменные
let projects = [];
let favorites = new Set(JSON.parse(localStorage.getItem(STORAGE_KEYS.FAVORITES) || '[]'));

// Загрузка и рендеринг проектов
async function reloadData() {
  const btn = document.getElementById('reloadBtn');
  if (btn) {
    btn.classList.add('loading');
    btn.innerHTML = '<span>⏳</span><span>Загрузка...</span>';
  }
  
  try {
    projects = await fetchProjects();
    
    // Преобразуем строки тегов в массивы
    projects = projects.map(p => ({
      ...p,
      tags: typeof p.tags === 'string' ? p.tags.split(',').map(t => t.trim()) : (p.tags || [])
    }));
    
    renderAll();
    showToast('✅ Данные обновлены', 'success');
    
  } catch (error) {
    console.error('Ошибка загрузки:', error);
    showToast('❌ Не удалось загрузить данные', 'error');
    renderEmptyState('popularProjects', 'Не удалось загрузить проекты');
    renderEmptyState('levelProjects', 'Проверьте подключение к серверу');
  } finally {
    if (btn) {
      btn.classList.remove('loading');
      btn.innerHTML = '<span>🔄</span><span>Обновить</span>';
    }
  }
}

// Рендеринг всех секций
function renderAll() {
  const popular = getFilteredProjects().slice(0, 6);
  renderProjects(popular, 'popularProjects');
  
  const levelFiltered = getFilteredProjects().filter(p => {
    const activeLevel = document.querySelector('.level-btn.active')?.dataset.level || 'all';
    return activeLevel === 'all' || p.level === activeLevel;
  });
  renderProjects(levelFiltered, 'levelProjects');
}

// Фильтрация проектов
function getFilteredProjects() {
  let filtered = [...projects];
  
  // Поиск
  const search = document.getElementById('searchInput')?.value.toLowerCase() || '';
  if (search) {
    filtered = filtered.filter(p => 
      p.title?.toLowerCase().includes(search) || 
      p.subtitle?.toLowerCase().includes(search)
    );
  }
  
  // Класс
  const classVal = document.getElementById('classFilter')?.value;
  if (classVal) {
    filtered = filtered.filter(p => p.tags?.some(t => t.includes(classVal + ' класс')));
  }
  
  // Предмет
  const subjectVal = document.getElementById('subjectFilter')?.value;
  if (subjectVal) {
    filtered = filtered.filter(p => p.tags?.some(t => t.includes(subjectVal)));
  }
  
  // Категория
  const category = document.querySelector('.category-tab.active')?.dataset.category;
  if (category && category !== 'all' && category !== 'popular') {
    filtered = filtered.filter(p => p.tags?.some(t => t.includes(category)));
  }
  
  return filtered;
}

// Рендеринг карточек проектов
function renderProjects(projectsToRender, containerId) {
  const container = document.getElementById(containerId);
  
  if (!projectsToRender || projectsToRender.length === 0) {
    renderEmptyState(containerId, 'Ничего не найдено по заданным фильтрам');
    return;
  }
  
  container.innerHTML = projectsToRender.map(project => {
    const isFavorite = favorites.has(project.id);
    
    return `
      <div class="project-card" data-id="${project.id}" onclick="openProject(${project.id})">
        <div class="project-header">
          <div class="level-badge">
            <span>${project.level === 'basic' ? '📚' : '🚀'}</span>
            <span>${project.level === 'basic' ? 'Базовый' : 'Углублённый'}</span>
          </div>
          <h3 class="project-title">${escapeHtml(project.title || 'Без названия')}</h3>
          <p class="project-subtitle">${escapeHtml(project.subtitle || '')}</p>
          <button class="favorite-btn ${isFavorite ? 'active' : ''}" 
                  onclick="toggleFavorite(${project.id}, event)">
            ${isFavorite ? '♥' : '♡'}
          </button>
        </div>
        <div class="project-body">
          <div class="project-tags">
            ${(project.tags || []).slice(0, 3).map(tag => 
              `<span class="tag ${tag === 'Исследовательский' ? 'primary' : ''}">${escapeHtml(tag)}</span>`
            ).join('')}
            ${project.tags?.length > 3 ? `<span class="tag">+${project.tags.length - 3}</span>` : ''}
          </div>
          <div class="project-meta">
            <div class="project-stats">
              <span class="stat">📦 ${escapeHtml(project.orders || '0')}</span>
              <span class="stat">👥 ${escapeHtml(project.participants || '0')}</span>
            </div>
            <div class="project-provider">
              <div class="provider-logo" style="background: ${project.providerColor || '#4F46E5'}">
                ${(project.provider || '?').charAt(0).toUpperCase()}
              </div>
              <span>${escapeHtml(project.provider || 'Неизвестно')}</span>
            </div>
          </div>
        </div>
      </div>
    `;
  }).join('');
}

function renderEmptyState(containerId, message) {
  document.getElementById(containerId).innerHTML = `
    <div class="empty-state" style="grid-column: 1 / -1;">
      <div class="empty-state-icon">🔍</div>
      <p>${message}</p>
    </div>
  `;
}

function renderLoading(containerId) {
  document.getElementById(containerId).innerHTML = `
    <div class="loading-container" style="grid-column: 1 / -1;">
      <div class="spinner"></div>
      <p>Загрузка проектов...</p>
    </div>
  `;
}

function toggleFavorite(id, event) {
  event.stopPropagation();
  
  if (favorites.has(id)) {
    favorites.delete(id);
    showToast('❌ Удалено из избранного');
  } else {
    favorites.add(id);
    showToast('⭐ Добавлено в избранное', 'success');
  }
  
  localStorage.setItem(STORAGE_KEYS.FAVORITES, JSON.stringify([...favorites]));
  renderAll();
}

function openProject(id) {
  const project = projects.find(p => p.id === id);
  if (!project) return;
  
  alert(`📋 ${project.title}\n\n${project.subtitle}\n\nУровень: ${project.level}\nТеги: ${project.tags?.join(', ')}`);
}

// Подписка на события фильтров
function initProjectFilters() {
  document.getElementById('searchInput')?.addEventListener('input', renderAll);
  document.getElementById('classFilter')?.addEventListener('change', renderAll);
  document.getElementById('subjectFilter')?.addEventListener('change', renderAll);
  
  document.querySelectorAll('.category-tab')?.forEach(tab => {
    tab.addEventListener('click', function() {
      document.querySelectorAll('.category-tab').forEach(t => t.classList.remove('active'));
      this.classList.add('active');
      renderAll();
    });
  });
  
  document.querySelectorAll('.level-btn')?.forEach(btn => {
    btn.addEventListener('click', function() {
      document.querySelectorAll('.level-btn').forEach(b => b.classList.remove('active'));
      this.classList.add('active');
      renderAll();
    });
  });
}