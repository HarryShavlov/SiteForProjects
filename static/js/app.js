// Инициализация приложения
document.addEventListener('DOMContentLoaded', async () => {
  // Инициализация фильтров
  if (typeof initProjectFilters === 'function') {
    initProjectFilters();
  }
  
  // Показываем лоадеры
  if (typeof renderLoading === 'function') {
    renderLoading('popularProjects');
    renderLoading('levelProjects');
  }
  
  // Инициализация авторизации
  await initAuth();
  
  // Загрузка проектов
  if (typeof reloadData === 'function') {
    await reloadData();
  }
});