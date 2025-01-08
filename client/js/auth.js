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

    const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
    });

    const message = await response.text();
    document.getElementById('message').innerText = message;
});

