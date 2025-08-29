class AITshirtShop {
    constructor() {
        this.baseURL = '/api/v1';
        this.token = localStorage.getItem('authToken');
        this.currentUser = JSON.parse(localStorage.getItem('currentUser') || 'null');
        
        this.init();
    }

    init() {
        this.bindEvents();
        this.checkAuthState();
        this.showSection('homeSection');
    }

    bindEvents() {
        // 导航菜单点击事件
        document.getElementById('navHome').addEventListener('click', () => this.showSection('homeSection'));
        document.getElementById('navDesign').addEventListener('click', () => this.showSection('designSection'));
        document.getElementById('navGallery').addEventListener('click', () => this.showSection('gallerySection'));
        
        // 认证按钮
        document.getElementById('btnLogin').addEventListener('click', () => this.showAuthModal('login'));
        document.getElementById('btnRegister').addEventListener('click', () => this.showAuthModal('register'));
        document.getElementById('btnLogout').addEventListener('click', () => this.logout());
        
        // 开始设计按钮
        document.getElementById('startDesignBtn').addEventListener('click', () => this.showSection('designSection'));
        
        // 生成设计按钮
        document.getElementById('generateBtn').addEventListener('click', () => this.generateDesign());
        
        // 下载和保存按钮
        document.getElementById('downloadBtn').addEventListener('click', () => this.downloadDesign());
        document.getElementById('saveDesignBtn').addEventListener('click', () => this.saveDesign());
        
        // 模态框事件
        document.getElementById('modalClose').addEventListener('click', () => this.hideAuthModal());
        document.getElementById('authForm').addEventListener('submit', (e) => this.handleAuthSubmit(e));
        document.getElementById('switchToRegister').addEventListener('click', (e) => {
            e.preventDefault();
            this.showAuthModal('register');
        });
        
        // 点击模态框背景关闭
        document.getElementById('authModal').addEventListener('click', (e) => {
            if (e.target.id === 'authModal') {
                this.hideAuthModal();
            }
        });
    }

    showSection(sectionId) {
        // 隐藏所有章节
        document.querySelectorAll('.section').forEach(section => {
            section.classList.remove('active');
        });
        
        // 显示目标章节
        document.getElementById(sectionId).classList.add('active');
        
        // 更新导航激活状态
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.remove('active');
        });
        
        const navItem = document.getElementById(`nav${sectionId.replace('Section', '')}`);
        if (navItem) {
            navItem.classList.add('active');
        }
    }

    checkAuthState() {
        const authSection = document.getElementById('navAuth');
        const userSection = document.getElementById('navUser');
        const welcomeText = document.getElementById('userWelcome');
        
        if (this.token && this.currentUser) {
            authSection.style.display = 'none';
            userSection.style.display = 'flex';
            welcomeText.textContent = `欢迎, ${this.currentUser.username}`;
        } else {
            authSection.style.display = 'flex';
            userSection.style.display = 'none';
        }
    }

    showAuthModal(mode = 'login') {
        const modal = document.getElementById('authModal');
        const title = document.getElementById('modalTitle');
        const emailGroup = document.getElementById('emailGroup');
        const switchText = document.getElementById('authSwitch');
        
        if (mode === 'login') {
            title.textContent = '登录';
            emailGroup.style.display = 'none';
            switchText.innerHTML = '没有账号？ <a href="#" id="switchToRegister">立即注册</a>';
        } else {
            title.textContent = '注册';
            emailGroup.style.display = 'block';
            switchText.innerHTML = '已有账号？ <a href="#" id="switchToLogin">立即登录</a>';
            
            // 重新绑定切换链接
            document.getElementById('switchToLogin')?.addEventListener('click', (e) => {
                e.preventDefault();
                this.showAuthModal('login');
            });
        }
        
        modal.style.display = 'flex';
        document.getElementById('username').focus();
    }

    hideAuthModal() {
        document.getElementById('authModal').style.display = 'none';
        document.getElementById('authForm').reset();
    }

    async handleAuthSubmit(e) {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const email = document.getElementById('email').value;
        const isRegister = document.getElementById('emailGroup').style.display === 'block';
        
        const endpoint = isRegister ? '/auth/register' : '/auth/login';
        const payload = isRegister ? { username, password, email } : { username, password };
        
        try {
            const response = await this.apiRequest(endpoint, 'POST', payload, false);
            
            if (response.token) {
                this.token = response.token;
                this.currentUser = { username: response.user?.username || username };
                
                localStorage.setItem('authToken', this.token);
                localStorage.setItem('currentUser', JSON.stringify(this.currentUser));
                
                this.hideAuthModal();
                this.checkAuthState();
                
                if (isRegister) {
                    this.showNotification('注册成功！已自动登录', 'success');
                    this.showSection('designSection');
                } else {
                    this.showNotification('登录成功！', 'success');
                }
            }
        } catch (error) {
            this.showNotification(error.message, 'error');
        }
    }

    logout() {
        localStorage.removeItem('authToken');
        localStorage.removeItem('currentUser');
        this.token = null;
        this.currentUser = null;
        this.checkAuthState();
        this.showNotification('已退出登录', 'success');
    }

    async generateDesign() {
        if (!this.token) {
            this.showNotification('请先登录以生成设计', 'error');
            this.showAuthModal('login');
            return;
        }
        
        const prompt = document.getElementById('promptInput').value.trim();
        const style = document.getElementById('styleSelect').value;
        
        if (!prompt) {
            this.showNotification('请输入设计描述', 'error');
            return;
        }
        
        // 构建完整的提示词
        let fullPrompt = prompt;
        if (style) {
            fullPrompt += `, ${style}风格`;
        }
        
        this.showLoading(true);
        
        try {
            const response = await this.apiRequest('/designs/generate', 'POST', {
                prompt: fullPrompt
            });
            
            // 显示生成的图片
            const img = document.getElementById('generatedImage');
            const placeholder = document.getElementById('designPlaceholder');
            
            img.src = response.image_url;
            img.style.display = 'block';
            placeholder.style.display = 'none';
            
            // 启用下载和保存按钮
            document.getElementById('downloadBtn').disabled = false;
            document.getElementById('saveDesignBtn').disabled = false;
            
            this.currentDesign = response;
            this.showNotification('设计生成成功！', 'success');
            
        } catch (error) {
            this.showNotification(error.message, 'error');
        } finally {
            this.showLoading(false);
        }
    }

    downloadDesign() {
        if (!this.currentDesign) return;
        
        const link = document.createElement('a');
        link.href = this.currentDesign.image_url;
        link.download = `tshirt-design-${Date.now()}.png`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        
        this.showNotification('设计已下载', 'success');
    }

    async saveDesign() {
        if (!this.currentDesign) return;
        
        try {
            // 这里可以添加保存设计到用户收藏的逻辑
            // 例如: await this.apiRequest('/designs/save', 'POST', { image_url: this.currentDesign.image_url });
            
            this.showNotification('设计已保存到收藏', 'success');
        } catch (error) {
            this.showNotification('保存失败: ' + error.message, 'error');
        }
    }

    async loadGallery() {
        if (!this.token) return;
        
        try {
            // 这里可以添加加载用户设计作品的逻辑
            // const designs = await this.apiRequest('/designs/my-designs', 'GET');
            // this.renderGallery(designs);
            
        } catch (error) {
            console.error('加载画廊失败:', error);
        }
    }

    renderGallery(designs) {
        const galleryGrid = document.getElementById('galleryGrid');
        
        if (designs.length === 0) {
            galleryGrid.innerHTML = `
                <div class="gallery-placeholder">
                    <i class="fas fa-paint-brush"></i>
                    <p>还没有设计作品，快去创作吧！</p>
                </div>
            `;
            return;
        }
        
        galleryGrid.innerHTML = designs.map(design => `
            <div class="gallery-item">
                <img src="${design.image_url}" alt="T恤设计">
                <div class="gallery-item-content">
                    <p>${design.created_at}</p>
                </div>
            </div>
        `).join('');
    }

    async apiRequest(endpoint, method = 'GET', data = null, auth = true) {
        const url = `${this.baseURL}${endpoint}`;
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json',
            },
        };
        
        if (auth && this.token) {
            options.headers['Authorization'] = `Bearer ${this.token}`;
        }
        
        if (data && method !== 'GET') {
            options.body = JSON.stringify(data);
        }
        
        const response = await fetch(url, options);
        const result = await response.json();
        
        if (!response.ok) {
            throw new Error(result.error || result.message || '请求失败');
        }
        
        return result;
    }

    showLoading(show) {
        document.getElementById('loading').style.display = show ? 'flex' : 'none';
    }

    showNotification(message, type = 'info') {
        const notification = document.getElementById('notification');
        notification.textContent = message;
        notification.className = `notification ${type} show`;
        
        setTimeout(() => {
            notification.classList.remove('show');
        }, 3000);
    }
}

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new AITshirtShop();
});

// 工具函数
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

function formatDate(date) {
    return new Date(date).toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}