import { useState, useEffect } from 'react'
import PlusIcon from './assets/PlusIcon.svg?react';
import PlayIcon from './assets/PlayIcon.svg?react';
import PauseIcon from './assets/PauseIcon.svg?react';

function Home({ trackPlay, playWave, loginData, playlist, album }) {
  const [tracks, setTracks] = useState([])
  const [tracksLoaded, setTracksLoaded] = useState(false)
  const [playlists, setPlaylists] = useState([])
  const [albums, setAlbums] = useState([])
  const [playlistsLoaded, setPlaylistsLoaded] = useState(false)
  const [albumsLoaded, setAlbumsLoaded] = useState(false)
  const [selectedAlbum, setSelectedAlbum] = useState();
  const [selectedPlaylist, setSelectedPlaylist] = useState();
  const [tracksGroup, setTracksGroup] = useState([])
  const [selectedGroupInfo, setSelectedGroupInfo] = useState({});

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

  useEffect(()=>{
    setSelectedAlbum(album);
    getGroupInfo(album, 'album')
  }, [album])
  
  useEffect(()=>{
    setSelectedPlaylist(playlist)
    getGroupInfo(playlist, 'playlist')
  }, [playlist])

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

  const getAllAlbums = async () => {
    console.log('getting all tracks');
    const response = await fetch(`http://localhost:8080/album/all/${loginData.user_id}`, {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    console.log(response);
    const responseData = await response.json();
    console.log(responseData);
    if (response.ok) {
      setAlbums(responseData);
      setAlbumsLoaded(true);
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
    if (!albumsLoaded) {
      getAllAlbums();
    }
  }, [])

  const formatDuration = (duration) => {
    const minutes = Math.floor(duration / 60);
    const seconds = duration % 60;
    return `${minutes}:${seconds.toString().padStart(2, '0')}`; // Форматирует секунды с ведущим нулём
  };

  const getGroupInfo = async (selectedid, type) => {
    console.log('getting playlist info and tracks');
    console.log(type)
    console.log(selectedid)
    setTracksGroup([]);
    setSelectedGroupInfo({});
    let link;
    let link2;
    if (type == 'album') {
      setSelectedPlaylist();
      setSelectedAlbum(selectedid);
      link =  `http://localhost:8080/album/${selectedid}`
      link2 =  `http://localhost:8080/tracks/by/album/${selectedid}`
    } else {
      setSelectedAlbum();
      setSelectedPlaylist(selectedid);
      link =  `http://localhost:8080/playlists/view/${selectedid}`
      link2 =  `http://localhost:8080/tracks/by/playlist/${selectedid}`
    }
    const response = await fetch(link, {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    console.log(response);
    const responseData = await response.json();
    console.log(responseData);
    if (response.ok) {
      setSelectedGroupInfo(responseData);
      const tracksResponse = await fetch(link2, {
        method: 'GET',
        credentials: 'include',
        withCredentials: true,
      })
      console.log(tracksResponse);
      const responseTracksData = await tracksResponse.json();
      console.log(responseTracksData);
      if (response.ok) {
        setTracksGroup(responseTracksData);
      } else {
        console.log('жопа');
      }
    } else {
      console.log('жопа');
    }
  }

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
                {loginData.role == 'artist' &&
                <button>
                  <PlusIcon />
                  Новый Альбом
                </button>
                }
              </div>
            <div id="myPlaylists">
            {albums && albums.length >=1 &&
              <>
              <h2>Мои альбомы: </h2>
              {albums.map((album)=>{
                return (
                  <>
                  <div className="album" onClick={()=>{getGroupInfo(album.id, 'album')}}>
                    <img src={album.cover} alt="" />
                    <div className="playlistTextInfo">
                      <span>{album.title}</span>
                    </div>
                  </div>
                  </>
                )
              })}
              </>
            }
            </div>
            <div id="myPlaylists">
            {playlists && playlists.length >=1 &&
              <>
              <h2>Мои плейлисты: </h2>
              {playlists.map((playlist)=>{
                return (
                  <>
                  <div className="playlist" onClick={()=>{getGroupInfo(playlist.id, 'playlist')}}>
                    <img src={playlist.cover} alt="" />
                    <div className="playlistTextInfo">
                      <span>{playlist.title}</span>
                    </div>
                  </div>
                  </>
                )
              })}
              </>
            }
            </div>
          </div>
          {!selectedAlbum && !selectedPlaylist &&
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
                        <span>{track.album_title}</span>
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
          }
          {(selectedAlbum || selectedPlaylist) &&
          <div id="selectedGroup">
              <div id='selectedGroupInfo'>
                <img src={selectedGroupInfo.cover} alt="" />
                <div id="selectedGroupTextInfo">
                  <h2>{selectedGroupInfo.title}</h2>
                  <h2>{selectedGroupInfo.authorLogin}</h2>
                </div>
              </div>
              <div id='selectedTracksGroup'>
                {tracksGroup && tracksGroup.length >= 1 &&
                  tracksGroup.map((track)=>{
                    return (
                      <div className="track" key={track.id}>
                        <div className="trackLeftInfo">
                          {track.title}
                        </div>
                        <div className="trackRightInfo">
                          {formatDuration(track.duration)}
                        </div>
                      </div>
                    )
                  })
                }
              </div>
          </div>
          }
        </div>
    </>
  )
}

export default Home
