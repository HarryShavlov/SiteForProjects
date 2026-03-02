// ==================== ADMIN FUNCTIONS ====================
// Все функции привязаны к window для доступа из HTML

window.openAdminModal = function() {
  if (currentUser?.role !== 'admin') return;
  document.getElementById('adminModal').classList.add('active');
  updateAdminStats();
}

window.closeAdminModal = function() {
  document.getElementById('adminModal').classList.remove('active');
  const statusDiv = document.getElementById('uploadStatus');
  if (statusDiv) {
    statusDiv.className = 'upload-status';
    statusDiv.textContent = '';
  }
}

window.updateAdminStats = function() {
  document.getElementById('totalProjects').textContent = projects?.length || 0;
  document.getElementById('totalUsers').textContent = '—';
}

window.handleUpload = async function(event) {
  event.preventDefault();
  
  const fileInput = document.getElementById('fileInput');
  const file = fileInput?.files[0];
  const statusDiv = document.getElementById('uploadStatus');
  const btn = document.getElementById('uploadBtn');
  
  if (!file) {
    statusDiv.textContent = '❌ Выберите файл';
    statusDiv.className = 'upload-status error';
    return;
  }
  
  if (!file.name.endsWith('.xlsx')) {
    statusDiv.textContent = '❌ Только файлы .xlsx';
    statusDiv.className = 'upload-status error';
    return;
  }
  
  btn.disabled = true;
  btn.textContent = '⏳ Загрузка...';
  statusDiv.textContent = '📤 Отправка файла...';
  statusDiv.className = 'upload-status loading';
  
  try {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch(`${API_BASE}/upload`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: formData
    });
    
    const result = await response.json();
    
    if (!response.ok) {
      throw new Error(result.error || 'Ошибка загрузки');
    }
    
    statusDiv.textContent = `✅ ${result.message}`;
    statusDiv.className = 'upload-status success';
    
    projects = await fetchProjects();
    projects = projects.map(p => ({
      ...p,
      tags: typeof p.tags === 'string' ? p.tags.split(',').map(t => t.trim()) : (p.tags || [])
    }));
    
    renderAll();
    updateAdminStats();
    
    fileInput.value = '';
    
    setTimeout(() => {
      statusDiv.className = 'upload-status';
      statusDiv.textContent = '';
    }, 3000);
    
  } catch (error) {
    console.error('Ошибка загрузки:', error);
    statusDiv.textContent = `❌ ${error.message}`;
    statusDiv.className = 'upload-status error';
  } finally {
    btn.disabled = false;
    btn.textContent = '📤 Загрузить и обновить';
  }
}

window.showAdminPanel = function() {
  if (currentUser?.role !== 'admin') return;
  openAdminModal();
}

window.showAdminFeatures = function() {
  const reloadBtn = document.getElementById('reloadBtn');
  if (reloadBtn) reloadBtn.style.display = 'flex';
  
  // Удаляем старую кнопку перед созданием новой (защита от дублирования)
//   const oldAdminBtn = document.getElementById('adminPanelBtn');
//   if (oldAdminBtn) oldAdminBtn.remove();
  
//   const headerActions = document.querySelector('.header-actions');
//   if (headerActions) {
//     const adminBtn = document.createElement('button');
//     adminBtn.id = 'adminPanelBtn';
//     adminBtn.className = 'btn btn-primary';
//     adminBtn.innerHTML = '⚙️ Панель';
//     adminBtn.onclick = showAdminPanel;
    
//     const logoutBtn = headerActions.querySelector('button:last-child');
//     headerActions.insertBefore(adminBtn, logoutBtn);
//   }
}

window.hideAdminFeatures = function() {
  const reloadBtn = document.getElementById('reloadBtn');
  if (reloadBtn) reloadBtn.style.display = 'none';
  const adminBtn = document.getElementById('adminPanelBtn');
  if (adminBtn) adminBtn.remove();
}