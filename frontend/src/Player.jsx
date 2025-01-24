import { useState, useRef, useEffect } from 'react';
import './player.css';
import PlayIcon from './assets/PlayIcon.svg?react';
import NextIcon from './assets/NextIcon.svg?react';
import PauseIcon from './assets/PauseIcon.svg?react';
import VolumeIcon from './assets/VolumeIcon.svg?react';
import VolumeOffIcon from './assets/VolumeOffIcon.svg?react';

function Player({ track, userProfile, isWavePlaying }) {
  const audioRef = useRef(null);
  const progressRef = useRef(null);
  const [progress, setProgress] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
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
  
  useEffect(() => {
    console.log(playingTrackInfo.track_link);

    if (audioRef.current && playingTrackInfo.track_link) {
      audioRef.current.src = playingTrackInfo.track_link;
      audioRef.current.load(); // Загрузка трека
      audioRef.current.play();
      setIsPlaying(true);
    }
  }, [playingTrackInfo]);

  

  useEffect(() => {
    if (track) {
      setPlayingTrackInfo(track);
    }
  }, [track]);

  useEffect(()=>{
    getRandomTrack();
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

  return (
    <div id="player">
      <div id="trackinfo">
        {trackLoaded && playingTrackInfo &&
          <>
            <img src={playingTrackInfo.cover} alt="" />
            <div id="trackTextInfo">
              <span id="tracktitle">{playingTrackInfo.title}</span>
              <span id="trackauthor" onClick={()=>{userProfile(playingTrackInfo.author_id)}}>{playingTrackInfo.author_login}</span>
            </div>
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
            {isPlaying ? 'Pause' : 'Play'}
          </button>
          <button
            id="nextBtn"
            onClick={()=>{
              getRandomTrack();
            }}
          >
            <NextIcon />
            Next
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
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          value={volume}
          onChange={handleVolumeChange}
          className="volume-slider"
        />
        <button onClick={toggleMute}>
          {isMuted ? <VolumeOffIcon/> : <VolumeIcon/>}
          {isMuted ? 'Unmute' : 'Mute'}
        </button>
      </div>
    </div>
  );
}

export default Player;
