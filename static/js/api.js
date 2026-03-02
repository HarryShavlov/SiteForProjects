// Получение списка проектов
async function fetchProjects() {
  const response = await fetch(`${API_BASE}/projects`);
  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
  return await response.json();
}

// Вход в систему
async function login(login, password) {
  const response = await fetch(`${API_BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, password })
  });
  
  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.error || 'Ошибка входа');
  }
  
  return await response.json();
}

// Получение данных текущего пользователя
async function getCurrentUser() {
  const response = await fetchWithAuth(`${API_BASE}/me`);
  if (!response.ok) throw new Error('Не авторизован');
  return await response.json();
}

// Загрузка Excel-файла (только админ)
async function uploadProjectsFile(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch(`${API_BASE}/upload`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem(STORAGE_KEYS.TOKEN)}`
    },
    body: formData
  });
  
  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.error || 'Ошибка загрузки');
  }
  
  return await response.json();
}