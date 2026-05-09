const PAGE_SIZE = 10;
let currentPage = 'products';
let products = [], brands = [], categories = [];
let productPage = 1;

const $ = s => document.querySelector(s);
const $$ = s => document.querySelectorAll(s);

async function api(url, opts = {}) {
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json', ...opts.headers },
    ...opts,
  });
  if (!res.ok) {
    const text = await res.text();
    let msg;
    try { msg = JSON.parse(text).error; } catch { msg = text; }
    throw new Error(msg || `HTTP ${res.status}`);
  }
  const ct = res.headers.get('content-type') || '';
  if (ct.includes('application/json')) return res.json();
  return res.text();
}

// Navigation
$$('.nav-btn').forEach(btn => {
  btn.onclick = () => {
    $$('.nav-btn').forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    currentPage = btn.dataset.page;
    productPage = 1;
    render();
  };
});

function render() {
  if (currentPage === 'products') renderProducts();
  else if (currentPage === 'brands') renderBrands();
  else if (currentPage === 'categories') renderCategories();
}

// Modal
function openModal(title, content) {
  $('#modal-title').textContent = title;
  $('#modal-content').innerHTML = content;
  $('#modal').classList.remove('hidden');
}

function closeModal() {
  $('#modal').classList.add('hidden');
}

$('.modal-close').onclick = closeModal;
$('.modal-backdrop').onclick = closeModal;

// ─── LOAD ────────────────────────────────────────────────
async function loadAll() {
  try {
    [products, brands, categories] = await Promise.all([
      api('/products'), api('/brands'), api('/categories'),
    ]);
    render();
  } catch (e) {
    console.error(e);
  }
}

// ─── PRODUCTS ────────────────────────────────────────────
function renderProducts() {
  const totalPages = Math.ceil(products.length / PAGE_SIZE) || 1;
  const start = (productPage - 1) * PAGE_SIZE;
  const page = products.slice(start, start + PAGE_SIZE);
  const cards = page.map(p => `
    <div class="product-card" onclick="editProduct(${p.ID})">
      <div class="card-img">
        ${p.MainPhotoURL ? `<img src="${esc(p.MainPhotoURL)}" alt="" onerror="this.style.display='none';this.nextElementSibling.style.display='flex'" loading="lazy"><div class="card-img-placeholder" style="display:none">📷</div>` : `<div class="card-img-placeholder">📷</div>`}
      </div>
      <div class="card-body">
        <div class="card-title">${esc(p.Title)}</div>
        <div class="card-brand">${p.Brand && p.Brand.Title ? esc(p.Brand.Title) : '—'}</div>
        <div class="card-price">${p.Price.toLocaleString()}₽</div>
      </div>
      <div class="card-actions" onclick="event.stopPropagation()">
        <button class="btn btn-sm btn-primary" onclick="editProduct(${p.ID})">✎</button>
        <button class="btn btn-sm btn-danger" onclick="deleteProduct(${p.ID})">✕</button>
      </div>
    </div>`).join('');

  let pagesHtml = '';
  for (let i = 1; i <= totalPages; i++) {
    pagesHtml += `<button class="page-btn${i === productPage ? ' active' : ''}" onclick="productPage=${i};renderProducts()">${i}</button>`;
  }

  const searchBar = `
    <div class="search-bar">
      <input class="form-control" placeholder="Поиск по названию..." oninput="filterProducts(this.value)" id="search-input">
    </div>`;

  $('#page-content').innerHTML = searchBar + `
    <div class="product-grid">${cards || '<div class="empty">Нет товаров</div>'}</div>
    ${totalPages > 1 ? `<div class="pagination">${pagesHtml}</div>` : ''}`;
  $('#page-title').textContent = 'Товары';
  $('#add-btn').onclick = showAddProduct;
  $('#add-btn').style.display = '';
}

function filterProducts(q) {
  if (!q) { loadAll(); return; }
  const filtered = products.filter(p => p.Title.toLowerCase().includes(q.toLowerCase()));
  const cards = filtered.map(p => `
    <div class="product-card" onclick="editProduct(${p.ID})">
      <div class="card-img">
        ${p.MainPhotoURL ? `<img src="${esc(p.MainPhotoURL)}" alt="" onerror="this.style.display='none';this.nextElementSibling.style.display='flex'" loading="lazy"><div class="card-img-placeholder" style="display:none">📷</div>` : `<div class="card-img-placeholder">📷</div>`}
      </div>
      <div class="card-body">
        <div class="card-title">${esc(p.Title)}</div>
        <div class="card-brand">${p.Brand && p.Brand.Title ? esc(p.Brand.Title) : '—'}</div>
        <div class="card-price">${p.Price.toLocaleString()}₽</div>
      </div>
      <div class="card-actions" onclick="event.stopPropagation()">
        <button class="btn btn-sm btn-primary" onclick="editProduct(${p.ID})">✎</button>
        <button class="btn btn-sm btn-danger" onclick="deleteProduct(${p.ID})">✕</button>
      </div>
    </div>`).join('');
  document.querySelector('#page-content .product-grid').innerHTML = cards || '<div class="empty">Нет совпадений</div>';
}

async function showAddProduct() {
  const brandOpts = '<option value="">— Без бренда —</option>' + brands.map(b => `<option value="${b.ID}">${esc(b.Title)}</option>`).join('');
  const catOpts = categories.map(c => `<label><input type="checkbox" value="${c.ID}"> ${esc(c.Title)}</label>`).join('');

  openModal('Добавить товар', `
    <div class="form-group"><label>Название</label><input class="form-control" id="f-title"></div>
    <div class="form-row">
      <div class="form-group">
        <label>Бренд</label>
        <div style="display:flex;gap:6px;align-items:center">
          <select class="form-control" id="f-brand" style="flex:1">${brandOpts}</select>
          <button type="button" class="btn btn-sm btn-primary" id="f-brand-add-btn" onclick="showNewBrandInput()" title="Создать новый бренд">+</button>
        </div>
        <div id="f-brand-new" style="display:none;margin-top:8px">
          <input class="form-control" id="f-brand-new-title" placeholder="Название нового бренда">
          <div style="margin-top:4px;font-size:12px"><a href="#" onclick="cancelNewBrand();return false" style="color:var(--text2)">Отмена</a></div>
        </div>
      </div>
      <div class="form-group"><label>Цена (₽)</label><input type="number" class="form-control" id="f-price"></div>
    </div>
    <div class="form-group"><label>Описание</label><textarea class="form-control" id="f-desc"></textarea></div>
    <div class="form-group"><label>Категории</label><div class="checkbox-group">${catOpts || '<span class="empty">Нет категорий</span>'}</div></div>
    <div class="form-group">
      <label>Фото</label>
      <div class="photo-source-tabs">
        <button type="button" class="tab-btn active" onclick="switchPhotoTab('file',this)">Файл</button>
        <button type="button" class="tab-btn" onclick="switchPhotoTab('url',this)">Ссылка</button>
      </div>
      <div id="f-photo-file">
        <div style="display:flex;gap:10px;align-items:center">
          <button type="button" class="btn btn-primary" onclick="document.getElementById('f-photo').click()">Выбрать файл</button>
          <input type="file" id="f-photo" accept="image/*" style="display:none" onchange="document.getElementById('f-photo-name').textContent=this.files[0]?.name||''">
          <span id="f-photo-name" style="font-size:13px;color:var(--text2)"></span>
        </div>
      </div>
      <div id="f-photo-url" style="display:none">
        <input type="text" class="form-control" id="f-photo-url-input" placeholder="https://example.com/photo.jpg">
      </div>
    </div>
    <div class="form-actions">
      <button class="btn" onclick="closeModal()">Отмена</button>
      <button class="btn btn-primary" onclick="createProduct()">Сохранить</button>
    </div>`);
}

async function createProduct() {
  const title = $('#f-title').value.trim();
  if (!title) return alert('Название обязательно');
  let productId;
  try {
    let brandId = null;
    const newBrandDiv = $('#f-brand-new');
    if (newBrandDiv.style.display !== 'none') {
      const brandTitle = $('#f-brand-new-title').value.trim();
      if (!brandTitle) { alert('Введите название нового бренда'); return; }
      const res = await api('/brands', { method: 'POST', body: JSON.stringify({ title: brandTitle }) });
      brandId = res.id;
      brands.push(res);
    } else {
      const brandEl = $('#f-brand');
      brandId = brandEl.value ? parseInt(brandEl.value) : null;
    }
    const data = {
      title,
      description: $('#f-desc').value.trim(),
      price: parseInt($('#f-price').value) || 0,
      brand_id: brandId,
      category_ids: [...document.querySelectorAll('#modal-content input[type=checkbox]:checked')].map(c => parseInt(c.value)),
    };
    const res = await api('/products', { method: 'POST', body: JSON.stringify(data) });
    productId = res.id;

    const photoInput = $('#f-photo');
    if (photoInput && photoInput.files[0]) {
      const fd = new FormData();
      fd.append('files', photoInput.files[0]);
      fd.append('product_id', productId);
      fd.append('is_main', 'true');
      await fetch('/upload', { method: 'POST', body: fd });
    }
    const photoURL = $('#f-photo-url-input');
    if (photoURL && photoURL.value.trim()) {
      await api(`/product/${productId}/photo-url`, { method: 'POST', body: JSON.stringify({ url: photoURL.value.trim() }) });
    }

    closeModal();

    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

async function editProduct(id) {
  const p = products.find(x => x.ID === id);
  if (!p) return;
  const brandOpts = '<option value="">— Без бренда —</option>' + brands.map(b => `<option value="${b.ID}"${p.Brand && p.Brand.ID === b.ID ? ' selected' : ''}>${esc(b.Title)}</option>`).join('');
  const catOpts = categories.map(c => {
    const checked = p.Categories && p.Categories.some(pc => pc.ID === c.ID);
    return `<label><input type="checkbox" value="${c.ID}"${checked ? ' checked' : ''}> ${esc(c.Title)}</label>`;
  }).join('');

  // Load photos
  let photosHtml = '<div class="empty">Загрузка...</div>';
  try {
    const photos = await api(`/product/${id}/photos`);
    photosHtml = photos.map(ph => `
      <div class="photo-item">
        <img src="${ph.URL}" alt="" style="object-fit:cover;width:100%;height:100%">
        ${ph.IsMain ? '<span class="photo-badge">Главное</span>' : ''}
        <div class="photo-actions">
          ${!ph.IsMain ? `<button style="background:var(--accent);color:#fff" onclick="setMainPhoto(${ph.ID},${id})">★</button>` : ''}
          <button style="background:var(--danger);color:#fff" onclick="deletePhoto(${ph.ID})">✕</button>
        </div>
      </div>`).join('');
  } catch { photosHtml = '<div class="empty">Ошибка загрузки</div>'; }

  openModal('Редактировать товар', `
    <input type="hidden" id="f-eid" value="${id}">
    <div class="form-group"><label>Название</label><input class="form-control" id="f-title" value="${esc(p.Title)}"></div>
    <div class="form-row">
      <div class="form-group"><label>Бренд</label>
      <div style="display:flex;gap:6px;align-items:center">
        <select class="form-control" id="f-brand" style="flex:1">${brandOpts}</select>
        <button type="button" class="btn btn-sm btn-primary" id="f-brand-add-btn" onclick="showNewBrandInput()" title="Создать новый бренд">+</button>
      </div>
      <div id="f-brand-new" style="display:none;margin-top:8px">
        <input class="form-control" id="f-brand-new-title" placeholder="Название нового бренда">
        <div style="margin-top:4px;font-size:12px"><a href="#" onclick="cancelNewBrand();return false" style="color:var(--text2)">Отмена</a></div>
      </div>
    </div>
      <div class="form-group"><label>Цена (₽)</label><input type="number" class="form-control" id="f-price" value="${p.Price}"></div>
    </div>
    <div class="form-group"><label>Описание</label><textarea class="form-control" id="f-desc">${esc(p.Description || '')}</textarea></div>
    <div class="form-group"><label>Категории</label><div class="checkbox-group">${catOpts || '<span class="empty">Нет категорий</span>'}</div></div>
    <div class="form-group">
      <label>Фотографии</label>
      <div class="photo-grid" id="photo-grid">${photosHtml}</div>
      <div style="margin-top:8px">
        <button type="button" class="btn btn-sm btn-primary" onclick="document.getElementById('photo-upload-input').click()">+ Загрузить файл</button>
        <input type="file" accept="image/*" style="display:none" id="photo-upload-input" onchange="uploadPhoto(${id})">
      </div>
      <div style="margin-top:8px;display:flex;gap:8px">
        <input type="text" class="form-control" id="f-photo-url-input-edit" placeholder="https://example.com/photo.jpg" style="flex:1">
        <button type="button" class="btn btn-sm btn-primary" onclick="savePhotoURL(${id})">+ По ссылке</button>
      </div>
    </div>
    <div class="form-actions">
      <button class="btn" onclick="closeModal()">Отмена</button>
      <button class="btn btn-primary" onclick="updateProduct()">Сохранить</button>
    </div>`);
}

async function updateProduct() {
  const id = $('#f-eid').value;
  const data = {};
  const title = $('#f-title').value.trim();
  if (!title) return alert('Название обязательно');
  data.title = title;
  const desc = $('#f-desc').value.trim();
  if (desc) data.description = desc;
  const price = parseInt($('#f-price').value);
  if (price) data.price = price;
  try {
    const newBrandDiv = $('#f-brand-new');
    if (newBrandDiv.style.display !== 'none') {
      const brandTitle = $('#f-brand-new-title').value.trim();
      if (!brandTitle) { alert('Введите название нового бренда'); return; }
      const res = await api('/brands', { method: 'POST', body: JSON.stringify({ title: brandTitle }) });
      data.brand_id = res.id;
      brands.push(res);
    } else {
      const brandEl = $('#f-brand');
      if (brandEl.value) data.brand_id = parseInt(brandEl.value);
    }
    const catIds = [...document.querySelectorAll('#modal-content input[type=checkbox]:checked')].map(c => parseInt(c.value));
    data.category_ids = catIds;
    await api(`/product/${id}`, { method: 'PUT', body: JSON.stringify(data) });
    closeModal();
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

async function deleteProduct(id) {
  if (!confirm('Удалить товар?')) return;
  try {
    await api(`/product/${id}`, { method: 'DELETE' });
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

// Photos
async function uploadPhoto(productId) {
  const input = $('#photo-upload-input');
  const file = input.files[0];
  if (!file) return;
  const fd = new FormData();
  fd.append('files', file);
  fd.append('product_id', productId);
  fd.append('is_main', 'false');
  try {
    await fetch('/upload', { method: 'POST', body: fd });
    input.value = '';
    editProduct(productId);
  } catch (e) { alert('Ошибка загрузки: ' + e.message); }
}

async function setMainPhoto(photoId, productId) {
  try {
    await api(`/photo/${photoId}/main/${productId}`, { method: 'PUT' });
    editProduct(productId);
  } catch (e) { alert('Ошибка: ' + e.message); }
}

async function deletePhoto(photoId) {
  if (!confirm('Удалить фото?')) return;
  try {
    await api(`/photo/${photoId}`, { method: 'DELETE' });
    // Refresh the product edit view
    const id = $('#f-eid')?.value;
    if (id) editProduct(parseInt(id));
  } catch (e) { alert('Ошибка: ' + e.message); }
}

function switchPhotoTab(tab, btn) {
  const parent = btn.closest('.form-group');
  parent.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
  btn.classList.add('active');
  parent.querySelector('#f-photo-file').style.display = tab === 'file' ? '' : 'none';
  parent.querySelector('#f-photo-url').style.display = tab === 'url' ? '' : 'none';
}

function showNewBrandInput() {
  $('#f-brand').style.display = 'none';
  $('#f-brand-add-btn').style.display = 'none';
  $('#f-brand-new').style.display = '';
  $('#f-brand-new-title').focus();
}

function cancelNewBrand() {
  $('#f-brand').style.display = '';
  $('#f-brand-add-btn').style.display = '';
  $('#f-brand-new').style.display = 'none';
  $('#f-brand-new-title').value = '';
}

async function savePhotoURL(productId) {
  const input = $(`#f-photo-url-input-edit`);
  const url = input?.value?.trim();
  if (!url) return alert('Введите ссылку на фото');
  try {
    await api(`/product/${productId}/photo-url`, { method: 'POST', body: JSON.stringify({ url }) });
    input.value = '';
    editProduct(productId);
  } catch (e) { alert('Ошибка: ' + e.message); }
}

// ─── BRANDS ──────────────────────────────────────────────
function renderBrands() {
  const cards = brands.map(b => `
    <div class="mini-card" onclick="editBrand(${b.ID})">
      <div class="mini-card-icon">🏢</div>
      <div class="mini-card-body">
        <div class="mini-card-title">${esc(b.Title)}</div>
        <div class="mini-card-desc">${esc(b.Description || 'Нет описания')}</div>
      </div>
      <div class="mini-card-actions" onclick="event.stopPropagation()">
        <button class="btn btn-sm btn-primary" onclick="editBrand(${b.ID})">✎</button>
        <button class="btn btn-sm btn-danger" onclick="deleteBrand(${b.ID})">✕</button>
      </div>
    </div>`).join('');

  $('#page-content').innerHTML = `
    <div class="mini-grid">${cards || '<div class="empty">Нет брендов</div>'}</div>`;
  $('#page-title').textContent = 'Бренды';
  $('#add-btn').onclick = showAddBrand;
  $('#add-btn').style.display = '';
}

function showAddBrand() {
  openModal('Добавить бренд', `
    <div class="form-group"><label>Название</label><input class="form-control" id="f-title"></div>
    <div class="form-group"><label>Описание</label><textarea class="form-control" id="f-desc"></textarea></div>
    <div class="form-actions">
      <button class="btn" onclick="closeModal()">Отмена</button>
      <button class="btn btn-primary" onclick="createBrand()">Сохранить</button>
    </div>`);
}

async function createBrand() {
  const title = $('#f-title').value.trim();
  if (!title) return alert('Название обязательно');
  try {
    await api('/brands', { method: 'POST', body: JSON.stringify({ title, description: $('#f-desc').value.trim() || null }) });
    closeModal();
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

async function editBrand(id) {
  const b = brands.find(x => x.ID === id);
  if (!b) return;
  openModal('Редактировать бренд', `
    <input type="hidden" id="f-eid" value="${id}">
    <div class="form-group"><label>Название</label><input class="form-control" id="f-title" value="${esc(b.Title)}"></div>
    <div class="form-group"><label>Описание</label><textarea class="form-control" id="f-desc">${esc(b.Description || '')}</textarea></div>
    <div class="form-actions">
      <button class="btn" onclick="closeModal()">Отмена</button>
      <button class="btn btn-primary" onclick="updateBrand()">Сохранить</button>
    </div>`);
}

async function updateBrand() {
  const id = $('#f-eid').value;
  const data = {};
  const title = $('#f-title').value.trim();
  if (title) data.title = title;
  const desc = $('#f-desc').value.trim();
  if (desc) data.description = desc;
  try {
    await api(`/brand/${id}`, { method: 'PUT', body: JSON.stringify(data) });
    closeModal();
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

async function deleteBrand(id) {
  if (!confirm('Удалить бренд?')) return;
  try {
    await api(`/brand/${id}`, { method: 'DELETE' });
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

// ─── CATEGORIES ──────────────────────────────────────────
function renderCategories() {
  const cards = categories.map(c => `
    <div class="mini-card" onclick="editCategory(${c.ID})">
      <div class="mini-card-icon">📂</div>
      <div class="mini-card-body">
        <div class="mini-card-title">${esc(c.Title)}</div>
      </div>
      <div class="mini-card-actions" onclick="event.stopPropagation()">
        <button class="btn btn-sm btn-primary" onclick="editCategory(${c.ID})">✎</button>
        <button class="btn btn-sm btn-danger" onclick="deleteCategory(${c.ID})">✕</button>
      </div>
    </div>`).join('');

  $('#page-content').innerHTML = `
    <div class="mini-grid">${cards || '<div class="empty">Нет категорий</div>'}</div>`;
  $('#page-title').textContent = 'Категории';
  $('#add-btn').onclick = showAddCategory;
  $('#add-btn').style.display = '';
}

function showAddCategory() {
  openModal('Добавить категорию', `
    <div class="form-group"><label>Название</label><input class="form-control" id="f-title"></div>
    <div class="form-actions">
      <button class="btn" onclick="closeModal()">Отмена</button>
      <button class="btn btn-primary" onclick="createCategory()">Сохранить</button>
    </div>`);
}

async function createCategory() {
  const title = $('#f-title').value.trim();
  if (!title) return alert('Название обязательно');
  try {
    await api('/categories', { method: 'POST', body: JSON.stringify({ title }) });
    closeModal();
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

async function editCategory(id) {
  const c = categories.find(x => x.ID === id);
  if (!c) return;
  openModal('Редактировать категорию', `
    <input type="hidden" id="f-eid" value="${id}">
    <div class="form-group"><label>Название</label><input class="form-control" id="f-title" value="${esc(c.Title)}"></div>
    <div class="form-actions">
      <button class="btn" onclick="closeModal()">Отмена</button>
      <button class="btn btn-primary" onclick="updateCategory()">Сохранить</button>
    </div>`);
}

async function updateCategory() {
  const id = $('#f-eid').value;
  const title = $('#f-title').value.trim();
  if (!title) return alert('Название обязательно');
  try {
    await api(`/category/${id}`, { method: 'PUT', body: JSON.stringify({ title }) });
    closeModal();
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

async function deleteCategory(id) {
  if (!confirm('Удалить категорию?')) return;
  try {
    await api(`/category/${id}`, { method: 'DELETE' });
    await loadAll();
  } catch (e) { alert('Ошибка: ' + e.message); }
}

// ─── UTILS ───────────────────────────────────────────────
function esc(s) {
  if (!s) return '';
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}

// ─── INIT ────────────────────────────────────────────────
loadAll();
