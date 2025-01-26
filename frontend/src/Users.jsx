import { useEffect, useState } from 'react'
import SearchIcon from './assets/SearchIcon.svg?react';

function Users({ userProfile, userData }) {
  const [users, setUsers] = useState([]);
  const [filteredUsers, setFilteredUsers] = useState([]);
  const [searchQuery, setSearchQuery] = useState(''); // Состояние для поиска
  
  const getAllUsers = async () => {
    console.log('getting all tracks');
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    console.log(response);
    const responseData = await response.json();
    console.log('users responseData');
    console.log(responseData);
    if (response.ok) {
      setUsers(responseData);
      setFilteredUsers(responseData);
    } else {
      console.log('жопа');
    }
  }

  useEffect(()=>{
    getAllUsers();
    console.log(userData)
  }, [])

  const updUserRole = async (role, id) => {
    try {
      const updUserRoleRes = await fetch(`http://localhost:8080/users/update-role`,{
        method: 'PATCH',
        credentials: 'include',
        withCredentials: true,
        headers: {
          "Content-Type": "application/json", 
        },
        body: JSON.stringify({
          role: role,
          user_id: id
        }),
      });
      const responseData = await updUserRoleRes.json();
      console.log("responseData");
      console.log(responseData);
      if (updUserRoleRes.ok) {
        console.log("salamalekum")
        getAllUsers();
      } else {
        const errorData = await updUserRoleRes.json();
        console.log(errorData.error);
      }
    } catch (error) {
      console.log(error);
    }
  }

  const filterFunc = (role = 'all', search = searchQuery) => {
    // Сброс классов кнопок
    const userAdminButton = Array.from(document.getElementsByClassName('userAdminButton'));
    userAdminButton.forEach((btn) => {
      btn.classList.remove('selected');
    });
  
    if (role !== 'all') {
      const selectedButton = document.getElementById(`${role}UsersBtn`);
      if (selectedButton) selectedButton.classList.add('selected');
    } else {
      const allButton = document.getElementById('allUsersBtn');
      if (allButton) allButton.classList.add('selected');
    }
  
    // Фильтрация пользователей
    let filtered = users;
  
    if (role !== 'all') {
      filtered = filtered.filter((user) => user.role === role);
    }
  
    if (search) {
      filtered = filtered.filter((user) =>
        user.login.toLowerCase().includes(search.toLowerCase())
      );
    }
  
    setFilteredUsers(filtered);
  };
  
  const handleSearch = (e) => {
    const query = e.target.value;
    setSearchQuery(query); // Обновляем значение поиска
    filterFunc('all', query); // Фильтрация по всем пользователям с обновлённым запросом
  };  

  return (
    <> 
      <div id="usersPage">
        <h2>Пользователи: </h2>
        <div id="userAdminBtns">
          <button id="allUsersBtn" className='userAdminButton selected' onClick={() => filterFunc('all')}>Все</button>
          <button id="adminUsersBtn" className='userAdminButton ' onClick={() => filterFunc('admin')}>Admin</button>
          <button id="artistUsersBtn" className='userAdminButton ' onClick={() => filterFunc('artist')}>Artist</button>
          <button id="userUsersBtn" className='userAdminButton ' onClick={() => filterFunc('user')}>User</button>
          <div id="searchUsersDiv">
            <input
              type="text"
              id="searchUsers"
              placeholder="Поиск пользователей..."
              value={searchQuery}
              onChange={handleSearch} // Обработчик изменения
            />
            <SearchIcon />
          </div>
        </div>
        <div id="usersList">
            {filteredUsers && filteredUsers.length >=1 &&
                filteredUsers.map((user)=>{
                    return (
                        <div className="user" key={user.id} onClick={()=>userProfile(user.id)}>
                            <div className="userLeftInfo">
                                <img src={user.pfp} alt="" />
                                <div className="userCardTextInfo">
                                    <span className='userLoginAndRole'>
                                        <span className='usrCardLogin'>{user.login}</span>
                                        <span className='usrCardRole'>{user.role}</span>
                                    </span>
                                    <span>{user.email}</span>
                                </div>
                            </div>
                            {userData.user_id != user.id &&
                            <div className="userAdminBtns">
                                <button 
                                    onClick={(e)=>{e.stopPropagation(); updUserRole('admin', user.id)}}
                                    className={user.role == 'admin' ? 'selected' : ''}
                                >
                                    Admin
                                </button>
                                <button 
                                    onClick={(e)=>{e.stopPropagation(); updUserRole('artist', user.id)}}
                                    className={user.role == 'artist' ? 'selected' : ''}
                                >
                                    Artist
                                </button>
                                <button 
                                    onClick={(e)=>{e.stopPropagation(); updUserRole('user', user.id)}}
                                    className={user.role == 'user' ? 'selected' : ''}
                                >
                                    User
                                </button>
                            </div>}
                        </div>
                    )
                })
            }
        </div>
      </div>
    </>
  )
}

export default Users
