import { useState, useRef, useEffect } from 'react';
import './player.css';

function Player() {
  const audioRef = useRef(null);
  const progressRef = useRef(null);
  const [progress, setProgress] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const [volume, setVolume] = useState(1); // volume from 0 to 1

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
    if (isPlaying) {
      audioRef.current.pause();
    } else {
      audioRef.current.play();
    }
    setIsPlaying(!isPlaying);
  };

  // Функция для выключения музыки
//   const stopMusic = () => {
//     audioRef.current.pause();
//     audioRef.current.currentTime = 0;
//     setIsPlaying(false);
//   };

  // Функция для изменения громкости
  const handleVolumeChange = (event) => {
    const newVolume = event.target.value;
    audioRef.current.volume = newVolume;
    setVolume(newVolume);
  };

  return (
    <div id="player">
      <div id="trackinfo">
        <img src="" alt="" />
        <div id="trackTextInfo">
          <span id="tracktitle">asd</span>
          <span id="trackauthor">fassdsa</span>
        </div>
      </div>
      <div id="audioplayer">
        <audio
          ref={audioRef}
          src="http://localhost:8080/uploads/tracks/cd424f05-f784-46b2-b045-27432e0deb87.mp3"
          preload="metadata"
        ></audio>
        <div id="playercontrols">
            <button onClick={togglePlay}>
            {isPlaying ? 'Pause' : 'Play'}
            </button>
            {/* <button onClick={stopMusic}>Stop</button> */}
        </div>
        <div
          id="progress-bar-container"
          ref={progressRef}
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
      </div>
    </div>
  );
}

export default Player;
