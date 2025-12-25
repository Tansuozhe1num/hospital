// 医生管理模块
class DoctorManager {
    constructor() {
        this.doctors = [];
    }

    // 加载医生数据
    async loadDoctors() {
        try {
            const response = await fetch(`${API_BASE_URL}/doctors`);
            this.doctors = await response.json();
            return this.doctors;
        } catch (error) {
            console.error('加载医生数据失败:', error);
            return [];
        }
    }

    // 按科室筛选
    filterByDepartment(department) {
        if (!department) return this.doctors;
        return this.doctors.filter(doctor => doctor.department === department);
    }

    // 获取所有科室
    getDepartments() {
        const departments = new Set();
        this.doctors.forEach(doctor => departments.add(doctor.department));
        return Array.from(departments);
    }

    // 获取医生的病种
    async getDoctorDiseases(doctorId) {
        try {
            const response = await fetch(`${API_BASE_URL}/doctors/${doctorId}/diseases`);
            return await response.json();
        } catch (error) {
            console.error('获取医生病种失败:', error);
            return [];
        }
    }
}

// 创建全局实例
window.doctorManager = new DoctorManager();