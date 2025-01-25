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
  const [backHome, setBackHome] = useState(false);
  const [track, setTrack] = useState(null);
  const [isWavePlaying, setIsWavePlaying] = useState(false); // состояние волны
  const [isTrackPlaying, setIsTrackPlaying] = useState(false); // состояние волны
  const [userProfile, setUserProfile] = useState('');
  const [loginData, setLoginData] = useState();
  const [groupInfo, setGroupInfo] = useState();
  const [playlist, setPlaylist] = useState();
  const [album, setAlbum] = useState();
  const [page, setPage] = useState('Home')

  const handleProfile = (userid) => {
    setUserProfile(userid);
    setPage('Profile');
  };
  
  const handleHome = () => {
    setPage('Home');
    setBackHome(true);
    setInterval(()=>{
      setBackHome(false);
    }, 0)
    setPlaylist();
    setAlbum();
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
  
  const playTrack = (newTrack, groupInfo) => {
    setIsWavePlaying(false); // Останавливаем волну
    setTrack(newTrack); // Устанавливаем текущий трек
    if (groupInfo) {
      setGroupInfo(groupInfo);
    }
  };   

  const playlistFunc = (id) => {
    console.log('playlistFunc, playlist.id: ' + id)
    setAlbum();
    setPlaylist(id); 
    setPage('Home');
  }

  const albumFunc = (id) => {
    console.log('albumFunc, album.id: ' + id)
    setPlaylist();
    setAlbum(id); 
    setPage('Home');
  }

  return (
    <>
      {!logged && 
        <LoginForm 
          logFunc={() => setLogged(true)} 
          loginData={(data)=>{setLoginData(data); console.log(loginData)}}
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
        <Profile 
          userId={userProfile} 
          myData={loginData} 
          handleTrackClick={playTrack}
          playlistFunc={playlistFunc}
          albumFunc={albumFunc}
        />
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
          loginData={loginData}
          playlist={playlist}
          album={album}
          backHome={backHome}
        />
      }
      {logged && 
        <Player 
          track={track} 
          isWavePlaying={isWavePlaying} 
          userProfile={handleProfile} 
          groupInfo={groupInfo}
        />
      }
    </>
  );
}

export default App
