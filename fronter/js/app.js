// 医院挂号系统主应用
document.addEventListener('DOMContentLoaded', function() {
    // 初始化
    initApp();

    // 设置当前日期
    setCurrentDate();

    // 绑定导航点击事件
    setupNavigation();

    // 加载首页统计数据
    loadDashboardStats();

    // 设置搜索功能
    setupSearch();

    // 初始化模态框
    initModals();
});

// API基础URL
const API_BASE_URL = '/api';
const DEPARTMENTS = ['内科', '外科', '儿科', '妇产科', '眼科', '耳鼻喉科', '口腔科', '皮肤科', '中医科', '其他'];

let currentPatients = [];
let currentDiseases = [];
let currentDoctors = [];
let currentRegistrations = [];

let editingPatientId = null;
let editingDiseaseId = null;
let editingDoctorId = null;
let editingRegistrationId = null;

// 初始化应用
function initApp() {
    console.log('医院挂号系统前端应用初始化');

    // 设置默认页面
    if (!window.location.hash) {
        window.location.hash = '#dashboard';
    }

    // 监听hash变化
    window.addEventListener('hashchange', function() {
        const page = window.location.hash.substring(1) || 'dashboard';
        showPage(page);
    });

    // 初始显示页面
    const initialPage = window.location.hash.substring(1) || 'dashboard';
    showPage(initialPage);
}

// 设置当前日期
function setCurrentDate() {
    const now = new Date();
    const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'long'
    };
    const dateString = now.toLocaleDateString('zh-CN', options);
    document.getElementById('current-date').textContent = dateString;
}

// 设置导航
function setupNavigation() {
    const navItems = document.querySelectorAll('.nav-item');

    navItems.forEach(item => {
        item.addEventListener('click', function() {
            // 移除所有active类
            navItems.forEach(nav => nav.classList.remove('active'));

            // 添加active类到当前项
            this.classList.add('active');

            // 显示对应页面
            const page = this.getAttribute('data-page');
            window.location.hash = page;
        });
    });
}

// 显示页面
function showPage(pageId) {
    // 隐藏所有页面
    const pages = document.querySelectorAll('.page');
    pages.forEach(page => {
        page.classList.remove('active');
    });

    // 显示目标页面
    const targetPage = document.getElementById(pageId);
    if (targetPage) {
        targetPage.classList.add('active');

        // 更新导航
        const navItems = document.querySelectorAll('.nav-item');
        navItems.forEach(item => {
            if (item.getAttribute('data-page') === pageId) {
                item.classList.add('active');
            } else {
                item.classList.remove('active');
            }
        });

        // 加载页面数据
        loadPageData(pageId);
    }
}

// 加载页面数据
function loadPageData(pageId) {
    switch(pageId) {
        case 'dashboard':
            loadDashboardStats();
            break;
        case 'patients':
            loadPatients();
            break;
        case 'diseases':
            loadDiseases();
            break;
        case 'doctors':
            loadDoctors();
            break;
        case 'registrations':
            loadRegistrations();
            break;
        case 'reports':
            loadReports();
            break;
    }
}

// 加载首页统计数据
async function loadDashboardStats() {
    try {
        // 模拟数据，实际应该从API获取
        document.getElementById('patient-count').textContent = '1,234';
        document.getElementById('doctor-count').textContent = '56';
        document.getElementById('disease-count').textContent = '89';
        document.getElementById('registration-count').textContent = '45';

        // 实际API调用示例
        /*
        const responses = await Promise.all([
            fetch(`${API_BASE_URL}/patients`),
            fetch(`${API_BASE_URL}/doctors`),
            fetch(`${API_BASE_URL}/diseases`),
            fetch(`${API_BASE_URL}/registrations/today`)
        ]);

        const data = await Promise.all(responses.map(r => r.json()));

        document.getElementById('patient-count').textContent = data[0].length;
        document.getElementById('doctor-count').textContent = data[1].length;
        document.getElementById('disease-count').textContent = data[2].length;
        document.getElementById('registration-count').textContent = data[3].length;
        */
    } catch (error) {
        console.error('加载统计数据失败:', error);
    }
}

// 设置搜索功能
function setupSearch() {
    const searchInput = document.querySelector('.search-bar input');
    const patientSearch = document.getElementById('patient-search');
    const patientFilter = document.getElementById('patient-filter');
    const registrationStatusFilter = document.getElementById('registration-status-filter');
    const registrationDateFilter = document.getElementById('registration-date-filter');

    if (searchInput) {
        searchInput.addEventListener('input', function(e) {
            const query = e.target.value.toLowerCase();
            // 根据当前页面执行搜索
            const currentPage = window.location.hash.substring(1);
            searchInPage(currentPage, query);
        });
    }

    if (patientSearch) {
        patientSearch.addEventListener('input', function(e) {
            const query = e.target.value.toLowerCase();
            filterPatients(query);
        });
    }

    if (patientFilter) {
        patientFilter.addEventListener('change', function() {
            applyPatientsFilterAndRender();
        });
    }

    if (registrationStatusFilter) {
        registrationStatusFilter.addEventListener('change', function() {
            applyRegistrationFilterAndRender();
        });
    }

    if (registrationDateFilter) {
        registrationDateFilter.addEventListener('change', function() {
            applyRegistrationFilterAndRender();
        });
    }
}

// 页面内搜索
function searchInPage(page, query) {
    // 根据页面执行不同搜索逻辑
    console.log(`在${page}页面搜索: ${query}`);
}

// 初始化模态框
function initModals() {
    // 绑定添加按钮事件
    document.getElementById('add-patient-btn')?.addEventListener('click', showAddPatientModal);
    document.getElementById('add-disease-btn')?.addEventListener('click', showAddDiseaseModal);
    document.getElementById('add-doctor-btn')?.addEventListener('click', showAddDoctorModal);
    document.getElementById('add-registration-btn')?.addEventListener('click', showAddRegistrationModal);
}

// 显示添加病人模态框
function showAddPatientModal() {
    editingPatientId = null;
    const modalHtml = `
        <div class="modal active" id="add-patient-modal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>新增病人</h2>
                    <button class="btn-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="patient-form">
                        <div class="form-row">
                            <div class="form-group">
                                <label for="patient-name">姓名 *</label>
                                <input type="text" id="patient-name" required>
                            </div>
                            <div class="form-group">
                                <label for="patient-gender">性别 *</label>
                                <select id="patient-gender" required>
                                    <option value="">请选择</option>
                                    <option value="男">男</option>
                                    <option value="女">女</option>
                                </select>
                            </div>
                        </div>
                        
                        <div class="form-row">
                            <div class="form-group">
                                <label for="patient-age">年龄 *</label>
                                <input type="number" id="patient-age" min="0" max="120" required>
                            </div>
                            <div class="form-group">
                                <label for="patient-phone">电话 *</label>
                                <input type="tel" id="patient-phone" required>
                            </div>
                        </div>
                        
                        <div class="form-group">
                            <label for="patient-idcard">身份证号 *</label>
                            <input type="text" id="patient-idcard" required>
                        </div>
                        
                        <div class="form-group">
                            <label for="patient-address">地址</label>
                            <input type="text" id="patient-address">
                        </div>
                        
                        <div class="form-row">
                            <div class="form-group">
                                <label for="patient-emergency-contact">紧急联系人</label>
                                <input type="text" id="patient-emergency-contact">
                            </div>
                            <div class="form-group">
                                <label for="patient-emergency-phone">紧急联系电话</label>
                                <input type="tel" id="patient-emergency-phone">
                            </div>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary" id="cancel-patient-btn">取消</button>
                    <button class="btn-primary" id="save-patient-btn">保存</button>
                </div>
            </div>
        </div>
    `;

    document.getElementById('modal-container').innerHTML = modalHtml;

    // 绑定事件
    document.querySelector('#add-patient-modal .btn-close').addEventListener('click', closeModal);
    document.getElementById('cancel-patient-btn').addEventListener('click', closeModal);
    document.getElementById('save-patient-btn').addEventListener('click', savePatient);
}

// 显示添加病种模态框
function showAddDiseaseModal() {
    editingDiseaseId = null;
    const modalHtml = `
        <div class="modal active" id="add-disease-modal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>新增病种</h2>
                    <button class="btn-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="disease-form">
                        <div class="form-group">
                            <label for="disease-name">病种名称 *</label>
                            <input type="text" id="disease-name" required>
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-category">分类 *</label>
                            <select id="disease-category" required>
                                <option value="">请选择</option>
                                <option value="内科">内科</option>
                                <option value="外科">外科</option>
                                <option value="儿科">儿科</option>
                                <option value="妇产科">妇产科</option>
                                <option value="眼科">眼科</option>
                                <option value="耳鼻喉科">耳鼻喉科</option>
                                <option value="口腔科">口腔科</option>
                                <option value="皮肤科">皮肤科</option>
                                <option value="中医科">中医科</option>
                            </select>
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-description">描述</label>
                            <textarea id="disease-description" rows="3"></textarea>
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-symptoms">症状</label>
                            <textarea id="disease-symptoms" rows="3"></textarea>
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-treatment">治疗方法</label>
                            <textarea id="disease-treatment" rows="3"></textarea>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary" id="cancel-disease-btn">取消</button>
                    <button class="btn-primary" id="save-disease-btn">保存</button>
                </div>
            </div>
        </div>
    `;

    document.getElementById('modal-container').innerHTML = modalHtml;

    // 绑定事件
    document.querySelector('#add-disease-modal .btn-close').addEventListener('click', closeModal);
    document.getElementById('cancel-disease-btn').addEventListener('click', closeModal);
    document.getElementById('save-disease-btn').addEventListener('click', saveDisease);
}

// 关闭模态框
function closeModal() {
    document.getElementById('modal-container').innerHTML = '';
}

// 保存病人
async function savePatient() {
    const patient = {
        name: document.getElementById('patient-name').value,
        gender: document.getElementById('patient-gender').value,
        age: parseInt(document.getElementById('patient-age').value),
        phone: document.getElementById('patient-phone').value,
        idCard: document.getElementById('patient-idcard').value,
        address: document.getElementById('patient-address').value,
        emergencyContact: document.getElementById('patient-emergency-contact').value,
        emergencyPhone: document.getElementById('patient-emergency-phone').value
    };

    // 验证必填字段
    if (!patient.name || !patient.gender || !patient.age || !patient.phone || !patient.idCard) {
        alert('请填写所有必填字段！');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/patients/createPatient`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(patient)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '保存失败');
        }

        closeModal();
        await loadPatients();
    } catch (error) {
        console.error('保存病人失败:', error);
        alert('保存失败，请重试！');
    }
}

async function updatePatient(id) {
    const patient = {
        name: document.getElementById('patient-name').value,
        gender: document.getElementById('patient-gender').value,
        age: parseInt(document.getElementById('patient-age').value),
        phone: document.getElementById('patient-phone').value,
        idCard: document.getElementById('patient-idcard').value,
        address: document.getElementById('patient-address').value,
        emergencyContact: document.getElementById('patient-emergency-contact').value,
        emergencyPhone: document.getElementById('patient-emergency-phone').value
    };

    if (!patient.name || !patient.gender || !patient.age || !patient.phone || !patient.idCard) {
        alert('请填写所有必填字段！');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/patients/updatePatient?id=${encodeURIComponent(id)}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(patient)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '保存失败');
        }

        closeModal();
        editingPatientId = null;
        await loadPatients();
    } catch (error) {
        console.error('更新病人失败:', error);
        alert('保存失败，请重试！');
    }
}

// 保存病种
async function saveDisease() {
    const disease = {
        name: document.getElementById('disease-name').value,
        category: document.getElementById('disease-category').value,
        description: document.getElementById('disease-description').value,
        symptoms: document.getElementById('disease-symptoms').value,
        treatment: document.getElementById('disease-treatment').value
    };

    // 验证必填字段
    if (!disease.name || !disease.category) {
        alert('请填写所有必填字段！');
        return;
    }

    try {
        const isEdit = !!editingDiseaseId;
        const url = isEdit
            ? `${API_BASE_URL}/diseases/updateDisease?id=${encodeURIComponent(editingDiseaseId)}`
            : `${API_BASE_URL}/diseases/createDisease`;

        const response = await fetch(url, {
            method: isEdit ? 'PUT' : 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(disease)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '保存失败');
        }

        closeModal();
        editingDiseaseId = null;
        await loadDiseases();
    } catch (error) {
        console.error('保存病种失败:', error);
        alert('保存失败，请重试！');
    }
}

// 加载病人数据
async function loadPatients() {
    try {
        const response = await fetch(`${API_BASE_URL}/patients/getPatients`, {
            method: 'GET'
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '加载失败');
        }

        const patients = await response.json();
        currentPatients = Array.isArray(patients) ? patients : [];

        applyPatientsFilterAndRender();

    } catch (error) {
        console.error('加载病人数据失败:', error);
        document.getElementById('patient-table-body').innerHTML = `
            <tr>
                <td colspan="8" style="text-align: center; padding: 50px;">
                    <i class="fas fa-exclamation-triangle" style="font-size: 48px; color: #ff9800; margin-bottom: 20px;"></i>
                    <p>加载数据失败，请刷新页面重试</p>
                </td>
            </tr>
        `;
    }
}

// 渲染病人表格
function renderPatientsTable(patients) {
    const tbody = document.getElementById('patient-table-body');

    if (patients.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="8" style="text-align: center; padding: 50px;">
                    <i class="fas fa-user-slash" style="font-size: 48px; color: #ccc; margin-bottom: 20px;"></i>
                    <p>暂无病人数据</p>
                </td>
            </tr>
        `;
        return;
    }

    let html = '';

    patients.forEach(patient => {
        html += `
            <tr>
                <td>${patient.id}</td>
                <td>${patient.name}</td>
                <td>${patient.gender}</td>
                <td>${patient.age}</td>
                <td>${patient.phone}</td>
                <td>${patient.idCard}</td>
                <td>${patient.address}</td>
                <td>
                    <button class="btn-action btn-edit" onclick="editPatient('${patient.id}')">
                        <i class="fas fa-edit"></i> 编辑
                    </button>
                    <button class="btn-action btn-delete" onclick="deletePatient('${patient.id}')">
                        <i class="fas fa-trash"></i> 删除
                    </button>
                </td>
            </tr>
        `;
    });

    tbody.innerHTML = html;
}

// 过滤病人
function filterPatients(query) {
    applyPatientsFilterAndRender(query);
}

// 编辑病人
function editPatient(id) {
    const patient = currentPatients.find(p => p.id === id);
    if (!patient) {
        alert('未找到该病人信息');
        return;
    }

    editingPatientId = id;

    const modalHtml = `
        <div class="modal active" id="edit-patient-modal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>编辑病人</h2>
                    <button class="btn-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="patient-form">
                        <div class="form-row">
                            <div class="form-group">
                                <label for="patient-name">姓名 *</label>
                                <input type="text" id="patient-name" required value="${patient.name ?? ''}">
                            </div>
                            <div class="form-group">
                                <label for="patient-gender">性别 *</label>
                                <select id="patient-gender" required>
                                    <option value="">请选择</option>
                                    <option value="男" ${patient.gender === '男' ? 'selected' : ''}>男</option>
                                    <option value="女" ${patient.gender === '女' ? 'selected' : ''}>女</option>
                                </select>
                            </div>
                        </div>
                        
                        <div class="form-row">
                            <div class="form-group">
                                <label for="patient-age">年龄 *</label>
                                <input type="number" id="patient-age" min="0" max="120" required value="${patient.age ?? ''}">
                            </div>
                            <div class="form-group">
                                <label for="patient-phone">电话 *</label>
                                <input type="tel" id="patient-phone" required value="${patient.phone ?? ''}">
                            </div>
                        </div>
                        
                        <div class="form-group">
                            <label for="patient-idcard">身份证号 *</label>
                            <input type="text" id="patient-idcard" required value="${patient.idCard ?? ''}">
                        </div>
                        
                        <div class="form-group">
                            <label for="patient-address">地址</label>
                            <input type="text" id="patient-address" value="${patient.address ?? ''}">
                        </div>
                        
                        <div class="form-row">
                            <div class="form-group">
                                <label for="patient-emergency-contact">紧急联系人</label>
                                <input type="text" id="patient-emergency-contact" value="${patient.emergencyContact ?? ''}">
                            </div>
                            <div class="form-group">
                                <label for="patient-emergency-phone">紧急联系电话</label>
                                <input type="tel" id="patient-emergency-phone" value="${patient.emergencyPhone ?? ''}">
                            </div>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary" id="cancel-patient-btn">取消</button>
                    <button class="btn-primary" id="save-patient-btn">保存</button>
                </div>
            </div>
        </div>
    `;

    document.getElementById('modal-container').innerHTML = modalHtml;

    document.querySelector('#edit-patient-modal .btn-close').addEventListener('click', closeModal);
    document.getElementById('cancel-patient-btn').addEventListener('click', closeModal);
    document.getElementById('save-patient-btn').addEventListener('click', () => updatePatient(id));
}

// 删除病人
async function deletePatient(id) {
    if (confirm('确定要删除这个病人吗？此操作不可恢复。')) {
        try {
            const response = await fetch(`${API_BASE_URL}/patients/deletePatient?id=${encodeURIComponent(id)}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || '删除失败');
            }

            await loadPatients();

        } catch (error) {
            console.error('删除病人失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

function applyPatientsFilterAndRender(searchQuery) {
    const query = (searchQuery ?? document.getElementById('patient-search')?.value ?? '').trim().toLowerCase();
    const filterValue = document.getElementById('patient-filter')?.value ?? 'all';

    const filtered = currentPatients.filter((p) => {
        if (filterValue === 'male' && p.gender !== '男') return false;
        if (filterValue === 'female' && p.gender !== '女') return false;

        if (!query) return true;
        const haystack = `${p.name ?? ''} ${p.phone ?? ''} ${p.idCard ?? ''} ${p.address ?? ''}`.toLowerCase();
        return haystack.includes(query);
    });

    renderPatientsTable(filtered);
}

// 加载病种数据
async function loadDiseases() {
    try {
        const response = await fetch(`${API_BASE_URL}/diseases/getDiseases`, {
            method: 'GET'
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '加载失败');
        }

        const diseases = await response.json();
        currentDiseases = Array.isArray(diseases) ? diseases : [];

        renderDiseasesCards(currentDiseases);

    } catch (error) {
        console.error('加载病种数据失败:', error);
        document.getElementById('disease-cards').innerHTML = `
            <div style="text-align: center; padding: 50px; width: 100%;">
                <i class="fas fa-exclamation-triangle" style="font-size: 48px; color: #ff9800; margin-bottom: 20px;"></i>
                <p>加载数据失败，请刷新页面重试</p>
            </div>
        `;
    }
}

// 渲染病种卡片
function renderDiseasesCards(diseases) {
    const container = document.getElementById('disease-cards');

    if (diseases.length === 0) {
        container.innerHTML = `
            <div style="text-align: center; padding: 50px; width: 100%;">
                <i class="fas fa-disease" style="font-size: 48px; color: #ccc; margin-bottom: 20px;"></i>
                <p>暂无病种数据</p>
            </div>
        `;
        return;
    }

    let html = '';

    diseases.forEach(disease => {
        html += `
            <div class="disease-card">
                <div class="disease-header">
                    <h3>${disease.name}</h3>
                    <span class="disease-category">${disease.category}</span>
                </div>
                <div class="disease-body">
                    <p class="disease-description">${disease.description}</p>
                    
                    <div class="disease-symptoms">
                        <h4>症状</h4>
                        <p>${disease.symptoms}</p>
                    </div>
                    
                    <div class="disease-symptoms">
                        <h4>治疗方法</h4>
                        <p>${disease.treatment}</p>
                    </div>
                </div>
                <div class="doctor-actions">
                    <button class="btn-action btn-edit" onclick="editDisease('${disease.id}')">
                        <i class="fas fa-edit"></i> 编辑
                    </button>
                    <button class="btn-action btn-delete" onclick="deleteDisease('${disease.id}')">
                        <i class="fas fa-trash"></i> 删除
                    </button>
                </div>
            </div>
        `;
    });

    container.innerHTML = html;
}

// 编辑病种
function editDisease(id) {
    const disease = currentDiseases.find(d => d.id === id);
    if (!disease) {
        alert('未找到该病种信息');
        return;
    }

    editingDiseaseId = id;

    const categories = ['内科', '外科', '儿科', '妇产科', '眼科', '耳鼻喉科', '口腔科', '皮肤科', '中医科'];
    const categoryOptions = categories.map((c) => `<option value="${c}" ${disease.category === c ? 'selected' : ''}>${c}</option>`).join('');

    const modalHtml = `
        <div class="modal active" id="edit-disease-modal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>编辑病种</h2>
                    <button class="btn-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="disease-form">
                        <div class="form-group">
                            <label for="disease-name">病种名称 *</label>
                            <input type="text" id="disease-name" required value="${disease.name ?? ''}">
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-category">分类 *</label>
                            <select id="disease-category" required>
                                <option value="">请选择</option>
                                ${categoryOptions}
                            </select>
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-description">描述</label>
                            <textarea id="disease-description" rows="3">${disease.description ?? ''}</textarea>
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-symptoms">症状</label>
                            <textarea id="disease-symptoms" rows="3">${disease.symptoms ?? ''}</textarea>
                        </div>
                        
                        <div class="form-group">
                            <label for="disease-treatment">治疗方法</label>
                            <textarea id="disease-treatment" rows="3">${disease.treatment ?? ''}</textarea>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary" id="cancel-disease-btn">取消</button>
                    <button class="btn-primary" id="save-disease-btn">保存</button>
                </div>
            </div>
        </div>
    `;

    document.getElementById('modal-container').innerHTML = modalHtml;
    document.querySelector('#edit-disease-modal .btn-close').addEventListener('click', closeModal);
    document.getElementById('cancel-disease-btn').addEventListener('click', closeModal);
    document.getElementById('save-disease-btn').addEventListener('click', saveDisease);
}

// 删除病种
async function deleteDisease(id) {
    if (confirm('确定要删除这个病种吗？此操作不可恢复。')) {
        try {
            const response = await fetch(`${API_BASE_URL}/diseases/deleteDisease?id=${encodeURIComponent(id)}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || '删除失败');
            }

            await loadDiseases();
        } catch (error) {
            console.error('删除病种失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

// 加载医生数据
async function loadDoctors() {
    try {
        if (!currentDiseases.length) {
            await loadDiseases();
        }

        const response = await fetch(`${API_BASE_URL}/doctors/getDoctors`, {
            method: 'GET'
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '加载失败');
        }

        const doctors = await response.json();
        currentDoctors = Array.isArray(doctors) ? doctors : [];

        renderDoctorsCards(currentDoctors);

    } catch (error) {
        console.error('加载医生数据失败:', error);
        document.getElementById('doctor-cards').innerHTML = `
            <div style="text-align: center; padding: 50px; width: 100%;">
                <i class="fas fa-exclamation-triangle" style="font-size: 48px; color: #ff9800; margin-bottom: 20px;"></i>
                <p>加载数据失败，请刷新页面重试</p>
            </div>
        `;
    }
}

// 渲染医生卡片
function renderDoctorsCards(doctors) {
    const container = document.getElementById('doctor-cards');
    const diseaseNameById = Object.fromEntries(currentDiseases.map(d => [d.id, d.name]));

    if (doctors.length === 0) {
        container.innerHTML = `
            <div style="text-align: center; padding: 50px; width: 100%;">
                <i class="fas fa-user-md" style="font-size: 48px; color: #ccc; margin-bottom: 20px;"></i>
                <p>暂无医生数据</p>
            </div>
        `;
        return;
    }

    let html = '';

    doctors.forEach(doctor => {
        // 获取工作日信息
        const workDays = (doctor.workSchedule ?? [])
            .filter(schedule => schedule.isAvailable)
            .map(schedule => schedule.dayOfWeek)
            .join('、');

        const diseaseTagsHtml = (doctor.diseases ?? [])
            .map((diseaseId) => `<span class="disease-tag">${diseaseNameById[diseaseId] || diseaseId}</span>`)
            .join('');

        html += `
            <div class="doctor-card">
                <div class="doctor-header">
                    <div class="doctor-avatar">
                        ${doctor.photo ? 
                            `<img src="${doctor.photo}" alt="${doctor.name}">` : 
                            `<i class="fas fa-user-md"></i>`
                        }
                    </div>
                    <div class="doctor-info">
                        <h3>${doctor.name}</h3>
                        <p>${doctor.department} · ${doctor.title}</p>
                    </div>
                </div>
                <div class="doctor-body">
                    <div class="doctor-details">
                        <p><i class="fas fa-stethoscope"></i> ${doctor.introduction}</p>
                        <p><i class="fas fa-calendar-alt"></i> 出诊时间: ${workDays}</p>
                        <p><i class="fas fa-users"></i> 每日限额: ${doctor.maxPatients}人</p>
                        <p><i class="fas fa-money-bill-wave"></i> 挂号费: ¥${doctor.fee.toFixed(2)}</p>
                    </div>
                    
                    <div class="disease-tags">
                        ${diseaseTagsHtml}
                    </div>
                </div>
                <div class="doctor-actions">
                    <button class="btn-action btn-edit" onclick="editDoctor('${doctor.id}')">
                        <i class="fas fa-edit"></i> 编辑
                    </button>
                    <button class="btn-action btn-delete" onclick="deleteDoctor('${doctor.id}')">
                        <i class="fas fa-trash"></i> 删除
                    </button>
                </div>
            </div>
        `;
    });

    container.innerHTML = html;
}

// 编辑医生
function editDoctor(id) {
    const doctor = currentDoctors.find(d => d.id === id);
    if (!doctor) {
        alert('未找到该医生信息');
        return;
    }

    editingDoctorId = id;
    openDoctorModal(doctor);
}

// 删除医生
async function deleteDoctor(id) {
    if (confirm('确定要删除这个医生吗？此操作不可恢复。')) {
        try {
            const response = await fetch(`${API_BASE_URL}/doctors/deleteDoctor?id=${encodeURIComponent(id)}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || '删除失败');
            }

            await loadDoctors();
        } catch (error) {
            console.error('删除医生失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

// 显示添加医生模态框
function showAddDoctorModal() {
    editingDoctorId = null;
    openDoctorModal();
}

function openDoctorModal(doctor) {
    const diseasesHtml = currentDiseases.map((d) => {
        const checked = doctor?.diseases?.includes(d.id) ? 'checked' : '';
        return `
            <label style="display: inline-flex; align-items: center; gap: 6px; margin-right: 12px; margin-bottom: 8px;">
                <input type="checkbox" name="doctor-disease" value="${d.id}" ${checked}>
                <span>${d.name}</span>
            </label>
        `;
    }).join('');

    const modalHtml = `
        <div class="modal active" id="doctor-modal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>${doctor ? '编辑医生' : '新增医生'}</h2>
                    <button class="btn-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="doctor-form">
                        <div class="form-row">
                            <div class="form-group">
                                <label for="doctor-name">姓名 *</label>
                                <input type="text" id="doctor-name" required value="${doctor?.name ?? ''}">
                            </div>
                            <div class="form-group">
                                <label for="doctor-department">科室 *</label>
                                <select id="doctor-department" required>
                                    <option value="">请选择</option>
                                    <option value="内科" ${doctor?.department === '内科' ? 'selected' : ''}>内科</option>
                                    <option value="外科" ${doctor?.department === '外科' ? 'selected' : ''}>外科</option>
                                    <option value="儿科" ${doctor?.department === '儿科' ? 'selected' : ''}>儿科</option>
                                    <option value="妇产科" ${doctor?.department === '妇产科' ? 'selected' : ''}>妇产科</option>
                                    <option value="眼科" ${doctor?.department === '眼科' ? 'selected' : ''}>眼科</option>
                                    <option value="耳鼻喉科" ${doctor?.department === '耳鼻喉科' ? 'selected' : ''}>耳鼻喉科</option>
                                    <option value="口腔科" ${doctor?.department === '口腔科' ? 'selected' : ''}>口腔科</option>
                                    <option value="皮肤科" ${doctor?.department === '皮肤科' ? 'selected' : ''}>皮肤科</option>
                                    <option value="中医科" ${doctor?.department === '中医科' ? 'selected' : ''}>中医科</option>
                                </select>
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label for="doctor-title">职称 *</label>
                                <input type="text" id="doctor-title" required value="${doctor?.title ?? ''}">
                            </div>
                            <div class="form-group">
                                <label for="doctor-maxPatients">每日限额</label>
                                <input type="number" id="doctor-maxPatients" min="1" max="200" value="${doctor?.maxPatients ?? 30}">
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="form-group" style="flex: 1;">
                                <label for="doctor-fee">挂号费</label>
                                <input type="number" id="doctor-fee" min="0" step="0.01" value="${doctor?.fee ?? 0}">
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="doctor-introduction">简介</label>
                            <textarea id="doctor-introduction" rows="3">${doctor?.introduction ?? ''}</textarea>
                        </div>

                        <div class="form-group">
                            <label>管理病种（1-3个） *</label>
                            <div id="doctor-diseases" style="display: flex; flex-wrap: wrap;">
                                ${diseasesHtml}
                            </div>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary" id="cancel-doctor-btn">取消</button>
                    <button class="btn-primary" id="save-doctor-btn">保存</button>
                </div>
            </div>
        </div>
    `;

    document.getElementById('modal-container').innerHTML = modalHtml;
    document.querySelector('#doctor-modal .btn-close').addEventListener('click', closeModal);
    document.getElementById('cancel-doctor-btn').addEventListener('click', closeModal);
    document.getElementById('save-doctor-btn').addEventListener('click', saveDoctor);
}

async function saveDoctor() {
    const selectedDiseaseIds = Array.from(document.querySelectorAll('input[name="doctor-disease"]:checked')).map((el) => el.value);
    if (selectedDiseaseIds.length < 1 || selectedDiseaseIds.length > 3) {
        alert('请选择1-3个病种');
        return;
    }

    const doctor = {
        name: document.getElementById('doctor-name').value,
        department: document.getElementById('doctor-department').value,
        title: document.getElementById('doctor-title').value,
        introduction: document.getElementById('doctor-introduction').value,
        photo: '',
        diseases: selectedDiseaseIds,
        maxPatients: parseInt(document.getElementById('doctor-maxPatients').value || '30', 10),
        fee: parseFloat(document.getElementById('doctor-fee').value || '0'),
        workSchedule: []
    };

    if (!doctor.name || !doctor.department || !doctor.title) {
        alert('请填写所有必填字段！');
        return;
    }

    try {
        const isEdit = !!editingDoctorId;
        const url = isEdit
            ? `${API_BASE_URL}/doctors/updateDoctor?id=${encodeURIComponent(editingDoctorId)}`
            : `${API_BASE_URL}/doctors/createDoctor`;

        const response = await fetch(url, {
            method: isEdit ? 'PUT' : 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(doctor)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '保存失败');
        }

        closeModal();
        editingDoctorId = null;
        await loadDoctors();
    } catch (error) {
        console.error('保存医生失败:', error);
        alert('保存失败，请重试！');
    }
}

// 加载挂号数据
async function loadRegistrations() {
    try {
        if (!currentPatients.length) {
            await loadPatients();
        }
        if (!currentDoctors.length) {
            await loadDoctors();
        }

        const response = await fetch(`${API_BASE_URL}/registrations/getRegistrations`, {
            method: 'GET'
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '加载失败');
        }

        const registrations = await response.json();
        currentRegistrations = Array.isArray(registrations) ? registrations : [];
        applyRegistrationFilterAndRender();

    } catch (error) {
        console.error('加载挂号数据失败:', error);
        document.getElementById('registration-table-body').innerHTML = `
            <tr>
                <td colspan="7" style="text-align: center; padding: 50px;">
                    <i class="fas fa-exclamation-triangle" style="font-size: 48px; color: #ff9800; margin-bottom: 20px;"></i>
                    <p>加载数据失败，请刷新页面重试</p>
                </td>
            </tr>
        `;
    }
}

function applyRegistrationFilterAndRender() {
    const status = document.getElementById('registration-status-filter')?.value ?? 'all';
    const dateStr = document.getElementById('registration-date-filter')?.value ?? '';
    const targetDate = dateStr ? new Date(`${dateStr}T00:00:00`).toISOString().slice(0, 10) : '';

    const filtered = currentRegistrations.filter((r) => {
        if (status !== 'all' && r.status !== status) return false;
        if (!targetDate) return true;
        const visitDate = new Date(r.visitDate);
        if (Number.isNaN(visitDate.getTime())) return false;
        const visitIso = visitDate.toISOString().slice(0, 10);
        return visitIso === targetDate;
    });

    renderRegistrationsTable(filtered);
}

// 渲染挂号表格
function renderRegistrationsTable(registrations) {
    const tbody = document.getElementById('registration-table-body');

    if (registrations.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="7" style="text-align: center; padding: 50px;">
                    <i class="fas fa-calendar-times" style="font-size: 48px; color: #ccc; margin-bottom: 20px;"></i>
                    <p>暂无挂号数据</p>
                </td>
            </tr>
        `;
        return;
    }

    const patientNames = Object.fromEntries(currentPatients.map(p => [p.id, p.name]));
    const doctorNames = Object.fromEntries(currentDoctors.map(d => [d.id, d.name]));

    let html = '';

    registrations.forEach(registration => {
        const statusText = {
            'pending': '待处理',
            'confirmed': '已确认',
            'completed': '已完成',
            'cancelled': '已取消'
        }[registration.status] || registration.status;

        const statusClass = `status-${registration.status}`;

        // 格式化日期
        const visitDate = new Date(registration.visitDate);
        const formattedDate = Number.isNaN(visitDate.getTime()) ? '' : visitDate.toLocaleDateString('zh-CN');
        const departments = Array.isArray(registration.departments) && registration.departments.length
            ? registration.departments
            : (registration.department ? [registration.department] : []);
        const departmentText = departments.join('、');

        html += `
            <tr>
                <td>${registration.id}</td>
                <td>${patientNames[registration.patientId] || '未知'}</td>
                <td>${doctorNames[registration.doctorId] || '未知'}</td>
                <td>${departmentText}</td>
                <td>${formattedDate} ${registration.timeSlot}</td>
                <td><span class="status-badge ${statusClass}">${statusText}</span></td>
                <td>
                    <button class="btn-action btn-edit" onclick="editRegistration('${registration.id}')">
                        <i class="fas fa-edit"></i> 编辑
                    </button>
                    <button class="btn-action btn-delete" onclick="deleteRegistration('${registration.id}')">
                        <i class="fas fa-trash"></i> 删除
                    </button>
                </td>
            </tr>
        `;
    });

    tbody.innerHTML = html;
}

// 编辑挂号
function editRegistration(id) {
    const registration = currentRegistrations.find(r => r.id === id);
    if (!registration) {
        alert('未找到该挂号信息');
        return;
    }

    editingRegistrationId = id;
    openRegistrationModal(registration);
}

// 删除挂号
async function deleteRegistration(id) {
    if (confirm('确定要删除这个挂号记录吗？此操作不可恢复。')) {
        try {
            const response = await fetch(`${API_BASE_URL}/registrations/deleteRegistration?id=${encodeURIComponent(id)}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || '删除失败');
            }

            await loadRegistrations();
        } catch (error) {
            console.error('删除挂号失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

// 显示添加挂号模态框
function showAddRegistrationModal() {
    editingRegistrationId = null;
    openRegistrationModal();
}

async function openRegistrationModal(registration) {
    if (!currentPatients.length) {
        await loadPatients();
    }
    if (!currentDoctors.length) {
        await loadDoctors();
    }

    const patientOptions = currentPatients
        .map((p) => `<option value="${p.id}" ${registration?.patientId === p.id ? 'selected' : ''}>${p.name}</option>`)
        .join('');

    const doctorOptions = currentDoctors
        .map((d) => `<option value="${d.id}" ${registration?.doctorId === d.id ? 'selected' : ''}>${d.name}（${d.department}）</option>`)
        .join('');

    const visitDateValue = registration?.visitDate ? new Date(registration.visitDate).toISOString().slice(0, 10) : '';
    const selectedDepartments = Array.isArray(registration?.departments) && registration.departments.length
        ? registration.departments
        : (registration?.department ? [registration.department] : []);
    const departmentsHtml = DEPARTMENTS.map((dept) => {
        const checked = selectedDepartments.includes(dept) ? 'checked' : '';
        return `
            <label style="display: inline-flex; align-items: center; gap: 6px; margin-right: 12px; margin-bottom: 8px;">
                <input type="checkbox" name="registration-department" value="${dept}" ${checked}>
                <span>${dept}</span>
            </label>
        `;
    }).join('');

    const modalHtml = `
        <div class="modal active" id="registration-modal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>${registration ? '编辑挂号' : '新增挂号'}</h2>
                    <button class="btn-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="registration-form">
                        <div class="form-row">
                            <div class="form-group">
                                <label for="registration-patient">病人 *</label>
                                <select id="registration-patient" required>
                                    <option value="">请选择</option>
                                    ${patientOptions}
                                </select>
                            </div>
                            <div class="form-group">
                                <label for="registration-doctor">医生 *</label>
                                <select id="registration-doctor" required>
                                    <option value="">请选择</option>
                                    ${doctorOptions}
                                </select>
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label for="registration-status">状态</label>
                                <select id="registration-status">
                                    <option value="pending" ${registration?.status === 'pending' ? 'selected' : ''}>待处理</option>
                                    <option value="confirmed" ${registration?.status === 'confirmed' ? 'selected' : ''}>已确认</option>
                                    <option value="completed" ${registration?.status === 'completed' ? 'selected' : ''}>已完成</option>
                                    <option value="cancelled" ${registration?.status === 'cancelled' ? 'selected' : ''}>已取消</option>
                                </select>
                            </div>
                        </div>

                        <div class="form-group">
                            <label>科室（1-n个） *</label>
                            <div id="registration-departments" style="display: flex; flex-wrap: wrap;">
                                ${departmentsHtml}
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label for="registration-visitDate">就诊日期 *</label>
                                <input type="date" id="registration-visitDate" required value="${visitDateValue}">
                            </div>
                            <div class="form-group">
                                <label for="registration-timeSlot">时间段 *</label>
                                <input type="text" id="registration-timeSlot" required placeholder="例如 09:00-09:30" value="${registration?.timeSlot ?? ''}">
                            </div>
                        </div>

                        <div class="form-group">
                            <label for="registration-symptoms">症状</label>
                            <textarea id="registration-symptoms" rows="2">${registration?.symptoms ?? ''}</textarea>
                        </div>

                        <div class="form-group">
                            <label for="registration-notes">备注</label>
                            <textarea id="registration-notes" rows="2">${registration?.notes ?? ''}</textarea>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary" id="cancel-registration-btn">取消</button>
                    <button class="btn-primary" id="save-registration-btn">保存</button>
                </div>
            </div>
        </div>
    `;

    document.getElementById('modal-container').innerHTML = modalHtml;
    document.querySelector('#registration-modal .btn-close').addEventListener('click', closeModal);
    document.getElementById('cancel-registration-btn').addEventListener('click', closeModal);
    document.getElementById('save-registration-btn').addEventListener('click', saveRegistration);

    const doctorSelect = document.getElementById('registration-doctor');
    const upsertAndCheckDepartment = (dept) => {
        if (!dept) return;
        const existing = Array.from(document.querySelectorAll('input[name="registration-department"]'))
            .find((el) => el.value === dept);
        if (existing) {
            existing.checked = true;
            return;
        }
        const container = document.getElementById('registration-departments');
        if (!container) return;
        const label = document.createElement('label');
        label.style.display = 'inline-flex';
        label.style.alignItems = 'center';
        label.style.gap = '6px';
        label.style.marginRight = '12px';
        label.style.marginBottom = '8px';
        const input = document.createElement('input');
        input.type = 'checkbox';
        input.name = 'registration-department';
        input.value = dept;
        input.checked = true;
        const span = document.createElement('span');
        span.textContent = dept;
        label.appendChild(input);
        label.appendChild(span);
        container.appendChild(label);
    };

    const ensureDepartmentSelection = () => {
        const checked = document.querySelectorAll('input[name="registration-department"]:checked');
        if (checked.length) return;
        const doctor = currentDoctors.find((d) => d.id === doctorSelect?.value);
        if (doctor?.department) {
            upsertAndCheckDepartment(doctor.department);
        }
    };

    doctorSelect?.addEventListener('change', ensureDepartmentSelection);
    ensureDepartmentSelection();
}

async function saveRegistration() {
    const patientId = document.getElementById('registration-patient').value;
    const doctorId = document.getElementById('registration-doctor').value;
    const departments = Array.from(document.querySelectorAll('input[name="registration-department"]:checked'))
        .map((el) => el.value);
    const status = document.getElementById('registration-status').value;
    const visitDateInput = document.getElementById('registration-visitDate').value;
    const timeSlot = document.getElementById('registration-timeSlot').value;
    const symptoms = document.getElementById('registration-symptoms').value;
    const notes = document.getElementById('registration-notes').value;

    if (!patientId || !doctorId || !visitDateInput || !timeSlot) {
        alert('请填写所有必填字段！');
        return;
    }
    if (departments.length < 1) {
        alert('请选择至少1个科室！');
        return;
    }

    const visitDate = new Date(`${visitDateInput}T00:00:00`);
    if (Number.isNaN(visitDate.getTime())) {
        alert('就诊日期格式不正确');
        return;
    }

    const payload = {
        patientId,
        doctorId,
        department: departments[0],
        departments,
        visitDate: visitDate.toISOString(),
        timeSlot,
        status,
        symptoms,
        notes
    };

    try {
        const isEdit = !!editingRegistrationId;
        const url = isEdit
            ? `${API_BASE_URL}/registrations/updateRegistration?id=${encodeURIComponent(editingRegistrationId)}`
            : `${API_BASE_URL}/registrations/createRegistration`;

        const response = await fetch(url, {
            method: isEdit ? 'PUT' : 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || '保存失败');
        }

        closeModal();
        editingRegistrationId = null;
        await loadRegistrations();
    } catch (error) {
        console.error('保存挂号失败:', error);
        alert('保存失败，请重试！');
    }
}

// 加载报表数据
function loadReports() {
    // 初始化图表
    initRegistrationChart();
    initDepartmentChart();
}

// 初始化挂号统计图表
function initRegistrationChart() {
    const ctx = document.getElementById('registrations-chart').getContext('2d');

    // 模拟数据
    const data = {
        labels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
        datasets: [{
            label: '挂号数量',
            data: [45, 52, 48, 60, 55, 30, 25],
            backgroundColor: 'rgba(76, 175, 80, 0.2)',
            borderColor: 'rgba(76, 175, 80, 1)',
            borderWidth: 2,
            tension: 0.4
        }]
    };

    new Chart(ctx, {
        type: 'line',
        data: data,
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: 'top',
                },
                title: {
                    display: true,
                    text: '本周挂号数量统计'
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: '挂号数量'
                    }
                }
            }
        }
    });
}

// 初始化科室分布图表
function initDepartmentChart() {
    const ctx = document.getElementById('departments-chart').getContext('2d');

    // 模拟数据
    const data = {
        labels: ['内科', '外科', '儿科', '妇产科', '眼科', '其他'],
        datasets: [{
            label: '挂号数量',
            data: [120, 85, 60, 45, 30, 40],
            backgroundColor: [
                'rgba(33, 150, 243, 0.6)',
                'rgba(76, 175, 80, 0.6)',
                'rgba(255, 152, 0, 0.6)',
                'rgba(156, 39, 176, 0.6)',
                'rgba(244, 67, 54, 0.6)',
                'rgba(158, 158, 158, 0.6)'
            ],
            borderColor: [
                'rgba(33, 150, 243, 1)',
                'rgba(76, 175, 80, 1)',
                'rgba(255, 152, 0, 1)',
                'rgba(156, 39, 176, 1)',
                'rgba(244, 67, 54, 1)',
                'rgba(158, 158, 158, 1)'
            ],
            borderWidth: 1
        }]
    };

    new Chart(ctx, {
        type: 'doughnut',
        data: data,
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: 'right',
                },
                title: {
                    display: true,
                    text: '科室挂号分布'
                }
            }
        }
    });
}
