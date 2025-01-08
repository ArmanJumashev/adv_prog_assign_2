const API_URL = 'http://localhost:8080';

async function renderCart() {
    console.log('loading cart...');
    const cart = JSON.parse(localStorage.getItem('cart')) || [];
    const cartItemsDiv = document.getElementById('cartItems');
    const totalPriceSpan = document.getElementById('totalPrice');
    let totalPrice = 0;

    if (cart.length === 0) {
        cartItemsDiv.innerHTML = '<p>Your cart is empty.</p>';
        totalPriceSpan.textContent = '0.00';
        return;
    }

    // Clear previous cart items
    cartItemsDiv.innerHTML = '';

    for (const item of cart) {
        const product = await getProductById(item.id); // Use await to get the product data
        console.log('setting products in cart list...');
        console.log(product);

        if (!product) continue;

        const totalItemPrice = product.price * item.quantity;
        totalPrice += totalItemPrice;

        const cartItemHTML = `
            <div class="cart-item">
                <img src="${product.image}" alt="${product.name}" class="cart-item-image"/>
                <div class="cart-item-info">
                    <h3>${product.name}</h3>
                    <p><strong>Price:</strong> $${product.price}</p>
                    <p><strong>Quantity:</strong>
                        <input type="number" min="1" value="${item.quantity}" class="quantity-input" data-id="${item.id}"/>
                    </p>
                    <p><strong>Total:</strong> $${totalItemPrice.toFixed(2)}</p>
                    <button class="remove-item" data-id="${item.id}">Remove</button>
                </div>
            </div>
        `;

        cartItemsDiv.innerHTML += cartItemHTML;
    }

    totalPriceSpan.textContent = totalPrice.toFixed(2);

    // Add event listeners for quantity input and remove item buttons
    document.querySelectorAll('.quantity-input').forEach(input => {
        input.addEventListener('input', (e) => updateItemQuantity(e.target));
    });

    document.querySelectorAll('.remove-item').forEach(button => {
        button.addEventListener('click', (e) => removeItemFromCart(e.target.getAttribute('data-id')));
    });
}


function updateItemQuantity(input) {
    const productId = input.getAttribute('data-id');
    const cart = JSON.parse(localStorage.getItem('cart')) || [];
    const product = cart.find(item => item.id === productId);
    if (product) {
        product.quantity = parseInt(input.value) || 1;
        localStorage.setItem('cart', JSON.stringify(cart));
        renderCart();
    }
}

function removeItemFromCart(productId) {
    let cart = JSON.parse(localStorage.getItem('cart')) || [];
    cart = cart.filter(item => item.id !== productId);
    localStorage.setItem('cart', JSON.stringify(cart));
    renderCart();
}

function getProductById(id) {
    return fetch(`${API_URL}/product?id=${id}`)
        .then(response => response.json())
        .then(data => {
            console.log('(product by id) data :', data);
            if (data) {
                console.log('Product:', data);
                return data;
            } else {
                console.error('Failed to fetch product:', data);
            }
        })
        .catch(error => {
            console.error('Error fetching product:', error);
        });
}


document.addEventListener('DOMContentLoaded', renderCart);
