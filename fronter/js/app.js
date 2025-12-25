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
const API_BASE_URL = 'http://localhost:8080/api';

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
        emergencyPhone: document.getElementById('patient-emergency-phone').value,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
    };

    // 验证必填字段
    if (!patient.name || !patient.gender || !patient.age || !patient.phone || !patient.idCard) {
        alert('请填写所有必填字段！');
        return;
    }

    try {
        // 实际API调用
        /*
        const response = await fetch(`${API_BASE_URL}/patients`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(patient)
        });

        if (response.ok) {
            closeModal();
            loadPatients();
        } else {
            throw new Error('保存失败');
        }
        */

        // 模拟成功
        console.log('保存病人:', patient);
        closeModal();
        loadPatients();

    } catch (error) {
        console.error('保存病人失败:', error);
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
        // 模拟成功
        console.log('保存病种:', disease);
        closeModal();
        loadDiseases();

    } catch (error) {
        console.error('保存病种失败:', error);
        alert('保存失败，请重试！');
    }
}

// 加载病人数据
async function loadPatients() {
    try {
        // 模拟数据
        const mockPatients = [
            {
                id: '1',
                name: '张三',
                gender: '男',
                age: 35,
                phone: '13800138000',
                idCard: '110101199001011234',
                address: '北京市朝阳区',
                emergencyContact: '李四',
                emergencyPhone: '13800138001',
                createdAt: '2023-01-15T10:30:00Z',
                updatedAt: '2023-01-15T10:30:00Z'
            },
            {
                id: '2',
                name: '王芳',
                gender: '女',
                age: 28,
                phone: '13900139000',
                idCard: '110101199501011235',
                address: '上海市浦东新区',
                emergencyContact: '王明',
                emergencyPhone: '13900139001',
                createdAt: '2023-02-20T14:20:00Z',
                updatedAt: '2023-02-20T14:20:00Z'
            },
            {
                id: '3',
                name: '李强',
                gender: '男',
                age: 45,
                phone: '13700137000',
                idCard: '110101197801011236',
                address: '广州市天河区',
                emergencyContact: '张红',
                emergencyPhone: '13700137001',
                createdAt: '2023-03-10T09:15:00Z',
                updatedAt: '2023-03-10T09:15:00Z'
            }
        ];

        // 实际API调用
        /*
        const response = await fetch(`${API_BASE_URL}/patients`);
        const patients = await response.json();
        */

        const patients = mockPatients;

        renderPatientsTable(patients);

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
    // 实现搜索过滤逻辑
    console.log('过滤病人:', query);
}

// 编辑病人
function editPatient(id) {
    console.log('编辑病人:', id);
    alert('编辑病人功能开发中...');
}

// 删除病人
async function deletePatient(id) {
    if (confirm('确定要删除这个病人吗？此操作不可恢复。')) {
        try {
            // 实际API调用
            /*
            const response = await fetch(`${API_BASE_URL}/patients/${id}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                loadPatients();
            } else {
                throw new Error('删除失败');
            }
            */

           // 模拟成功
           console.log('删除病人:', id);
           loadPatients();

        } catch (error) {
            console.error('删除病人失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

// 加载病种数据
async function loadDiseases() {
    try {
        // 模拟数据
        const mockDiseases = [
            {
                id: '1',
                name: '高血压',
                category: '内科',
                description: '血压持续升高的一种疾病',
                symptoms: '头痛、头晕、心悸',
                treatment: '药物治疗、饮食控制'
            },
            {
                id: '2',
                name: '糖尿病',
                category: '内科',
                description: '血糖调节异常导致的代谢疾病',
                symptoms: '多饮、多尿、体重下降',
                treatment: '胰岛素治疗、饮食管理'
            },
            {
                id: '3',
                name: '急性阑尾炎',
                category: '外科',
                description: '阑尾急性炎症',
                symptoms: '右下腹疼痛、发热',
                treatment: '手术切除'
            }
        ];

        // 实际API调用
        /*
        const response = await fetch(`${API_BASE_URL}/diseases`);
        const diseases = await response.json();
        */

        const diseases = mockDiseases;

        renderDiseasesCards(diseases);

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
    console.log('编辑病种:', id);
    alert('编辑病种功能开发中...');
}

// 删除病种
async function deleteDisease(id) {
    if (confirm('确定要删除这个病种吗？此操作不可恢复。')) {
        try {
            // 模拟成功
            console.log('删除病种:', id);
            loadDiseases();

        } catch (error) {
            console.error('删除病种失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

// 加载医生数据
async function loadDoctors() {
    try {
        // 模拟数据
        const mockDoctors = [
            {
                id: '1',
                name: '王医生',
                department: '内科',
                title: '主任医师',
                introduction: '擅长高血压、糖尿病等内科疾病的诊断与治疗',
                photo: '',
                diseases: ['1', '2'],
                maxPatients: 30,
                fee: 50.00,
                workSchedule: [
                    { dayOfWeek: '周一', startTime: '08:30', endTime: '17:00', isAvailable: true },
                    { dayOfWeek: '周三', startTime: '08:30', endTime: '17:00', isAvailable: true },
                    { dayOfWeek: '周五', startTime: '08:30', endTime: '17:00', isAvailable: true }
                ]
            },
            {
                id: '2',
                name: '李医生',
                department: '外科',
                title: '副主任医师',
                introduction: '擅长普外科手术，包括阑尾炎、胆囊炎等',
                photo: '',
                diseases: ['3'],
                maxPatients: 20,
                fee: 80.00,
                workSchedule: [
                    { dayOfWeek: '周二', startTime: '09:00', endTime: '17:00', isAvailable: true },
                    { dayOfWeek: '周四', startTime: '09:00', endTime: '17:00', isAvailable: true }
                ]
            }
        ];

        // 实际API调用
        /*
        const response = await fetch(`${API_BASE_URL}/doctors`);
        const doctors = await response.json();
        */

        const doctors = mockDoctors;

        renderDoctorsCards(doctors);

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
        const workDays = doctor.workSchedule
            .filter(schedule => schedule.isAvailable)
            .map(schedule => schedule.dayOfWeek)
            .join('、');

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
                        <span class="disease-tag">高血压</span>
                        <span class="disease-tag">糖尿病</span>
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
    console.log('编辑医生:', id);
    alert('编辑医生功能开发中...');
}

// 删除医生
async function deleteDoctor(id) {
    if (confirm('确定要删除这个医生吗？此操作不可恢复。')) {
        try {
            // 模拟成功
            console.log('删除医生:', id);
            loadDoctors();

        } catch (error) {
            console.error('删除医生失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

// 显示添加医生模态框
function showAddDoctorModal() {
    alert('添加医生功能开发中...');
}

// 加载挂号数据
async function loadRegistrations() {
    try {
        // 模拟数据
        const mockRegistrations = [
            {
                id: 'R001',
                patientId: '1',
                doctorId: '1',
                department: '内科',
                registrationDate: '2023-10-15T09:30:00Z',
                visitDate: '2023-10-16T14:30:00Z',
                timeSlot: '14:30-15:00',
                status: 'confirmed',
                symptoms: '头痛、头晕',
                notes: '高血压复诊',
                createdAt: '2023-10-15T09:30:00Z'
            },
            {
                id: 'R002',
                patientId: '2',
                doctorId: '2',
                department: '外科',
                registrationDate: '2023-10-16T10:15:00Z',
                visitDate: '2023-10-17T10:00:00Z',
                timeSlot: '10:00-10:30',
                status: 'pending',
                symptoms: '右下腹疼痛',
                notes: '疑似阑尾炎',
                createdAt: '2023-10-16T10:15:00Z'
            },
            {
                id: 'R003',
                patientId: '3',
                doctorId: '1',
                department: '内科',
                registrationDate: '2023-10-14T14:20:00Z',
                visitDate: '2023-10-15T09:00:00Z',
                timeSlot: '09:00-09:30',
                status: 'completed',
                symptoms: '多饮、多尿',
                notes: '糖尿病检查',
                createdAt: '2023-10-14T14:20:00Z'
            }
        ];

        // 实际API调用
        /*
        const response = await fetch(`${API_BASE_URL}/registrations`);
        const registrations = await response.json();
        */

        const registrations = mockRegistrations;

        renderRegistrationsTable(registrations);

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

    // 模拟病人和医生数据
    const patientNames = {
        '1': '张三',
        '2': '王芳',
        '3': '李强'
    };

    const doctorNames = {
        '1': '王医生',
        '2': '李医生'
    };

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
        const formattedDate = visitDate.toLocaleDateString('zh-CN');

        html += `
            <tr>
                <td>${registration.id}</td>
                <td>${patientNames[registration.patientId] || '未知'}</td>
                <td>${doctorNames[registration.doctorId] || '未知'}</td>
                <td>${registration.department}</td>
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
    console.log('编辑挂号:', id);
    alert('编辑挂号功能开发中...');
}

// 删除挂号
async function deleteRegistration(id) {
    if (confirm('确定要删除这个挂号记录吗？此操作不可恢复。')) {
        try {
            // 模拟成功
            console.log('删除挂号:', id);
            loadRegistrations();

        } catch (error) {
            console.error('删除挂号失败:', error);
            alert('删除失败，请重试！');
        }
    }
}

// 显示添加挂号模态框
function showAddRegistrationModal() {
    alert('添加挂号功能开发中...');
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