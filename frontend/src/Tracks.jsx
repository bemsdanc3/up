import { useState, useEffect } from 'react'

function Tracks({ handleTrackClick}) {
  const [tracks, setTracks] = useState(true);

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
    getAllTracks();
  }, [])

  const formatDuration = (duration) => {
    const minutes = Math.floor(duration / 60);
    const seconds = duration % 60;
    return `${minutes}:${seconds.toString().padStart(2, '0')}`; // Форматирует секунды с ведущим нулём
  };

  const trackDelete = async (trackid) => {
    try {
      const trackDelRes = await fetch(`http://localhost:8080/tracks/delete/${trackid}`,{
        method: 'DELETE',
        credentials: 'include',
        withCredentials: true,
      });
      const responseData = await trackDelRes.json();
      console.log("responseData");
      console.log(responseData);
      if (trackDelRes.ok) {
        console.log("track deleted");
        getAllTracks();
      } else {
        const errorData = await loginRes.json();
        console.log(errorData.error);
      }
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <> 
      <div id="tracksPage">
        <h2>Треки: </h2>
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
                    <button onClick={(e)=>{e.stopPropagation(); trackDelete(track.id)}}>
                        Удалить
                    </button>
                    </div>
                </div>
                )
            })
            } 
      </div>
    </>
  )
}

export default Tracks
