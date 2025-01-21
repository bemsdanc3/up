import { useState } from 'react'
import LoginForm from './loginForm.jsx'
import Header from './Header.jsx'
import Profile from './Profile.jsx'
import Home from './Home.jsx'
import Player from './Player.jsx'

function App() {
  const [logged, setLogged] = useState(false);
  const [isProfile, setIsProfile] = useState(false);

  const handleProfile = () => {
    console.log("Profile function called");
    setIsProfile(true);
  };

  const handleHome = () => {
    console.log("Home function called");
    setIsProfile(false);
  };

  return (
    <>
      {!logged &&
        <LoginForm logFunc={()=>{setLogged(true)}}/>
      }
      {logged && 
        <Header profileFunc={handleProfile} homeFunc={handleHome} />
      }
      {logged && isProfile &&
        <Profile />
      }
      {logged && !isProfile &&
        <Home />
      }
      {logged &&
      <Player/>}
    </>
  )
}

export default App
