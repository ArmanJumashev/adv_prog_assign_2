const API_URL = 'http://localhost:8080';

document.addEventListener('DOMContentLoaded', () => {
    console.log('DOMContentLoaded');
    const user = JSON.parse(localStorage.getItem('user'));
    console.log('user data:', user);
    if (!user) {
        alert('You must be logged in to access your profile.');
        window.location.href = '/static/login.html';
        return;
    }

    document.getElementById('email').value = user.email;
});

// Обновление информации пользователя
document.getElementById('updateInfoForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const fullName = document.getElementById('full_name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const dateOfBirth = document.getElementById('date_of_birth').value;

    try {
        const response = await fetch(`${API_URL}/update-user`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ full_name: fullName, email, password, date_of_birth: dateOfBirth })
        });

        if (response.ok) {
            const message = await response.text();
            alert('User info updated successfully!');
            localStorage.setItem('user', JSON.stringify({ full_name: fullName, email, date_of_birth: dateOfBirth })); // Обновляем все данные в localStorage
        } else {
            const error = await response.text();
            alert(`Failed to update user info: ${error}`);
        }
    } catch (error) {
        console.error('Error updating user info:', error);
        alert('An error occurred while updating your info.');
    }
});


// Загрузка списка заказов
async function loadOrders() {
    try {
        const response = await fetch(`${API_URL}/orders`);
        if (response.ok) {
            const orders = await response.json();
            const orderList = document.getElementById('orderList');
            orderList.innerHTML = '';

            orders.forEach((order) => {
                const li = document.createElement('li');
                li.textContent = `Order #${order.id} - ${order.total}`;
                orderList.appendChild(li);
            });
        }
    } catch (error) {
        console.error('Error fetching orders:', error);
    }
}
loadOrders();

// Отправка сообщения в поддержку
document.getElementById('supportForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const message = document.getElementById('message').value;
    const fileInput = document.getElementById('file');
    const formData = new FormData();

    formData.append('message', message);
    if (fileInput.files[0]) {
        formData.append('file', fileInput.files[0]);
    }

    try {
        const response = await fetch(`${API_URL}/support`, {
            method: 'POST',
            body: formData,
        });

        if (response.ok) {
            document.getElementById('supportMessage').innerText = 'Message sent successfully!';
        } else {
            const error = await response.text();
            document.getElementById('supportMessage').innerText = `Failed to send message: ${error}`;
        }
    } catch (error) {
        console.error('Error sending support message:', error);
        document.getElementById('supportMessage').innerText = 'An error occurred while sending the message.';
    }
});
