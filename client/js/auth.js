const API_URL = 'https://adv-prog-assign-2.onrender.com';

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
            localStorage.setItem('user', JSON.stringify({ full_name: fullName, email, date_of_birth: dateOfBirth }));

            document.getElementById('message').innerText = 'Registration successful. Please confirm your email.';

            document.getElementById('registerForm').style.display = 'none';
            document.getElementById('emailConfirmationForm').style.display = 'block';

            document.getElementById('formTitle').innerText = 'Confirm Your Email';

            const emailData = {
                to: email,
                subject: 'Confirm',
                body: 'Your code:'          
            };

            try {
                const response = await fetch('https://adv-prog-assign-2.onrender.com', { // Замените на ваш реальный серверный URL
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(emailData),
                });

                if (response.ok) {
                    console.log('successfully sent confirmation code')
                } else {
                    const errorText = await response.text();
                }
            } catch (error) {
                console.error('Error sending confirmation message:', error);
            }
        }
    } catch (error) {
        console.error('Error during registration:', error);
        document.getElementById('message').innerText = 'An error occurred during registration.';
    }
});


document.getElementById('emailConfirmationForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const confirmationCode = document.getElementById('confirmation_code').value;

    // Извлекаем email пользователя из localStorage
    const user = JSON.parse(localStorage.getItem('user'));
    if (!user || !user.email) {
        document.getElementById('message').innerText = 'Ошибка: Email пользователя не найден.';
        return;
    }

    const userEmail = user.email;
    console.log('extracted user email from local storage: ', userEmail);

    try {
        const response = await fetch(`${API_URL}/confirm-email`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-User-Email': userEmail,
            },
            body: JSON.stringify({ confirmation_code: confirmationCode }),
        });

        const message = await response.text();

        if (response.ok) {
            document.getElementById('message').innerText = 'Email confirmed successfully!';
            document.getElementById('emailConfirmationForm').style.display = 'none';
        } else {
            document.getElementById('message').innerText = message;
        }
    } catch (error) {
        console.error('Error during email confirmation:', error);
        document.getElementById('message').innerText = 'An error occurred during email confirmation.';
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

