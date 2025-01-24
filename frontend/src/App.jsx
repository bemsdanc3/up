import { useState } from 'react'
import LoginForm from './loginForm.jsx'
import Header from './Header.jsx'
import Profile from './Profile.jsx'
import Home from './Home.jsx'
import Player from './Player.jsx'

function App() {
  const [logged, setLogged] = useState(false);
  const [track, setTrack] = useState(null);
  const [isWavePlaying, setIsWavePlaying] = useState(false); // состояние волны
  const [isTrackPlaying, setIsTrackPlaying] = useState(false); // состояние волны
  const [userProfile, setUserProfile] = useState('');
  const [loginData, setLoginData] = useState();

  const handleProfile = (userid) => {
    setUserProfile(userid);
  };

  const handleHome = () => {
    setUserProfile('');
  };

  const playWave = () => {
    setIsWavePlaying(!isWavePlaying); // Переключаем состояние волны
    setTrack(null); // Сбрасываем текущий трек
  };
  
  const playTrack = (newTrack) => {
    setIsWavePlaying(false); // Останавливаем волну
    setTrack(newTrack); // Устанавливаем текущий трек
  };   

  return (
    <>
      {!logged && <LoginForm logFunc={() => setLogged(true)} loginData={(data)=>setLoginData(data)}/>}
      {logged && <Header profileFunc={() => handleProfile('self')} homeFunc={handleHome} loginData={loginData}/>}
      {logged && userProfile && <Profile userId={userProfile} />}
      {logged && !userProfile && (
        <Home 
          trackPlay={playTrack} 
          playWave={playWave} 
        />
      )}
      {logged && (
        <Player 
          track={track} 
          isWavePlaying={isWavePlaying} 
          userProfile={handleProfile} 
        />
      )}
    </>
  );
}

export default App
