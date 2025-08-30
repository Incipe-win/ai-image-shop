class AICreativeStudio {
    constructor() {
        this.baseURL = '/api/v1';
        this.token = localStorage.getItem('authToken');
        this.refreshToken = localStorage.getItem('refreshToken');
        this.currentUser = JSON.parse(localStorage.getItem('currentUser') || 'null');
        
        // 电商相关属性
        this.cartItems = JSON.parse(localStorage.getItem('cartItems') || '[]');
        this.currentProduct = null;
        this.selectedCartItems = [];
        
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
        document.getElementById('navShop').addEventListener('click', () => this.showSection('shopSection'));
        document.getElementById('navCart').addEventListener('click', () => this.showSection('cartSection'));
        document.getElementById('navOrders').addEventListener('click', () => this.showSection('ordersSection'));
        
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
        document.getElementById('publishToShopBtn').addEventListener('click', () => this.showPublishModal());

        // 电商功能事件绑定
        document.getElementById('applyFilters').addEventListener('click', () => this.loadProducts());
        document.getElementById('applyGalleryFilters').addEventListener('click', () => this.loadGallery());
        document.getElementById('addToCartForm').addEventListener('submit', (e) => this.handleAddToCart(e));
        document.getElementById('publishForm').addEventListener('submit', (e) => this.handlePublishToShop(e));
        document.getElementById('checkoutBtn').addEventListener('click', () => this.showCheckoutModal());
        document.getElementById('confirmOrderBtn').addEventListener('click', () => this.createOrder());
        
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
        
        // 产品详情模态框背景点击关闭
        document.addEventListener('click', (e) => {
            if (e.target.id === 'productDetailsModal') {
                this.hideProductDetailsModal();
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

        // 页面特定加载逻辑
        if (sectionId === 'gallerySection') {
            this.loadGallery();
        } else if (sectionId === 'shopSection') {
            this.loadProducts();
        } else if (sectionId === 'cartSection') {
            this.loadCart();
        } else if (sectionId === 'ordersSection') {
            this.loadOrders();
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
            
            // 重新绑定切换到注册的链接
            document.getElementById('switchToRegister')?.addEventListener('click', (e) => {
                e.preventDefault();
                this.showAuthModal('register');
            });
        } else {
            title.textContent = '注册';
            emailGroup.style.display = 'block';
            switchText.innerHTML = '已有账号？ <a href="#" id="switchToLogin">立即登录</a>';
            
            // 重新绑定切换到登录的链接
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
                this.refreshToken = response.refresh_token;
                this.currentUser = { username: response.user?.username || username };
                
                localStorage.setItem('authToken', this.token);
                localStorage.setItem('refreshToken', this.refreshToken);
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
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('currentUser');
        this.token = null;
        this.refreshToken = null;
        this.currentUser = null;
        this.checkAuthState();
        this.showNotification('已退出登录', 'success');
    }

    async generateDesign() {
        if (!this.token) {
            this.showNotification('请先登录以生成创意作品', 'error');
            this.showAuthModal('login');
            return;
        }
        
        const prompt = document.getElementById('promptInput').value.trim();
        const style = document.getElementById('styleSelect').value;
        const category = document.getElementById('categorySelect').value;
        
        if (!prompt) {
            this.showNotification('请输入创意描述', 'error');
            return;
        }
        
        // 构建完整的提示词
        let fullPrompt = prompt;
        if (category) {
            const categoryMap = {
                'poster': '海报印刷',
                'sticker': '贴纸定制',
                'canvas': '画布装饰',
                'tshirt': 'T恤图案'
            };
            fullPrompt += `，适用于${categoryMap[category] || category}`;
        }
        if (style) {
            fullPrompt += `, ${style}风格`;
        }
        
        this.showLoading(true);
        
        try {
            const response = await this.apiRequest('/designs/generate', 'POST', {
                prompt: fullPrompt,
                category: category || 'general',
                style: style || ''
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
            document.getElementById('publishToShopBtn').disabled = false;
            
            this.currentDesign = response;
            this.showNotification('创意作品生成成功！', 'success');
            
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
        link.download = `creative-artwork-${Date.now()}.png`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        
        this.showNotification('创意作品已下载', 'success');
    }

    async saveDesign() {
        if (!this.currentDesign) return;
        
        try {
            // 这里可以添加保存设计到用户收藏的逻辑
            // 例如: await this.apiRequest('/designs/save', 'POST', { image_url: this.currentDesign.image_url });
            
            this.showNotification('创意作品已保存到收藏', 'success');
        } catch (error) {
            this.showNotification('保存失败: ' + error.message, 'error');
        }
    }

    async loadGallery() {
        const galleryGrid = document.getElementById('galleryGrid');
        
        if (!this.token) {
            // 未登录状态显示提示信息
            galleryGrid.innerHTML = `
                <div class="gallery-placeholder">
                    <i class="fas fa-user-lock"></i>
                    <p>登录后查看你的创意作品</p>
                </div>
            `;
            return;
        }
        
        try {
            const category = document.getElementById('galleryFilter')?.value || '';
            let endpoint = '/designs/my-designs';
            if (category) {
                endpoint = `/designs/my-designs?category=${encodeURIComponent(category)}`;
            }
            
            const response = await this.apiRequest(endpoint, 'GET');
            this.renderGallery(response.designs);
        } catch (error) {
            console.error('加载画廊失败:', error);
            galleryGrid.innerHTML = `
                <div class="gallery-placeholder">
                    <i class="fas fa-exclamation-triangle"></i>
                    <p>加载失败: ${error.message}</p>
                </div>
            `;
        }
    }

    renderGallery(designs) {
        const galleryGrid = document.getElementById('galleryGrid');
        
        if (!designs || designs.length === 0) {
            galleryGrid.innerHTML = `
                <div class="gallery-placeholder">
                    <i class="fas fa-paint-brush"></i>
                    <p>还没有创意作品，快去创作吧！</p>
                </div>
            `;
            return;
        }
        
        galleryGrid.innerHTML = designs.map(design => `
            <div class="gallery-item">
                <img src="${design.image_url}" alt="创意作品">
                <div class="gallery-item-content">
                    <div class="gallery-item-meta">
                        <span class="category-tag">${this.getCategoryName(design.category)}</span>
                        ${design.style ? `<span class="style-tag">${design.style}</span>` : ''}
                    </div>
                    <p class="design-prompt">${design.prompt || '创意作品'}</p>
                    <p class="creation-time">${this.formatDate(design.created_at) || '刚刚创建'}</p>
                </div>
            </div>
        `).join('');
    }

    getCategoryName(category) {
        const categoryMap = {
            'poster': '海报印刷',
            'sticker': '贴纸定制', 
            'canvas': '画布装饰',
            'tshirt': 'T恤图案'
        };
        return categoryMap[category] || '其他分类';
    }

    formatDate(dateString) {
        if (!dateString) return '';
        try {
            return new Date(dateString).toLocaleDateString('zh-CN', {
                year: 'numeric',
                month: 'short',
                day: 'numeric'
            });
        } catch (e) {
            return '';
        }
    }

    showPublishModal() {
        if (!this.currentDesign) {
            this.showNotification('请先生成一个设计作品', 'error');
            return;
        }
        
        document.getElementById('publishModal').style.display = 'flex';
        document.getElementById('productName').focus();
    }

    hidePublishModal() {
        document.getElementById('publishModal').style.display = 'none';
        document.getElementById('publishForm').reset();
    }

    async showProductDetails(productId) {
        try {
            const response = await this.apiRequest(`/products/${productId}`, 'GET', null, false);
            const product = response.product || response;
            
            const modal = document.getElementById('productDetailsModal');
            document.getElementById('detailProductName').textContent = product.name;
            document.getElementById('detailProductImage').src = product.image_url ? `/static${product.image_url}` : '';
            document.getElementById('detailProductDescription').textContent = product.description || '暂无描述';
            document.getElementById('detailProductPrice').textContent = `¥${product.base_price.toFixed(2)}`;
            document.getElementById('detailProductCategory').textContent = this.getCategoryName(product.category) || '通用';
            document.getElementById('detailProductCreator').textContent = product.creator_name || '匿名';
            
            if (product.design_prompt) {
                document.getElementById('detailDesignPrompt').textContent = `"${product.design_prompt}"`;
                document.getElementById('detailDesignPromptGroup').style.display = 'block';
            } else {
                document.getElementById('detailDesignPromptGroup').style.display = 'none';
            }
            
            if (product.design_style) {
                document.getElementById('detailDesignStyle').textContent = product.design_style;
                document.getElementById('detailDesignStyleGroup').style.display = 'block';
            } else {
                document.getElementById('detailDesignStyleGroup').style.display = 'none';
            }
            
            // 存储当前产品ID用于添加到购物车
            this.currentProductDetailId = product.id;
            modal.style.display = 'flex';
            
        } catch (error) {
            this.showNotification('加载产品详情失败: ' + error.message, 'error');
        }
    }

    hideProductDetailsModal() {
        document.getElementById('productDetailsModal').style.display = 'none';
        this.currentProductDetailId = null;
    }

    async handlePublishToShop(e) {
        e.preventDefault();
        
        if (!this.currentDesign) {
            this.showNotification('没有可发布的设计作品', 'error');
            return;
        }

        const productName = document.getElementById('productName').value.trim();
        const description = document.getElementById('productDescription').value.trim();
        const price = parseFloat(document.getElementById('productPrice').value);

        if (!productName) {
            this.showNotification('请输入商品名称', 'error');
            return;
        }

        if (!price || price <= 0) {
            this.showNotification('请输入有效的价格', 'error');
            return;
        }

        try {
            await this.apiRequest('/designs/publish', 'POST', {
                design_id: this.currentDesign.id,
                product_name: productName,
                description: description,
                price: price
            });

            this.showNotification('作品已成功发布到商店！', 'success');
            this.hidePublishModal();
            
            // 可选：跳转到商店页面
            setTimeout(() => {
                this.showSection('shopSection');
            }, 1500);

        } catch (error) {
            this.showNotification('发布失败: ' + error.message, 'error');
        }
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
        
        // If token expired and we have a refresh token, try to refresh
        if (response.status === 401 && auth && this.refreshToken) {
            const refreshed = await this.refreshAuthToken();
            if (refreshed) {
                // Retry the original request with new token
                options.headers['Authorization'] = `Bearer ${this.token}`;
                const retryResponse = await fetch(url, options);
                const retryResult = await retryResponse.json();
                
                if (!retryResponse.ok) {
                    throw new Error(retryResult.error || retryResult.message || '请求失败');
                }
                
                return retryResult;
            } else {
                // Refresh failed, redirect to login
                this.logout();
                this.showNotification('登录已过期，请重新登录', 'error');
                this.showAuthModal('login');
                throw new Error('登录已过期');
            }
        }
        
        const result = await response.json();
        
        if (!response.ok) {
            throw new Error(result.error || result.message || '请求失败');
        }
        
        return result;
    }

    async refreshAuthToken() {
        if (!this.refreshToken) {
            return false;
        }

        try {
            const response = await fetch(`${this.baseURL}/auth/refresh`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    refresh_token: this.refreshToken
                })
            });

            if (response.ok) {
                const result = await response.json();
                this.token = result.token;
                this.refreshToken = result.refresh_token;
                
                localStorage.setItem('authToken', this.token);
                localStorage.setItem('refreshToken', this.refreshToken);
                
                return true;
            }
        } catch (error) {
            console.error('Token refresh failed:', error);
        }

        return false;
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

    // ==================== 电商功能方法 ====================

    async loadProducts() {
        const productsGrid = document.getElementById('productsGrid');
        const category = document.getElementById('categoryFilter').value;
        
        productsGrid.innerHTML = `
            <div class="loading-placeholder">
                <i class="fas fa-spinner fa-spin"></i>
                <p>正在加载创意产品...</p>
            </div>
        `;

        try {
            let endpoint = '/products';
            if (category) {
                endpoint = `/products/category?category=${encodeURIComponent(category)}`;
            }
            
            const response = await this.apiRequest(endpoint, 'GET', null, false);
            
            // 确保返回的是数组格式
            let products = response.products || response;
            if (!Array.isArray(products)) {
                products = [];
            }
            
            this.renderProducts(products);
        } catch (error) {
            console.error('加载创意产品失败:', error);
            productsGrid.innerHTML = `
                <div class="error-placeholder">
                    <i class="fas fa-exclamation-triangle"></i>
                    <p>加载失败: ${error.message}</p>
                </div>
            `;
        }
    }

    renderProducts(products) {
        const productsGrid = document.getElementById('productsGrid');
        
        if (!products || products.length === 0) {
            productsGrid.innerHTML = `
                <div class="empty-placeholder">
                    <i class="fas fa-box-open"></i>
                    <p>暂无创意产品</p>
                </div>
            `;
            return;
        }

        productsGrid.innerHTML = products.map(product => {
            const imageContent = product.image_url ? 
                `<img src="/static${product.image_url}" alt="${product.name}" onload="this.style.display='block'" onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';">
                 <i class="fas fa-palette" style="display: none;"></i>` : 
                '<i class="fas fa-palette"></i>';
            
            return `
                <div class="product-card" data-product-id="${product.id}">
                    <div class="product-image" onclick="app.showProductDetails(${product.id})" style="cursor: pointer;">
                        ${imageContent}
                    </div>
                    <div class="product-info">
                        <div class="product-summary">
                            <div class="product-main-info">
                                <h3 class="product-name">${product.name}</h3>
                                <span class="product-creator">by ${product.creator_name || '匿名'}</span>
                                <span class="product-category-badge">${this.getCategoryName(product.category) || '通用'}</span>
                            </div>
                            <button class="cart-icon-btn" onclick="app.addToCartDirectly(${product.id})" title="添加到购物车">
                                <i class="fas fa-shopping-cart"></i>
                            </button>
                        </div>
                    </div>
                </div>
            `;
        }).join('');
    }

    async showAddToCartModal(productId) {
        if (!this.token) {
            this.showNotification('请先登录以添加创意产品到购物车', 'error');
            this.showAuthModal('login');
            return;
        }

        try {
            const response = await this.apiRequest(`/products/${productId}`, 'GET', null, false);
            this.currentProduct = response.product || response;
            
            // 填充模态框信息
            document.getElementById('cartProductName').textContent = this.currentProduct.name;
            document.getElementById('cartProductPrice').textContent = `¥${this.currentProduct.base_price.toFixed(2)}`;
            
            // 显示模态框
            document.getElementById('addToCartModal').style.display = 'flex';
            document.getElementById('sizeSelect').focus();
            
        } catch (error) {
            this.showNotification('加载创意产品信息失败: ' + error.message, 'error');
        }
    }

    hideAddToCartModal() {
        document.getElementById('addToCartModal').style.display = 'none';
        document.getElementById('addToCartForm').reset();
        this.currentProduct = null;
    }

    async addToCartDirectly(productId) {
        if (!this.token) {
            this.showNotification('请先登录以添加创意产品到购物车', 'error');
            this.showAuthModal('login');
            return;
        }

        try {
            // 获取用户的设计作品用于选择
            const designsResponse = await this.apiRequest('/designs/my-designs', 'GET');
            const designs = designsResponse.designs || [];

            if (designs.length === 0) {
                this.showNotification('请先创建一个设计作品才能购买商品', 'error');
                this.showSection('designSection');
                return;
            }

            // 使用最新的设计（最后一个）
            const design = designs[designs.length - 1];

            const cartItem = {
                product_id: parseInt(productId),
                design_id: parseInt(design.id),
                quantity: 1
            };

            await this.apiRequest('/cart/add', 'POST', cartItem);
            
            this.showNotification('创意产品已添加到购物车', 'success');
            this.updateCartBadge();
            
        } catch (error) {
            console.error('添加到购物车详细错误:', error);
            this.showNotification('添加到购物车失败: ' + error.message, 'error');
        }
    }

    async loadCart() {
        if (!this.token) {
            this.showCartLoginPrompt();
            return;
        }

        try {
            const response = await this.apiRequest('/cart', 'GET');
            this.renderCart(response);
        } catch (error) {
            console.error('加载购物车失败:', error);
            this.showCartError(error.message);
        }
    }

    renderCart(cartData) {
        const cartContent = document.getElementById('cartContent');
        const cartSummary = document.getElementById('cartSummary');
        
        if (!cartData.items || cartData.items.length === 0) {
            cartContent.innerHTML = `
                <div class="cart-empty">
                    <i class="fas fa-shopping-cart"></i>
                    <h3>购物车是空的</h3>
                    <p>快去商店挑选喜欢的商品吧！</p>
                    <button class="btn btn-primary" onclick="app.showSection('shopSection')">去购物</button>
                </div>
            `;
            cartSummary.style.display = 'none';
            return;
        }

        cartContent.innerHTML = cartData.items.map(item => `
            <div class="cart-item" data-item-id="${item.id}">
                <div class="cart-item-image">
                    <img src="${item.design?.image_url || '/static/images/placeholder-artwork.png'}" alt="创意设计">
                </div>
                <div class="cart-item-details">
                    <h4>${item.product?.name || '创意产品'}</h4>
                    <p class="design-prompt">${item.design?.prompt || '自定义创意'}</p>
                    <div class="item-price">¥${(item.product?.base_price * item.quantity).toFixed(2)}</div>
                </div>
                <div class="cart-item-controls">
                    <div class="quantity-controls">
                        <button class="btn-quantity" onclick="app.updateCartItemQuantity(${item.id}, ${item.quantity - 1})">-</button>
                        <span class="quantity">${item.quantity}</span>
                        <button class="btn-quantity" onclick="app.updateCartItemQuantity(${item.id}, ${item.quantity + 1})">+</button>
                    </div>
                    <button class="btn-remove" onclick="app.removeFromCart(${item.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `).join('');

        // 更新汇总信息
        document.getElementById('cartTotalItems').textContent = cartData.total_items || 0;
        document.getElementById('cartTotalAmount').textContent = `¥${(cartData.total_value || 0).toFixed(2)}`;
        cartSummary.style.display = 'block';
    }

    async updateCartItemQuantity(itemId, newQuantity) {
        if (newQuantity < 1) {
            await this.removeFromCart(itemId);
            return;
        }

        try {
            await this.apiRequest(`/cart/${itemId}`, 'PUT', { quantity: newQuantity });
            this.loadCart(); // 重新加载购物车
            this.updateCartBadge();
        } catch (error) {
            this.showNotification('更新数量失败: ' + error.message, 'error');
        }
    }

    async removeFromCart(itemId) {
        try {
            await this.apiRequest(`/cart/${itemId}`, 'DELETE');
            this.showNotification('商品已从购物车移除', 'success');
            this.loadCart(); // 重新加载购物车
            this.updateCartBadge();
        } catch (error) {
            this.showNotification('移除创意产品失败: ' + error.message, 'error');
        }
    }

    showCartLoginPrompt() {
        const cartContent = document.getElementById('cartContent');
        cartContent.innerHTML = `
            <div class="cart-login-prompt">
                <i class="fas fa-user-lock"></i>
                <h3>请先登录</h3>
                <p>登录后查看和管理购物车</p>
                <button class="btn btn-primary" onclick="app.showAuthModal('login')">立即登录</button>
            </div>
        `;
        document.getElementById('cartSummary').style.display = 'none';
    }

    showCartError(message) {
        const cartContent = document.getElementById('cartContent');
        cartContent.innerHTML = `
            <div class="cart-error">
                <i class="fas fa-exclamation-triangle"></i>
                <h3>加载失败</h3>
                <p>${message}</p>
                <button class="btn btn-outline" onclick="app.loadCart()">重试</button>
            </div>
        `;
        document.getElementById('cartSummary').style.display = 'none';
    }

    async updateCartBadge() {
        const badge = document.getElementById('cartBadge');
        if (this.token) {
            try {
                const cartResponse = await this.apiRequest('/cart', 'GET');
                const totalItems = cartResponse.total_items || 0;
                if (totalItems > 0) {
                    badge.style.display = 'inline';
                    badge.textContent = totalItems > 99 ? '99+' : totalItems.toString();
                } else {
                    badge.style.display = 'none';
                }
            } catch (error) {
                console.error('获取购物车数量失败:', error);
                badge.style.display = 'none';
            }
        } else {
            badge.style.display = 'none';
        }
    }

    showCheckoutModal() {
        document.getElementById('checkoutModal').style.display = 'flex';
        this.renderOrderSummary();
    }

    hideCheckoutModal() {
        document.getElementById('checkoutModal').style.display = 'none';
    }

    async renderOrderSummary() {
        try {
            const cartResponse = await this.apiRequest('/cart', 'GET');
            const orderSummary = document.getElementById('orderSummary');
            
            orderSummary.innerHTML = cartResponse.items.map(item => `
                <div class="order-item">
                    <img src="${item.design?.image_url || '/static/images/placeholder-artwork.png'}" alt="创意设计">
                    <div class="order-item-info">
                        <h4>${item.product?.name}</h4>
                        <p>数量: x${item.quantity}</p>
                    </div>
                    <div class="order-item-price">
                        ¥${(item.product?.base_price * item.quantity).toFixed(2)}
                    </div>
                </div>
            `).join('');

            document.getElementById('orderTotalAmount').textContent = 
                `¥${(cartResponse.total_value || 0).toFixed(2)}`;

        } catch (error) {
            console.error('加载订单汇总失败:', error);
        }
    }

    async createOrder() {
        try {
            const cartResponse = await this.apiRequest('/cart', 'GET');
            const cartItemIds = cartResponse.items.map(item => item.id);

            const orderData = {
                cart_item_ids: cartItemIds
            };

            await this.apiRequest('/orders', 'POST', orderData);
            
            this.showNotification('订单创建成功！', 'success');
            this.hideCheckoutModal();
            this.loadCart(); // 清空购物车
            this.updateCartBadge();
            
            // 跳转到订单页面
            setTimeout(() => {
                this.showSection('ordersSection');
            }, 1500);

        } catch (error) {
            this.showNotification('创建订单失败: ' + error.message, 'error');
        }
    }

    async loadOrders() {
        if (!this.token) {
            this.showOrdersLoginPrompt();
            return;
        }

        try {
            const response = await this.apiRequest('/orders', 'GET');
            
            // 确保返回的是数组格式
            let orders = response.orders || response;
            if (!Array.isArray(orders)) {
                orders = [];
            }
            
            this.renderOrders(orders);
        } catch (error) {
            console.error('加载订单失败:', error);
            this.showOrdersError(error.message);
        }
    }

    renderOrders(orders) {
        const ordersList = document.getElementById('ordersList');
        
        if (!orders || orders.length === 0) {
            ordersList.innerHTML = `
                <div class="orders-empty">
                    <i class="fas fa-receipt"></i>
                    <h3>暂无订单</h3>
                    <p>您还没有任何订单记录</p>
                    <button class="btn btn-primary" onclick="app.showSection('shopSection')">去购物</button>
                </div>
            `;
            return;
        }

        ordersList.innerHTML = orders.map(order => `
            <div class="order-card">
                <div class="order-header">
                    <div class="order-info">
                        <h3>订单号: ${order.order_sn}</h3>
                        <span class="order-date">${formatDate(order.created_at)}</span>
                    </div>
                    <div class="order-status ${order.status}">
                        ${this.getOrderStatusText(order.status)}
                    </div>
                </div>
                <div class="order-items">
                    ${order.order_items?.map(item => `
                        <div class="order-item">
                            <img src="${item.design_image_url}" alt="设计图案">
                            <div class="item-info">
                                <h4>${item.product_name}</h4>
                                <p>数量: x${item.quantity}</p>
                            </div>
                            <div class="item-price">¥${item.price.toFixed(2)}</div>
                        </div>
                    `).join('') || ''}
                </div>
                <div class="order-footer">
                    <div class="order-total">
                        总计: ¥${order.total_amount.toFixed(2)}
                    </div>
                </div>
            </div>
        `).join('');
    }

    getOrderStatusText(status) {
        const statusMap = {
            'pending': '待支付',
            'paid': '已支付',
            'shipped': '已发货',
            'completed': '已完成',
            'cancelled': '已取消'
        };
        return statusMap[status] || status;
    }

    showOrdersLoginPrompt() {
        const ordersList = document.getElementById('ordersList');
        ordersList.innerHTML = `
            <div class="orders-login-prompt">
                <i class="fas fa-user-lock"></i>
                <h3>请先登录</h3>
                <p>登录后查看订单记录</p>
                <button class="btn btn-primary" onclick="app.showAuthModal('login')">立即登录</button>
            </div>
        `;
    }

    showOrdersError(message) {
        const ordersList = document.getElementById('ordersList');
        ordersList.innerHTML = `
            <div class="orders-error">
                <i class="fas fa-exclamation-triangle"></i>
                <h3>加载失败</h3>
                <p>${message}</p>
                <button class="btn btn-outline" onclick="app.loadOrders()">重试</button>
            </div>
        `;
    }
}

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    window.app = new AICreativeStudio();
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