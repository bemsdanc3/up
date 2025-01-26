import { useState } from 'react'

function AddWin({ groupAdded, trackAdded, addInfo, close }) {
  const [file, setFile] = useState();
  const [title, setTitle] = useState('');
  const [genre, setGenre] = useState(1);

    const groupCrate = async () => {
      try {
        if (title.length <= 0) {
            return;
        } 
        let link;
        const formdata = new FormData();
        formdata.append('cover', file);
        formdata.append('title', title);
        if (addInfo.type == 'album') {
            link = `/album/create`
            formdata.append('type_of', 'type_of будет тут');
            formdata.append('label', 'label будет тут');
        } else {
            link = `/playlists/create`
            formdata.append('description', 'description будет тут');
            formdata.append('is_public', true);
        }
        const response = await fetch(`http://localhost:8080${link}`,{
          method: 'POST',
          credentials: 'include',
          withCredentials: true,
          body: formdata,
        });
        const responseData = await response.json();
        groupAdded();
        if (response.ok) {
          console.log("salamalekum")
          logFunc();
        } else {
          console.log(responseData.error);
        }
      } catch (error) {
        console.log(error);
      }
    } 

    const trackAdd = async () => {
      try {
        if (title.length <= 0 || !file) {
          return;
        }
    
        // Создаем временный URL для аудиофайла
        const audioUrl = URL.createObjectURL(file);
        const audio = new Audio(audioUrl);
    
        // Ждем загрузки метаданных
        audio.addEventListener('loadedmetadata', async () => {
          const duration = audio.duration; // Длительность в секундах
    
          const formdata = new FormData();
          formdata.append('trackFile', file);
          formdata.append('title', title);
          formdata.append('duration', Math.ceil(duration)); // Округляем до секунд
          formdata.append('album_id', addInfo.groupId);
          formdata.append('genre_id', genre);
    
          // Отправляем данные на сервер
          const response = await fetch(`http://localhost:8080/tracks/create`, {
            method: 'POST',
            credentials: 'include',
            withCredentials: true,
            body: formdata,
          });
    
          const responseData = await response.json();
    
          // Обработка ответа
          trackAdded();
          if (response.ok) {
            console.log('salamalekum');
            logFunc();
          } else {
            console.log(responseData.error);
          }
    
          // Освобождаем временный URL
          URL.revokeObjectURL(audioUrl);
        });
      } catch (error) {
        console.log(error);
      }
    };    

  return (
    <> 
        <div id="groupAdd">
            <div id="groupAddWin">
                <div id="groupAddTitleAndClose">
                    <h1>Создание {addInfo.type}</h1>
                    <button onClick={close}>X</button>
                </div>
                <input id='groupImgInput' type="file" placeholder='Выберите файл' onChange={(e)=>{setFile(e.target.files[0])}}/>
                <input id='groupTitleInput' type="text" placeholder='Введите название' onChange={(e)=>setTitle(e.target.value)}/>
                {(addInfo.type == 'album' || addInfo.type == 'playlist') &&
                  <button type='button' onClick={()=>groupCrate()}>Создать</button>
                }
                {addInfo.type == 'track' &&
                  <>
                    <select name="" id="" defaultValue={1} onChange={(e)=>setGenre(e.target.value)}>
                      <option value="1">Hip-Hop</option>
                      <option value="2">Gospel</option>
                      <option value="3">R&B</option>
                      <option value="4">Rock</option>
                      <option value="5">Pop</option>
                      <option value="6">Classic</option>
                      <option value="7">Jazz</option>
                      <option value="8">Metal</option>
                    </select>
                    <button type='button' onClick={()=>trackAdd()}>Добавить трек</button>
                  </>
                }
            </div>
        </div>
    </>
  )
}

export default AddWin
