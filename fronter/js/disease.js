// 病种管理模块
class DiseaseManager {
    constructor() {
        this.diseases = [];
    }

    // 加载病种数据
    async loadDiseases() {
        try {
            const response = await fetch(`${API_BASE_URL}/diseases`);
            this.diseases = await response.json();
            return this.diseases;
        } catch (error) {
            console.error('加载病种数据失败:', error);
            return [];
        }
    }

    // 按分类筛选
    filterByCategory(category) {
        if (!category) return this.diseases;
        return this.diseases.filter(disease => disease.category === category);
    }

    // 获取所有分类
    getCategories() {
        const categories = new Set();
        this.diseases.forEach(disease => categories.add(disease.category));
        return Array.from(categories);
    }
}

// 创建全局实例
window.diseaseManager = new DiseaseManager();