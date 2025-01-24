import { useState, useEffect } from 'react'
import PlusIcon from './assets/PlusIcon.svg?react';
import PlayIcon from './assets/PlayIcon.svg?react';
import PauseIcon from './assets/PauseIcon.svg?react';

function Home({ trackPlay, playWave }) {
  const [tracks, setTracks] = useState([])
  const [tracksLoaded, setTracksLoaded] = useState(false)
  const [playlists, setPlaylists] = useState([])
  const [playlistsLoaded, setPlaylistsLoaded] = useState(false)

  const getAllTracks = async () => {
    console.log('getting all tracks');
    const response = await fetch('http://localhost:8080/tracks/all', {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    console.log(response);
    const responseData = await response.json();
    console.log(responseData);
    if (response.ok) {
      setTracks(responseData);
      setTracksLoaded(true);
    } else {
      console.log('жопа');
    }
  }

  const getAllPlaylists = async () => {
    console.log('getting all tracks');
    const response = await fetch('http://localhost:8080/playlists/my/all', {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    console.log(response);
    const responseData = await response.json();
    console.log(responseData);
    if (response.ok) {
      setPlaylists(responseData);
      setPlaylistsLoaded(true);
    } else {
      console.log('жопа');
    }
  }

  useEffect(()=>{
    if (!tracksLoaded) {
      getAllTracks();
    }
    if (!playlistsLoaded) {
      getAllPlaylists();
    }
  }, [])

  const formatDuration = (duration) => {
    const minutes = Math.floor(duration / 60);
    const seconds = duration % 60;
    return `${minutes}:${seconds.toString().padStart(2, '0')}`; // Форматирует секунды с ведущим нулём
  };

  const handleWavePlay = () => {
    playWave();
  };
  
  const handleTrackClick = (track) => {
    trackPlay(track);
  };  

  return (
    <>
        <div id="homePage">
          <div id="myMedia">
              <h2>Моя медиатека</h2>
              <div id="createBtns">
                <button>
                  <PlusIcon />
                  Новый плейлист
                </button>
                <button>
                  <PlusIcon />
                  Новый Альбом
                </button>
              </div>
            <div id="myPlaylists">
            {playlists && playlists.length >=1 &&
              playlists.map((playlist)=>{
                return (
                  <>
                  <div className="playlist">
                    <img src={playlist.cover} alt="" />
                    <div className="playlistTextInfo">
                      <span>{playlist.title}</span>
                    </div>
                  </div>
                  </>
                )
              })
            }
            </div>
          </div>
          <div id="homeRight">
            <div id="myWave">
              <button id="wavePlayBtn" onClick={()=>{handleWavePlay()}}>
                {/* {!isWavePlaying ? <PlayIcon /> : <PauseIcon />} */}
                <PlayIcon />
                Играть
              </button>
              <div id="homeGenres">
                <button id="RockGanre">Рок</button>
                <button id="PopGanre">Поп</button>
                <button id="MetalGanre">Метал</button>
                <button id="ClassicalGanre">Классика</button>
                <button id="HipHopGanre">Хип-Хоп</button>
                {/* <button id="MusicForSexGanre">Music for sex</button> */}
              </div>
            </div>
            <div id="allTracks">
              {tracks && tracks.length >= 1 &&
                tracks.map((track)=>{
                  return (
                    <div className="track"
                      onClick={()=>{handleTrackClick(track)}}
                    >
                      <div className="trackLeftInfo">
                        <img src={track.cover} alt="" />
                        <div className="trackTextInfo">
                          <span>{track.title}</span>
                          <span>{track.author_login}</span>
                        </div>
                      </div>
                      <div className="trackAlbumInfo">
                        <span> AlbumTitle </span>
                      </div>
                      <div className="trackDurationInfo">
                        <span>
                          {formatDuration(track.duration)}
                        </span>
                      </div>
                    </div>
                  )
                })
              } 
            </div>
          </div>
        </div>
    </>
  )
}

export default Home
