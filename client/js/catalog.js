const API_URL = 'http://localhost:8080';
let currentPage = 1;

// Function to fetch products
async function fetchProducts(page = 1, category = '', sort = '', order = '') {
    try {
        const response = await fetch(
            `${API_URL}/catalog?page=${page}&category=${category}&sort_by=${sort}&order=${order}`
        );
        const data = await response.json();

        console.log(data);

        const productsDiv = document.getElementById('products');
        if (data.products && data.products.length > 0) {
            productsDiv.innerHTML = data.products
                .map(
                    (p) =>
                        `<div class="product-card">
                            <img src="${p.image}" alt="${p.name} image" class="product-image"/>
                            <div class="product-info">
                                <h3>${p.name}</h3>
                                <p><strong>Price:</strong> $${p.price}</p>
                                <p><strong>Category:</strong> ${p.category}</p>
                                <p><strong>Description:</strong> ${p.description}</p>
                                <button class="add-to-cart" data-id="${p.id}">Add to Cart</button>
                            </div>
                        </div>`
                )
                .join('');
        } else {
            productsDiv.innerHTML = '<p>No products found.</p>';
        }

        document.getElementById('prevPage').disabled = data.prevPage === 0;
        document.getElementById('nextPage').disabled = data.nextPage === 0;

        document.querySelectorAll('.add-to-cart').forEach((button) => {
            button.addEventListener('click', (e) => {
                const productId = e.target.getAttribute('data-id');
                addToCart(productId);
            });
        });

    } catch (error) {
        console.error('Error fetching products:', error);
    }
}

// Функция для добавления товара в корзину
function addToCart(productId) {
    let cart = JSON.parse(localStorage.getItem('cart')) || [];

    // Проверяем, есть ли товар в корзине
    const existingProduct = cart.find(item => item.id === productId);
    if (existingProduct) {
        existingProduct.quantity += 1; // Увеличиваем количество, если товар уже в корзине
    } else {
        cart.push({ id: productId, quantity: 1 });
    }

    // Сохраняем обновленную корзину в localStorage
    localStorage.setItem('cart', JSON.stringify(cart));

    alert('Product added to cart!');
}

// Handle filters
document.getElementById('filterForm')?.addEventListener('submit', (e) => {
    e.preventDefault();
    const category = document.getElementById('category').value;
    const sort = document.getElementById('sort').value;
    const order = document.getElementById('order').value;

    currentPage = 1; // Reset to the first page
    fetchProducts(currentPage, category, sort, order);
});

// Handle pagination
document.getElementById('prevPage')?.addEventListener('click', () => {
    if (currentPage > 1) {
        currentPage--;
        fetchProducts(currentPage);
    }
});

document.getElementById('nextPage')?.addEventListener('click', () => {
    currentPage++;
    fetchProducts(currentPage);
});

// Initial fetch to display all products
fetchProducts();
