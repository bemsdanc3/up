import { useState, useRef, useEffect } from 'react';
import './player.css';
import PlayIcon from './assets/PlayIcon.svg?react';
import NextIcon from './assets/NextIcon.svg?react';
import PauseIcon from './assets/PauseIcon.svg?react';
import VolumeIcon from './assets/VolumeIcon.svg?react';
import VolumeOffIcon from './assets/VolumeOffIcon.svg?react';
import HeartIcon from './assets/HeartIcon.svg?react';
import HeartFilledIcon from './assets/HeartFilledIcon.svg?react';
import PlusIcon from './assets/PlusIcon.svg?react';

function Player({ track, userProfile, isWavePlaying, groupInfo, playlistUpd }) {
  const audioRef = useRef(null);
  const progressRef = useRef(null);
  const [playlists, setPlaylists] = useState([]);
  const [playlistsLoaded, setPlaylistsLoaded] = useState(false);
  const [progress, setProgress] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const [isTrackLiked, setIsTrackLiked] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const [trackLoaded, setTrackLoaded] = useState(false);
  const [volume, setVolume] = useState(1); // volume from 0 to 1
  const [isMuted, setIsMuted] = useState(false);
  const [playingTrackInfo, setPlayingTrackInfo] = useState({});

  // Обновляем прогресс каждую секунду
  useEffect(() => {
    const updateProgress = () => {
      if (audioRef.current && !isDragging) {
        const progress = (audioRef.current.currentTime / audioRef.current.duration) * 100;
        setProgress(progress);
      }
      requestAnimationFrame(updateProgress);
    };
    requestAnimationFrame(updateProgress);

    return () => cancelAnimationFrame(updateProgress);
  }, [isDragging]);

  // Обрабатываем клик по ползунку для изменения времени воспроизведения
  const handleProgressClick = (event) => {
    const { left, width } = progressRef.current.getBoundingClientRect();
    const clickPosition = (event.clientX - left) / width;
    const seekTime = clickPosition * audioRef.current.duration;
    audioRef.current.currentTime = seekTime;
  };

  // Обрабатываем начало перетаскивания ползунка
  const handleMouseDown = () => {
    setIsDragging(true);
  };

  // Обрабатываем завершение перетаскивания
  const handleMouseUp = () => {
    setIsDragging(false);
  };

  // Обрабатываем перетаскивание ползунка
  const handleMouseMove = (event) => {
    if (isDragging) {
      const { left, width } = progressRef.current.getBoundingClientRect();
      const movePosition = (event.clientX - left) / width;
      const seekTime = movePosition * audioRef.current.duration;
      audioRef.current.currentTime = seekTime;
      setProgress(movePosition * 100);
    }
  };

  // Функция для воспроизведения/паузы
  const togglePlay = () => {
    if (playingTrackInfo?.track_link) {
      if (isPlaying) {
        audioRef.current.pause();
      } else {
        audioRef.current.play();
      }
      setIsPlaying(!isPlaying);
    } else {
      getRandomTrack();
    }
  };

  const isTrackLikedFunc = async (trackid) => {
    const response = await fetch('http://localhost:8080/playlists/track/check', {
      method: 'POST',
      credentials: 'include',
      withCredentials: true,
      headers: {
        "Content-Type": "application/json", 
      },
      body: JSON.stringify({
        track_id: trackid
      })
    })
    if (response.ok) {
      const responseData = await response.json();
      if (responseData.is_in_playlist == true) {
        setIsTrackLiked(true);
        console.log('track liked')
      } else {
        setIsTrackLiked(false);
        console.log('track not liked')
      }
    } else {
      console.log('жопа');
    }
  }
  
  useEffect(() => {
    console.log(playingTrackInfo.track_link);
    isTrackLikedFunc(playingTrackInfo.id);

    if (audioRef.current && playingTrackInfo.track_link) {
      audioRef.current.src = playingTrackInfo.track_link;
      audioRef.current.load(); // Загрузка трека
      audioRef.current.play();
      setIsPlaying(true);
    }
  }, [playingTrackInfo]);

  

  useEffect(() => {
    if (track) {
      if (groupInfo) {
        setPlayingTrackInfo({ 
          ...track, 
          author_id: groupInfo.author_id, 
          author_login: track.author_login || 
          groupInfo.authorLogin, 
          cover: track.cover || 
          groupInfo.cover });
        isTrackLikedFunc(track.id);
      } else {
        setPlayingTrackInfo(track);
      }
    }
  }, [track]);

  useEffect(()=>{
    if (isWavePlaying) {
      getRandomTrack();
    }
  }, [isWavePlaying])

  // Функция для изменения громкости
  const handleVolumeChange = (event) => {
    const newVolume = event.target.value;
    audioRef.current.volume = newVolume;
    setVolume(newVolume);
    if (newVolume > 0) setIsMuted(false);
  };

  // Функция для включения/выключения звука
  const toggleMute = () => {
    if (isMuted) {
      audioRef.current.volume = volume; // Восстановить предыдущую громкость
    } else {
      audioRef.current.volume = 0; // Выключить звук
    }
    setIsMuted(!isMuted);
  };

  const getRandomTrack = async () => {
    const response = await fetch('http://localhost:8080/tracks/random', {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    if (response.ok) {
      const responseData = await response.json();
      console.log(responseData);
      setIsPlaying(true);
      setPlayingTrackInfo(responseData);
      setTrackLoaded(true);
      isTrackLikedFunc();
    } else {
      console.log('жопа');
    }
  }

  useEffect(() => {
    const audio = audioRef.current;
    const handleError = () => {
      console.error('Ошибка загрузки трека:', audio?.error);
    };
    audio.addEventListener('error', handleError);
  
    return () => {
      audio.removeEventListener('error', handleError);
    };
  }, []); 
  
  useEffect(() => {
    const audio = audioRef.current;
  
    const handleTrackEnd = () => {
      console.log("Трек завершён, запускаем новый трек...");
      getRandomTrack(); // Запуск нового трека
    };
  
    const handleError = () => {
      console.error('Ошибка загрузки трека:', audio?.error);
    };
  
    if (audio) {
      audio.addEventListener('ended', handleTrackEnd);
      audio.addEventListener('error', handleError);
    }
  
    return () => {
      if (audio) {
        audio.removeEventListener('ended', handleTrackEnd);
        audio.removeEventListener('error', handleError);
      }
    };
  }, [getRandomTrack]);  

  const trackLike = async (trackId) => {
    try {
      const loginRes = await fetch(`http://localhost:8080/tracks/like`,{
        method: 'POST',
        credentials: 'include',
        withCredentials: true,
        headers: {
          "Content-Type": "application/json", 
        },
        body: JSON.stringify({
          track_id: trackId,
        }),
      });
      const responseData = await loginRes.json();
      console.log("responseData");
      console.log(responseData);
      if (loginRes.ok) {
        console.log("salamalekum")
        isTrackLikedFunc(playingTrackInfo.id);
      } else {
        const errorData = await loginRes.json();
        console.log(errorData.error);
      }
    } catch (error) {
      console.log(error);
    }
  }  
  
  const getAllPlaylists = async () => {
    console.log('getting all tracks');
    const response = await fetch('http://localhost:8080/playlists/my/all', {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    // console.log(response);
    const responseData = await response.json();
    // console.log(responseData);
    if (response.ok) {
      setPlaylists(responseData);
      setPlaylistsLoaded(true);
    } else {
      console.log('жопа');
    }
  }

  useEffect(()=>{
    if (playlistUpd) {
      getAllPlaylists();
    }
  }, [playlistUpd]);

  useEffect(()=>{
    getAllPlaylists();
  }, [])

  const playlistAddTrack = async (playlist_id) => {
    try {
      const playlistAddTrackSelect = document.getElementById('playlistAddTrackSelect');
      console.log("playlistAddTrackSelect.value: " + playlistAddTrackSelect.value);
      console.log("playingTrackInfo.id: " + playingTrackInfo.id)
      const response = await fetch(`http://localhost:8080/playlists/add/track`,{
        method: 'POST',
        credentials: 'include',
        withCredentials: true,
        headers: {
          "Content-Type": "application/json", 
        },
        body: JSON.stringify({
          playlist_id: Number(playlistAddTrackSelect.value),
          track_id: playingTrackInfo.id,
        }),
      });
      const responseData = await response.json();
      console.log("responseData");
      console.log(responseData);
      if (response.ok) {
        console.log("salamalekum");
      } else {
        const errorData = await response.text();
        console.log(errorData.error);
        setErrText(errorData);
      }
    } catch (error) {
      console.log(error);
    }
  }    

  return (
    <div id="player">
      <div id="trackinfo">
        {playingTrackInfo &&
          <>
            {playingTrackInfo.cover && <img src={playingTrackInfo.cover} alt="" />}
            <div id="trackTextInfo">
              {playingTrackInfo.title && <span id="tracktitle">{playingTrackInfo.title}</span>}
              <span id="trackauthor" onClick={()=>{userProfile(playingTrackInfo.author_id)}}>{playingTrackInfo.author_login}</span>
            </div>
            <div id="trackLikeDiv">
            {isTrackLiked && playingTrackInfo.title && 
              <HeartFilledIcon onClick={()=>{trackLike(playingTrackInfo.id)}}/>
            }
            {!isTrackLiked && playingTrackInfo.title && 
              <HeartIcon onClick={()=>{trackLike(playingTrackInfo.id)}}/>
            }
            </div>
            {playingTrackInfo.title && playlists && playlists.length >= 1 && playlists.find((el)=>el.title != 'Любимые треки') &&
            <div id="playlistAddTrackDiv">
              <select name="" id="playlistAddTrackSelect">
                {playlists.map((playlist)=>{
                  if (playlist.title != 'Любимые треки') {
                    return (
                      <option value={playlist.id}>{playlist.title}</option>
                    )
                  }
                })}
              </select>
              <button 
                onClick={playlistAddTrack}
              >
                <PlusIcon />
              </button>
            </div>
            }
          </>
        }
      </div>
      <div id="audioplayer">
        {playingTrackInfo &&
        <audio
          ref={audioRef}
          src={playingTrackInfo.track_link}
          preload="metadata"
        ></audio>}
        <div id="playercontrols">
          <button onClick={togglePlay}>
            {isPlaying ? <PauseIcon /> : <PlayIcon />}
            {isPlaying ? 'Пауза' : 'Играть'}
          </button>
          <button
            id="nextBtn"
            onClick={()=>{
              getRandomTrack();
            }}
          >
            <NextIcon />
            Далее
          </button>
        </div>
        <div
          id="progress-bar-container"
          ref={progressRef}
          className={(isPlaying ? 'isPlaying ' : 'notPlaying ') + (playingTrackInfo ? 'trackLoaded' : 'trackNotLoaded')}
          onClick={handleProgressClick}
          onMouseDown={handleMouseDown}
          onMouseUp={handleMouseUp}
          onMouseMove={handleMouseMove}
        >
            <div
            id="progress-bar"
            style={{ width: `${progress}%` }}
          ></div>
        </div>
      </div>
      <div id="playersettings">
        <div
          className="volume-slider"
          onMouseDown={(e) => {
            const slider = e.currentTarget;
            const { left, width } = slider.getBoundingClientRect();
            const updateVolume = (event) => {
              const clickPosition = (event.clientX - left) / width;
              const newVolume = Math.min(Math.max(clickPosition, 0), 1); // Ограничиваем от 0 до 1
              setVolume(newVolume);
              audioRef.current.volume = newVolume;
              if (newVolume > 0) setIsMuted(false);
            };

            const handleMouseMove = (event) => updateVolume(event);
            const handleMouseUp = () => {
              window.removeEventListener('mousemove', handleMouseMove);
              window.removeEventListener('mouseup', handleMouseUp);
            };

            updateVolume(e); // Обновляем громкость при первом клике
            window.addEventListener('mousemove', handleMouseMove);
            window.addEventListener('mouseup', handleMouseUp);
          }}
        >
          <div
            className="volume-progress"
            style={{ width: `${volume * 100}%` }}
          ></div>
        </div>
        <button onClick={toggleMute}>
          {isMuted ? <VolumeOffIcon /> : <VolumeIcon />}
          {isMuted ? 'Звук выкл.' : 'Звук'}
        </button>
      </div>
    </div>
  );
}

export default Player;
