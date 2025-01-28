import { useEffect, useState } from 'react'

function LoginForm({ loginData }) {
  const [isLogin, setIsLogin] = useState(true);
  const [errText, setErrText] = useState('');

  function validateString(str) {
    // Проверка на отсутствие символов * & { } | +
    const forbiddenChars = /[*&{}|+]/;
    if (forbiddenChars.test(str)) {
      setErrText('В пароле не должно быть символов * & { } | +')
      return false;
    }

    // Проверка на наличие хотя бы одной заглавной буквы
    const hasUpperCase = /[A-Z]/;
    if (!hasUpperCase.test(str)) {
      setErrText('В пароле должна содержаться заглавная буква')
      return false;
    }

    // Проверка на наличие хотя бы одной цифры
    const hasNumber = /\d/;
    if (!hasNumber.test(str)) {
      setErrText('В пароле должна содержаться цифра')
      return false;
    }

    return true;
  }

  const Register = async () => {
    const login = document.getElementById('loginInput').value
    const email = document.getElementById('emailInput').value
    const password = document.getElementById('passwordInput').value
    if (login.length < 1 || email.length < 1 || password.length < 1) {
      setErrText('Заполните все данные')
    } else if (!validateString(password)) {
      return
    } else {
      const loginRes = await fetch(`http://localhost:8080/register`,{
        method: 'POST',
        credentials: 'include',
        withCredentials: true,
        headers: {
          "Content-Type": "application/json", 
        },
        body: JSON.stringify({
          email: email,
          pass: password,
          login: login
        }),
      });
      const responseData = await loginRes.json();
      console.log(responseData);
      if (loginRes.ok) {
        console.log("salamalekum")
        loginData(responseData);
      } else {
        const errorData = await loginRes.text();
        setErrText(errorData);
      }
    }
  }

  const authTry = async () => {
    const response = await fetch(`http://localhost:8080/refresh-token`,{
      method: 'POST',
      credentials: 'include',
      withCredentials: true,
    });
    const responseData = await response.json();
    console.log(responseData);
    if (response.ok) {
      console.log("salamalekum")
      loginData(responseData);
    } else {
      const errorData = await response.text();
      setErrText(errorData);
    }
  }

  useEffect(()=>{
    authTry();
  },[])

  const Login = async () => {
    try {
        const email = document.getElementById('emailInput')
        const password = document.getElementById('passwordInput')
        if (email.value.length < 1 || password.value.length < 1) {
          setErrText('Заполните все данные')
        } else {
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
          console.log("responseData");
          console.log(responseData);
          if (loginRes.ok) {
            console.log("salamalekum")
            loginData(responseData);
          } else {
            const errorData = await loginRes.text();
            console.log(errorData.error);
            setErrText(errorData);
          }
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
        <span id='regErr' className={`${errText ? 'visible' : 'hidden'}`}>{errText}</span>
        <input id='emailInput' type="email" placeholder='Введите почту' />
        <input id='loginInput' type="text" placeholder='Введите логин' />
        <input id='passwordInput' type="password" placeholder='Введите пароль' />
        <button type='button' onClick={()=>Register()}>Зарегистрироваться</button>
        <span onClick={()=> {setIsLogin(!isLogin); setErrText('')}}>Есть аккаунт</span>
        </div>
      }
      {isLogin &&
        <div id="registration">
        <h1>Авторизация</h1>
        <span id='regErr' className={`${errText ? 'visible' : 'hidden'}`}>{errText}</span>
        <input id='emailInput' type="email" placeholder='Введите почту' />
        <input id='passwordInput' type="password" placeholder='Введите пароль' />
        <button type='button' onClick={()=>Login()}>Войти</button>
        <span onClick={()=> {setIsLogin(!isLogin); setErrText('')}}>Нет аккаунта</span>
      </div>
      }
    </>
  )
}

export default LoginForm
