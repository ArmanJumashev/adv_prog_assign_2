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

async function submitOrder() {
    const cart = JSON.parse(localStorage.getItem('cart')) || [];
    const user = JSON.parse(localStorage.getItem('user'));

    if (cart.length === 0) {
        alert('Your cart is empty.');
        return;
    }
    const order = {
        user: user,
        items: cart.map(item => ({
            product_id: item.id,
            quantity: item.quantity
        })),
        totalPrice: calculateTotalPrice(cart),
        shippingAddress: getShippingAddress(),
    };

    try {
        const response = await fetch(`${API_URL}/order`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(order),
        });

        if (response.ok) {
            const data = await response.json();
            alert('Your order has been placed successfully!');
            localStorage.removeItem('cart');  // Очищаем корзину после отправки заказа
            window.location.href = '/static/order-confirmation.html';  // Перенаправляем на страницу подтверждения
        } else {
            const error = await response.json();
            alert(`Failed to submit order: ${error.message}`);
        }
    } catch (error) {
        console.error('Error submitting order:', error);
        alert('An error occurred while submitting your order.');
    }
}

function calculateTotalPrice(cart) {
    return cart.reduce((total, item) => {
        const product = getProductById(item.id);  // Получаем товар по ID
        return total + product.price * item.quantity;
    }, 0).toFixed(2);
}

function getShippingAddress() {
    return {
        street: document.getElementById('street').value,
        city: document.getElementById('city').value,
        zipCode: document.getElementById('zipCode').value,
    };
}

document.getElementById('checkoutButton').addEventListener('click', (e) => {
    e.preventDefault();
    submitOrder();
});


document.addEventListener('DOMContentLoaded', renderCart);
