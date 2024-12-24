let currentPage = 1;
let currentAuthor = '';
let currentSort = 'DESC';

// Функция для отправки GET-запроса с параметрами
function fetchPosts() {
    var url = `/getPost?`
    if (currentAuthor !== '') {
      url = url + "author=" + currentAuthor
    }

    if ((currentAuthor !== "") && currentSort !== ""){
      url +="&"
    }
    if (currentSort !== ""){
      url = url + "sort=" + currentSort
    }
    if ((currentAuthor !== "" || currentSort !== "") && currentPage !== 0){
      url +="&"
    }
    if (currentPage !== 0){
      url = url + "page=" +currentPage
    }
    // const url = `/getPost?author=${currentAuthor}&sortDate=${currentSort}&page=${currentPage}`;
    fetch(url)
        .then(response => response.json())
        .then(data => {
            displayPosts(data.posts); // Отображение полученных постов
            })
        .catch(error => console.error('Error fetching posts:', error));
        }

        // Функция для отображения постов
function displayPosts(posts) {
    const postsDiv = document.getElementById('posts');
    if (posts.length <= 0) {
      console.log("no page (")
      document.getElementById("secret_p").textContent = "no page :("
      return
    }
    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.innerHTML = `<div class="post"><p>${post.content}</p><p>Author: ${post.author}</p><p>Date: ${post.created_at}</p></div>`;
        postsDiv.appendChild(postElement);
    });
    currentPage++;
}

        // Применение фильтров и сортировки
function applyFilters() {
    const postsDiv = document.getElementById('posts');
  
    currentAuthor = document.getElementById('author').value; // Получаем значение фильтра автора
    currentSort = document.getElementById('sortDate').value; // Получаем значение сортировки
    postsDiv.innerHTML = ''; // clear previuos posts
    currentPage = 1 // reload page
    fetchPosts(); // Загружаем посты с новыми параметрами
}

        // Функция для смены страницы пагинации
function changePage() {
    fetchPosts(); // Загружаем посты для новой страницы
}
        // Инициализация страницы с начальными параметрами
fetchPosts();