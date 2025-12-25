// 挂号管理模块
class RegistrationManager {
    constructor() {
        this.registrations = [];
        this.filteredRegistrations = [];
    }

    // 加载挂号数据
    async loadRegistrations() {
        try {
            const response = await fetch(`${API_BASE_URL}/registrations`);
            this.registrations = await response.json();
            this.filteredRegistrations = [...this.registrations];
            return this.registrations;
        } catch (error) {
            console.error('加载挂号数据失败:', error);
            return [];
        }
    }

    // 按状态筛选
    filterByStatus(status) {
        if (status === 'all') {
            this.filteredRegistrations = [...this.registrations];
        } else {
            this.filteredRegistrations = this.registrations.filter(
                reg => reg.status === status
            );
        }
        return this.filteredRegistrations;
    }

    // 按日期筛选
    filterByDate(date) {
        if (!date) {
            this.filteredRegistrations = [...this.registrations];
            return this.filteredRegistrations;
        }

        const targetDate = new Date(date).toISOString().split('T')[0];
        this.filteredRegistrations = this.registrations.filter(reg => {
            const regDate = new Date(reg.visitDate).toISOString().split('T')[0];
            return regDate === targetDate;
        });

        return this.filteredRegistrations;
    }

    // 获取病人挂号记录
    async getPatientRegistrations(patientId) {
        try {
            const response = await fetch(`${API_BASE_URL}/registrations/patient/${patientId}`);
            return await response.json();
        } catch (error) {
            console.error('获取病人挂号记录失败:', error);
            return [];
        }
    }
}

// 创建全局实例
window.registrationManager = new RegistrationManager();