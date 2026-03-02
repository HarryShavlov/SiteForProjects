let currentUser = null;

// Открыть модальное окно
window.openModal = function() {
  document.getElementById('loginModal').classList.add('active');
}

// Закрыть модальное окно
window.closeModal = function() {
  document.getElementById('loginModal').classList.remove('active');
  document.getElementById('loginError').textContent = '';
}

// Обработка формы входа
window.handleLogin = async function(event) {
  event.preventDefault();
  
  const loginInput = document.getElementById('loginInput').value;
  const password = document.getElementById('passwordInput').value;
  const btn = document.getElementById('loginBtn');
  const errorDiv = document.getElementById('loginError');
  
  btn.disabled = true;
  btn.textContent = 'Вход...';
  errorDiv.textContent = '';
  
  try {
    const data = await login(loginInput, password);
    
    localStorage.setItem(STORAGE_KEYS.TOKEN, data.token);
    currentUser = data.user;
    
    showToast('✅ Успешный вход!', 'success');
    closeModal();
    updateHeader(currentUser);
    
    if (currentUser.role === 'admin') {
      showAdminFeatures();
    }
    
    if (typeof renderAll === 'function') {
      renderAll();
    }
    
  } catch (error) {
    errorDiv.textContent = error.message;
  } finally {
    btn.disabled = false;
    btn.textContent = 'Войти';
  }
}

// Выход из системы
window.logout = function() {
  localStorage.removeItem(STORAGE_KEYS.TOKEN);
  currentUser = null;
  showToast('👋 Вы вышли из системы');
  updateHeader(null);
  hideAdminFeatures();
  openModal();
}

// Обновление хедера
window.updateHeader = function(user) {
  const container = document.getElementById('headerActions');
  
  if (user) {
    container.innerHTML = `
      <span style="color:white;opacity:0.9">
        👤 ${escapeHtml(user.name)} 
        ${user.class ? `(${escapeHtml(user.class)})` : ''}
        <span class="role-badge ${user.role}">
          ${user.role === 'admin' ? '👨‍ Админ' : '👨‍ Ученик'}
        </span>
      </span>
      <button class="btn btn-outline" onclick="logout()">Выйти</button>
      ${user.role === 'admin' ? '<button class="btn btn-primary" onclick="showAdminPanel()">⚙️ Панель</button>' : ''}
    `;
  } else {
    container.innerHTML = `
      <button class="btn btn-outline" onclick="openModal()">🔐 Войти</button>
      <button class="btn btn-primary" onclick="alert('Регистрация пока недоступна')">📝 Регистрация</button>
    `;
  }
}

// Инициализация авторизации
async function initAuth() {
  const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
  
  if (token) {
    try {
      currentUser = await getCurrentUser();
      updateHeader(currentUser);
      if (currentUser.role === 'admin') {
        showAdminFeatures();
      }
      return true;
    } catch {
      localStorage.removeItem(STORAGE_KEYS.TOKEN);
    }
  }
  
  openModal();
  return false;
}