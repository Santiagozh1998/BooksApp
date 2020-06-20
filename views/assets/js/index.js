const app = new Vue({
    el: '#app',
    data: {
        load: true,
        errors: [],
        books: [],
        book: {
            ID: null,
            Title: null,
            Author: null
        },
        currentState: 0
    },
    created() {
        axios.get("/api/books/all")
        .then(response => {
            console.log(response.data)
            this.books = response.data;
            this.load = false;
        })
        .catch(error => console.log(error))
    },
    methods: {
        checkForm(e) {
            
            e.preventDefault();

            if(this.book.Title != null && this.book.Author != null ) {

                if(this.currentState == 0) { //CREATE 

                    axios.post("/api/books/create", {
                        Title: this.book.Title,
                        Author: this.book.Author
                    })
                    .then(response => {

                        this.books.push({ID: response.data.ID, Title: response.data.Title, Author: response.data.Author})
                        this.book = {
                            ID: null,
                            Title: null,
                            Author: null
                        };
                    })
                    .catch(error => console.log(error))
                } 
                else { //UPDATE
    
                    this.currentState = 0;

                    axios.put("/api/books/update/" + this.book.ID, {
                        Title: this.book.Title,
                        Author: this.book.Author
                    })
                    .then(response => {

                        for(let i = 0; i< this.books.length; i++){

                            if(this.books[i].ID == response.data.ID){

                                this.books[i].Title = response.data.Title,
                                this.books[i].Author = response.data.Author
                                break;
                            }
                        }
                          
                        this.book = {
                            ID: null,
                            Title: null,
                            Author: null
                        };
                    })
                    .catch(error => console.log(error))
                }
            }
        },
        updateBook(idbook) {
            this.currentState = 1;
            this.book.ID = idbook;
            let book = this.books.find(book => book.ID == idbook);
            this.book.Title = book.Title;
            this.book.Author = book.Author;
        },
        deleteBook(idbook) {
            axios.delete("/api/books/delete/" + idbook)
            .then(response => {
                var temp = []

                for(let i = 0; i< this.books.length; i++){
                    if(this.books[i].ID != response.data.ID){
                        temp.push(this.books[i])
                    }
                }
                this.books = temp
            })
            .catch(error => console.log(error))
        }
    }
})
