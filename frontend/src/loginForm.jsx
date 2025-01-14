import { useState } from 'react'

function LoginForm() {
    const [showRegWin, setShowRegWin] = useState(true);

    const Register = async () => {
    try {
        const login = document.getElementById('loginInput')
        const email = document.getElementById('emailInput')
        const password = document.getElementById('passwordInput')
      const loginRes = await fetch('http://localhost:8080/register',{
        method: 'POST',
        credentials: 'include',
        withCredentials: true,
        headers: {
          "Content-Type": "application/json", 
        },
        body: JSON.stringify({
          email: email.value,
          pass: password.value,
          login: login.value
        }),
      });
      if (loginRes.ok) {
        console.log("salamalekum")
        setShowRegWin(false);
      } else {
        const errorData = await loginRes.json();
        console.log(errorData.error);
      }
    } catch (error) {
      console.log(error);
    }
    }

  return (
    <>
      {showRegWin && 
      <div id="registration">
        <h1>Регистрация</h1>
        <input id='emailInput' type="email" placeholder='Введите почту' />
        <input id='loginInput' type="text" placeholder='Введите логин' />
        <input id='passwordInput' type="password" placeholder='Введите пароль' />
        <button type='button' onClick={()=>Register()}>Зарегистрироваться</button>
      </div>}
      {!showRegWin && 
      <h1>мрилр</h1>
      }
    </>
  )
}

export default LoginForm
