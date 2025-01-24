import { useState } from 'react'

function LoginForm({ logFunc, loginData }) {
  const [isLogin, setIsLogin] = useState(true);

    const Register = async () => {
    try {
        const login = document.getElementById('loginInput')
        const email = document.getElementById('emailInput')
        const password = document.getElementById('passwordInput')
      const loginRes = await fetch(`http://localhost:8080/register`,{
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
      const responseData = await loginRes.json();
      if (loginRes.ok) {
        console.log("salamalekum")
        logFunc();
        loginData(responseData);
      } else {
        const errorData = await loginRes.json();
        console.log(errorData.error);
      }
    } catch (error) {
      console.log(error);
    }
    }

    const Login = async () => {
      try {
          const email = document.getElementById('emailInput')
          const password = document.getElementById('passwordInput')
        const loginRes = await fetch(`http://localhost:8080/login`,{
          method: 'POST',
          credentials: 'include',
          withCredentials: true,
          headers: {
            "Content-Type": "application/json", 
          },
          body: JSON.stringify({
            email: email.value,
            pass: password.value,
          }),
        });
        const responseData = await loginRes.json();
        loginData(responseData);
        console.log("responseData");
        console.log(responseData);
        if (loginRes.ok) {
          console.log("salamalekum")
          logFunc();
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
      {!isLogin &&
        <div id="registration">
        <h1>Регистрация</h1>
        <input id='emailInput' type="email" placeholder='Введите почту' />
        <input id='loginInput' type="text" placeholder='Введите логин' />
        <input id='passwordInput' type="password" placeholder='Введите пароль' />
        <button type='button' onClick={()=>Register()}>Зарегистрироваться</button>
        <span onClick={()=> setIsLogin(!isLogin)}>Есть аккаунт</span>
        </div>
      }
      {isLogin &&
        <div id="registration">
        <h1>Авторизация</h1>
        <input id='emailInput' type="email" placeholder='Введите почту' />
        <input id='passwordInput' type="password" placeholder='Введите пароль' />
        <button type='button' onClick={()=>Login()}>Войти</button>
        <span onClick={()=> setIsLogin(!isLogin)}>Нет аккаунта</span>
      </div>
      }
    </>
  )
}

export default LoginForm
