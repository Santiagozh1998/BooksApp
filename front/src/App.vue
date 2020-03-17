<template>
  <div id="app">
        <b-navbar type="dark" variant="danger">
            <b-navbar-brand href="#">VueGO</b-navbar-brand>
        </b-navbar>
        <b-container class="container-books">
          <b-row>
            <b-col lg="4">
              <div class="card text-center">
                <div class="card-header">
                  <h3>Book</h3>
                </div>
                <div class="card-body">
                  <b-form @submit="bookInfo">  
                    <input v-model="book.Title" class="form-control-lg form-control input-form" type="text" placeholder="Title">
                    <input v-model="book.Author" class="form-control-lg form-control input-form" type="text" placeholder="Author">
                    <b-button block variant="danger" type="submit" class="button-form">Enviar</b-button>
                  </b-form>
                </div>
              </div>
            </b-col>
            <b-col lg="8">
              <div class="card text-center">
                <div class="card-header">
                  <h3>Books</h3>
                </div>
                <div class="card-body">
                  <table class="table table-striped table-bordered">
                    <thead>
                      <tr>
                        <th>Title</th>
                        <th>Author</th>
                        <th>Options</th>
                      </tr>
                    </thead>
                    <tbody v-if="Books.length == 0">
                        <tr>
                          <td>Vacio</td>
                          <td>Vacio</td>
                          <td>Vacio</td>
                        </tr>
                    </tbody>
                    <tbody v-if="Books.length > 0">
                        <tr v-bind:key="book._id" v-for="book in Books">
                          <td>{{book.title}}</td>
                          <td>{{book.author}}</td>
                          <td>
                            <b-button v-on:click="updateBook(book._id)" class="table-button" variant="primary">Edit</b-button>
                            <b-button v-on:click="deleteBook(book._id)" class="table-button" variant="danger">Delete</b-button>
                          </td>
                        </tr>
                    </tbody>
                  </table>
                </div>
            </div>
            </b-col>
          </b-row>
        </b-container>
    </div>
</template>

<script>
import axios from "axios"

export default {
  data() {
    return{
      ApiUrl: "https://restapi-books.herokuapp.com",
      Books: [],
      book: {
        Id: null,
        Title: null,
        Author: null
      },
      state: 0
    }
  },
  created() {

    let vue = this;
    axios.get(vue.ApiUrl + "/api/books/")
    .then(function(response) {
      vue.Books = response.data;
    })
    .catch(function(error) {
      console.log(error)
    })
  },
  methods: {
    bookInfo(e) {

      e.preventDefault()

      if(this.state == 0){

        if(this.book.Title != null && this.book.Author != null){
        
          let vue = this;
          axios.post(vue.ApiUrl + "/api/books/", {
            title: vue.book.Title,
            author: vue.book.Author
          },
          {    
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json'
            }
          })
          .then(function(response) {
            vue.Books.push({_id: response.data, title: vue.book.Title, author: vue.book.Author})
            vue.book.Title = null
            vue.book.Author = null
          })
          .catch(function(error) {
            console.log(error)
          })
        }
      }
      else {

        if(this.book.Title != null && this.book.Author != null){

          this.state = 0
          let vue = this;
          axios.put(vue.ApiUrl + "/api/books/" + vue.book.Id, {
            title: vue.book.Title,
            author: vue.book.Author
          },
          {    
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json'
            }
          })
          .then(function(response) {
            
            if(response.data == 1){

              for(var i = 0; i< vue.Books.length; i++){
                if(vue.Books[i]._id == vue.book.Id){
                  vue.Books[i].title = vue.book.Title,
                  vue.Books[i].author = vue.book.Author
                  break;
                }
              }
              
              vue.book.Id = null
              vue.book.Title = null
              vue.book.Author = null
            }
          })
          .catch(function(error) {
            console.log(error)
          })
        }
      }
      
    },
    updateBook(idBook){
      this.state = 1
      this.book.Id = idBook
      var book = this.Books.find(book => book._id == idBook)
      this.book.Title = book.title
      this.book.Author = book.author
    },
    deleteBook(idBook){

      let vue = this;
      axios.delete(vue.ApiUrl + "/api/books/" + idBook)
      .then(function(response) {
        if(response.data == 1){
          var temp = []

          for(var i = 0; i< vue.Books.length; i++){
            if(vue.Books[i]._id != idBook){
              temp.push(vue.Books[i])
            }
          }
          vue.Books = temp
        }
      })
      .catch(function(error) {
        console.log(error)
      })
    }
  }
}
</script>

<style>
.container-books{
    margin-top: 30px;
}

.input-form{
    font-family: 'Roboto', sans-serif;
    margin-top: 10px;
    margin-bottom: 10px;
}

.button-form{
    margin-top: 20px;
}

.table-button{
  margin-left: 5px;
  margin-right: 5px;
}
</style>
