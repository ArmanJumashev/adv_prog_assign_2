const API_URL = 'http://localhost:8080';

// Регистрация
document.getElementById('registerForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Получаем значения из полей формы
    const fullName = document.getElementById('full_name').value; // full_name
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const dateOfBirth = document.getElementById('date_of_birth').value; // date_of_birth

    try {
        const response = await fetch(`${API_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ full_name: fullName, email, password, date_of_birth: dateOfBirth }),
        });

        const message = await response.text();
        document.getElementById('message').innerText = message;

        if (response.ok) {
            // Сохраняем данные пользователя в localStorage
            localStorage.setItem('user', JSON.stringify({ full_name: fullName, email, date_of_birth: dateOfBirth }));
        }
    } catch (error) {
        console.error('Error during registration:', error);
        document.getElementById('message').innerText = 'An error occurred during registration.';
    }
});

// Авторизация
document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        // Проверяем статус ответа
        if (response.ok) {
            // Получаем JSON с данными пользователя
            const user = await response.json();
            console.log('response: ', response);
            console.log(user);
            // Сохраняем данные пользователя в localStorage
            localStorage.setItem('user', JSON.stringify(user));
            document.getElementById('message').innerText = 'Login successful!';
        } else {
            // Если ошибка, читаем текст ответа
            const errorMessage = await response.text();
            document.getElementById('message').innerText = errorMessage;
        }
    } catch (error) {
        console.error('Error during login:', error);
        document.getElementById('message').innerText = 'An error occurred during login.';
    }
});

