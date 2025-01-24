import { useState } from 'react'
import LoginForm from './loginForm.jsx'
import Header from './Header.jsx'
import Profile from './Profile.jsx'
import Home from './Home.jsx'
import Player from './Player.jsx'
import Users from './Users.jsx'
import Tracks from './Tracks.jsx'

function App() {
  const [logged, setLogged] = useState(false);
  const [track, setTrack] = useState(null);
  const [isWavePlaying, setIsWavePlaying] = useState(false); // состояние волны
  const [isTrackPlaying, setIsTrackPlaying] = useState(false); // состояние волны
  const [userProfile, setUserProfile] = useState('');
  const [loginData, setLoginData] = useState();
  const [page, setPage] = useState('Home')

  const handleProfile = (userid) => {
    setUserProfile(userid);
    setPage('Profile');
  };
  
  const handleHome = () => {
    setPage('Home');
    setUserProfile('');
  };
  
  const handleUsers = () => {
    setPage('Users');
    setUserProfile('');
  };
  
  const handleTracks = () => {
    setPage('Tracks');
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
      {!logged && 
        <LoginForm 
          logFunc={() => setLogged(true)} 
          loginData={(data)=>setLoginData(data)}
        />
      }
      {logged && 
        <Header 
          profileFunc={() => handleProfile(loginData.user_id)} 
          homeFunc={handleHome} 
          loginData={loginData}
          usersFunc={handleUsers}
          tracksFunc={handleTracks}
        />
      }
      {logged && userProfile && page == 'Profile' &&
        <Profile userId={userProfile} myData={loginData} handleTrackClick={playTrack}/>
      }
      {logged && page == 'Users' &&
        <Users 
          userProfile={(id)=>handleProfile(id)}
          userData={loginData}
        />
      }
      {logged && page == 'Tracks' &&
        <Tracks 
        handleTrackClick={playTrack}
        />
      }
      {logged && page == 'Home' && 
        <Home 
          trackPlay={playTrack} 
          playWave={playWave} 
        />
      }
      {logged && 
        <Player 
          track={track} 
          isWavePlaying={isWavePlaying} 
          userProfile={handleProfile} 
        />
      }
    </>
  );
}

export default App
