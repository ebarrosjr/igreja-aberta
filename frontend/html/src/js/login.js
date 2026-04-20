        const signinBtn = document.getElementById('signin-btn');
        const forgotPassword = document.getElementById('forgot-password');

        if (signinBtn) {
            signinBtn.addEventListener('click', () => {
                const email = document.getElementById('login-email').value.trim();
                const password = document.getElementById('login-password').value.trim();
                const rememberMe = document.getElementById('remember-me').checked;

                if (!email || !password) {
                    alert('Preencha o email e a senha.');
                    return;
                }

                alert(`Login
Email: ${email}
Senha: ${'*'.repeat(password.length)}
Mantenha-me logado: ${rememberMe ? 'Sim' : 'Não'}`);
            });
        }

        if (forgotPassword) {
            forgotPassword.addEventListener('click', (e) => {
                e.preventDefault();
                const email = document.getElementById('login-email').value.trim();
                if (!email) {
                    alert('Digite seu email para redefinir a senha.');
                    return;
                }
                alert(`Link de recuperação enviado para ${email}.`);
            });
        }